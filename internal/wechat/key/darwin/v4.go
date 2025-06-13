package darwin

import (
	"bytes"
	"context"
	"encoding/hex"
	"runtime"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/sjzar/chatlog/internal/errors"
	"github.com/sjzar/chatlog/internal/wechat/decrypt"
	"github.com/sjzar/chatlog/internal/wechat/key/darwin/glance"
	"github.com/sjzar/chatlog/internal/wechat/model"
)

const (
	MaxWorkers        = 8
	MinChunkSize      = 1 * 1024 * 1024 // 1MB
	ChunkOverlapBytes = 1024            // Greater than all offsets
	ChunkMultiplier   = 2               // Number of chunks = MaxWorkers * ChunkMultiplier
)

var V4KeyPatterns = []KeyPatternInfo{
	{
		Pattern: []byte{0x20, 0x66, 0x74, 0x73, 0x35, 0x28, 0x25, 0x00},
		Offsets: []int{16, -80, 64},
	},
}

type V4Extractor struct {
	validator   *decrypt.Validator
	keyPatterns []KeyPatternInfo
}

func NewV4Extractor() *V4Extractor {
	return &V4Extractor{
		keyPatterns: V4KeyPatterns,
	}
}

func (e *V4Extractor) Extract(ctx context.Context, proc *model.Process) (string, error) {
	if proc.Status == model.StatusOffline {
		return "", errors.ErrWeChatOffline
	}

	// Check if SIP is disabled, as it's required for memory reading on macOS
	if !glance.IsSIPDisabled() {
		return "", errors.ErrSIPEnabled
	}

	if e.validator == nil {
		return "", errors.ErrValidatorNotSet
	}

	// Create context to control all goroutines
	searchCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create channels for memory data and results
	memoryChannel := make(chan []byte, 100)
	resultChannel := make(chan string, 1)

	// Determine number of worker goroutines
	workerCount := runtime.NumCPU()
	if workerCount < 2 {
		workerCount = 2
	}
	if workerCount > MaxWorkers {
		workerCount = MaxWorkers
	}
	log.Debug().Msgf("Starting %d workers for V4 key search", workerCount)

	// Start consumer goroutines
	var workerWaitGroup sync.WaitGroup
	workerWaitGroup.Add(workerCount)
	for index := 0; index < workerCount; index++ {
		go func() {
			defer workerWaitGroup.Done()
			e.worker(searchCtx, memoryChannel, resultChannel)
		}()
	}

	// Start producer goroutine
	var producerWaitGroup sync.WaitGroup
	producerWaitGroup.Add(1)
	go func() {
		defer producerWaitGroup.Done()
		defer close(memoryChannel) // Close channel when producer is done
		err := e.findMemory(searchCtx, uint32(proc.PID), memoryChannel)
		if err != nil {
			log.Err(err).Msg("Failed to read memory")
		}
	}()

	// Wait for producer and consumers to complete
	go func() {
		producerWaitGroup.Wait()
		workerWaitGroup.Wait()
		close(resultChannel)
	}()

	// Wait for result
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case result, ok := <-resultChannel:
		if ok && result != "" {
			return result, nil
		}
	}

	return "", errors.ErrNoValidKey
}

// findMemory searches for memory regions using Glance
func (e *V4Extractor) findMemory(ctx context.Context, pid uint32, memoryChannel chan<- []byte) error {
	// Initialize a Glance instance to read process memory
	g := glance.NewGlance(pid)

	// Read memory data
	memory, err := g.Read()
	if err != nil {
		return err
	}

	totalSize := len(memory)
	log.Debug().Msgf("Read memory region, size: %d bytes", totalSize)

	// If memory is small enough, process it as a single chunk
	if totalSize <= MinChunkSize {
		select {
		case memoryChannel <- memory:
			log.Debug().Msg("Memory sent as a single chunk for analysis")
		case <-ctx.Done():
			return ctx.Err()
		}
		return nil
	}

	chunkCount := MaxWorkers * ChunkMultiplier

	// Calculate chunk size based on fixed chunk count
	chunkSize := totalSize / chunkCount
	if chunkSize < MinChunkSize {
		// Reduce number of chunks if each would be too small
		chunkCount = totalSize / MinChunkSize
		if chunkCount == 0 {
			chunkCount = 1
		}
		chunkSize = totalSize / chunkCount
	}

	// Process memory in chunks from end to beginning
	for i := chunkCount - 1; i >= 0; i-- {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Calculate start and end positions for this chunk
			start := i * chunkSize
			end := (i + 1) * chunkSize

			// Ensure the last chunk includes all remaining memory
			if i == chunkCount-1 {
				end = totalSize
			}

			// Add overlap area to catch patterns at chunk boundaries
			if i > 0 {
				start -= ChunkOverlapBytes
				if start < 0 {
					start = 0
				}
			}

			chunk := memory[start:end]

			log.Debug().
				Int("chunk_index", i+1).
				Int("total_chunks", chunkCount).
				Int("chunk_size", len(chunk)).
				Int("start_offset", start).
				Int("end_offset", end).
				Msg("Processing memory chunk")

			select {
			case memoryChannel <- chunk:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return nil
}

// worker processes memory regions to find V4 version key
func (e *V4Extractor) worker(ctx context.Context, memoryChannel <-chan []byte, resultChannel chan<- string) {
	for {
		select {
		case <-ctx.Done():
			return
		case memory, ok := <-memoryChannel:
			if !ok {
				return
			}

			if key, ok := e.SearchKey(ctx, memory); ok {
				select {
				case resultChannel <- key:
				default:
				}
			}
		}
	}
}

func (e *V4Extractor) SearchKey(ctx context.Context, memory []byte) (string, bool) {
	for _, keyPattern := range e.keyPatterns {
		index := len(memory)

		for {
			select {
			case <-ctx.Done():
				return "", false
			default:
			}

			// Find pattern from end to beginning
			index = bytes.LastIndex(memory[:index], keyPattern.Pattern)
			if index == -1 {
				break // No more matches found
			}

			// Try each offset for this pattern
			for _, offset := range keyPattern.Offsets {
				// Check if we have enough space for the key
				keyOffset := index + offset
				if keyOffset < 0 || keyOffset+32 > len(memory) {
					continue
				}

				// Extract the key data, which is at the offset position and 32 bytes long
				keyData := memory[keyOffset : keyOffset+32]

				// Validate key against database header
				if keyData, ok := e.validate(ctx, keyData); ok {
					log.Debug().
						Str("pattern", hex.EncodeToString(keyPattern.Pattern)).
						Int("offset", offset).
						Str("key", hex.EncodeToString(keyData)).
						Msg("Key found")
					return hex.EncodeToString(keyData), true
				}
			}

			index -= 1
		}
	}

	return "", false
}

func (e *V4Extractor) validate(ctx context.Context, keyDate []byte) ([]byte, bool) {
	if e.validator.Validate(keyDate) {
		return keyDate, true
	}
	// Try to find a valid key by ***
	return nil, false
}

func (e *V4Extractor) SetValidate(validator *decrypt.Validator) {
	e.validator = validator
}

type KeyPatternInfo struct {
	Pattern []byte
	Offsets []int
}
