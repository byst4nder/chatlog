package model

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

type MediaMsg struct {
	XMLName xml.Name `xml:"msg"`
	Image   Image    `xml:"img,omitempty"`
	Video   Video    `xml:"videomsg,omitempty"`
	App     App      `xml:"appmsg,omitempty"`
}

type Image struct {
	MD5 string `xml:"md5,attr"`
	// HdLength            string `xml:"hdlength,attr"`
	// Length              string `xml:"length,attr"`
	// AesKey              string `xml:"aeskey,attr"`
	// EncryVer            string `xml:"encryver,attr"`
	// OriginSourceMd5     string `xml:"originsourcemd5,attr"`
	// FileKey             string `xml:"filekey,attr"`
	// UploadContinueCount string `xml:"uploadcontinuecount,attr"`
	// ImgSourceUrl        string `xml:"imgsourceurl,attr"`
	// HevcMidSize         string `xml:"hevc_mid_size,attr"`
	// CdnBigImgUrl        string `xml:"cdnbigimgurl,attr"`
	// CdnMidImgUrl        string `xml:"cdnmidimgurl,attr"`
	// CdnThumbUrl         string `xml:"cdnthumburl,attr"`
	// CdnThumbLength      string `xml:"cdnthumblength,attr"`
	// CdnThumbWidth       string `xml:"cdnthumbwidth,attr"`
	// CdnThumbHeight      string `xml:"cdnthumbheight,attr"`
	// CdnThumbAesKey      string `xml:"cdnthumbaeskey,attr"`
}

type Video struct {
	Md5    string `xml:"md5,attr"`
	RawMd5 string `xml:"rawmd5,attr"`
	// Length            string `xml:"length,attr"`
	// PlayLength        string `xml:"playlength,attr"`
	// Offset            string `xml:"offset,attr"`
	// FromUserName      string `xml:"fromusername,attr"`
	// Status            string `xml:"status,attr"`
	// Compress          string `xml:"compress,attr"`
	// CameraType        string `xml:"cameratype,attr"`
	// Source            string `xml:"source,attr"`
	// AesKey            string `xml:"aeskey,attr"`
	// CdnVideoUrl       string `xml:"cdnvideourl,attr"`
	// CdnThumbUrl       string `xml:"cdnthumburl,attr"`
	// CdnThumbLength    string `xml:"cdnthumblength,attr"`
	// CdnThumbWidth     string `xml:"cdnthumbwidth,attr"`
	// CdnThumbHeight    string `xml:"cdnthumbheight,attr"`
	// CdnThumbAesKey    string `xml:"cdnthumbaeskey,attr"`
	// EncryVer          string `xml:"encryver,attr"`
	// RawLength         string `xml:"rawlength,attr"`
	// CdnRawVideoUrl    string `xml:"cdnrawvideourl,attr"`
	// CdnRawVideoAesKey string `xml:"cdnrawvideoaeskey,attr"`
}

type App struct {
	Type              int         `xml:"type"`
	Title             string      `xml:"title"`
	Des               string      `xml:"des"`
	URL               string      `xml:"url"`                         // type 5 分享
	AppAttach         *AppAttach  `xml:"appattach,omitempty"`         // type 6 文件
	MD5               string      `xml:"md5,omitempty"`               // type 6 文件
	RecordItem        *RecordItem `xml:"recorditem,omitempty"`        // type 19 合并转发
	SourceDisplayName string      `xml:"sourcedisplayname,omitempty"` // type 33 小程序
	FinderFeed        *FinderFeed `xml:"finderFeed,omitempty"`        // type 51 视频号
	ReferMsg          *ReferMsg   `xml:"refermsg,omitempty"`          // type 57 引用
	PatMsg            *PatMsg     `xml:"patMsg,omitempty"`            // type 62 拍一拍
	WCPayInfo         *WCPayInfo  `xml:"wcpayinfo,omitempty"`         // type 2000 微信转账
}

// ReferMsg 表示引用消息
type ReferMsg struct {
	Type        int64  `xml:"type"`
	SvrID       string `xml:"svrid"`
	FromUsr     string `xml:"fromusr"`
	ChatUsr     string `xml:"chatusr"`
	DisplayName string `xml:"displayname"`
	MsgSource   string `xml:"msgsource"`
	Content     string `xml:"content"`
	StrID       string `xml:"strid"`
	CreateTime  int64  `xml:"createtime"`
}

// AppAttach 表示应用附件
type AppAttach struct {
	TotalLen       string `xml:"totallen"`
	AttachID       string `xml:"attachid"`
	CDNAttachURL   string `xml:"cdnattachurl"`
	EmoticonMD5    string `xml:"emoticonmd5"`
	AESKey         string `xml:"aeskey"`
	FileExt        string `xml:"fileext"`
	IsLargeFileMsg string `xml:"islargefilemsg"`
}

type RecordItem struct {
	CDATA string `xml:",cdata"`

	// 解析后的记录信息
	RecordInfo *RecordInfo
}

// RecordInfo 表示聊天记录信息
type RecordInfo struct {
	XMLName       xml.Name `xml:"recordinfo"`
	FromScene     string   `xml:"fromscene,omitempty"`
	FavUsername   string   `xml:"favusername,omitempty"`
	FavCreateTime string   `xml:"favcreatetime,omitempty"`
	IsChatRoom    string   `xml:"isChatRoom,omitempty"`
	Title         string   `xml:"title,omitempty"`
	Desc          string   `xml:"desc,omitempty"`
	Info          string   `xml:"info,omitempty"`
	DataList      DataList `xml:"datalist,omitempty"`
}

// DataList 表示数据列表
type DataList struct {
	Count     string     `xml:"count,attr,omitempty"`
	DataItems []DataItem `xml:"dataitem,omitempty"`
}

// DataItem 表示数据项
type DataItem struct {
	DataType      string `xml:"datatype,attr,omitempty"`
	DataID        string `xml:"dataid,attr,omitempty"`
	HTMLID        string `xml:"htmlid,attr,omitempty"`
	DataFmt       string `xml:"datafmt,omitempty"`
	SourceName    string `xml:"sourcename,omitempty"`
	SourceTime    string `xml:"sourcetime,omitempty"`
	SourceHeadURL string `xml:"sourceheadurl,omitempty"`
	DataDesc      string `xml:"datadesc,omitempty"`

	// 图片特有字段
	ThumbSourcePath  string `xml:"thumbsourcepath,omitempty"`
	ThumbSize        string `xml:"thumbsize,omitempty"`
	CDNDataURL       string `xml:"cdndataurl,omitempty"`
	CDNDataKey       string `xml:"cdndatakey,omitempty"`
	CDNThumbURL      string `xml:"cdnthumburl,omitempty"`
	CDNThumbKey      string `xml:"cdnthumbkey,omitempty"`
	DataSourcePath   string `xml:"datasourcepath,omitempty"`
	FullMD5          string `xml:"fullmd5,omitempty"`
	ThumbFullMD5     string `xml:"thumbfullmd5,omitempty"`
	ThumbHead256MD5  string `xml:"thumbhead256md5,omitempty"`
	DataSize         string `xml:"datasize,omitempty"`
	CDNEncryVer      string `xml:"cdnencryver,omitempty"`
	SrcChatname      string `xml:"srcChatname,omitempty"`
	SrcMsgLocalID    string `xml:"srcMsgLocalid,omitempty"`
	SrcMsgCreateTime string `xml:"srcMsgCreateTime,omitempty"`
	MessageUUID      string `xml:"messageuuid,omitempty"`
	FromNewMsgID     string `xml:"fromnewmsgid,omitempty"`

	// 套娃合并转发
	DataTitle string     `xml:"datatitle,omitempty"`
	RecordXML *RecordXML `xml:"recordxml,omitempty"`
}

type RecordXML struct {
	RecordInfo RecordInfo `xml:"recordinfo,omitempty"`
}

func (r *RecordInfo) String(title, host string) string {
	buf := strings.Builder{}
	if title == "" {
		title = r.Title
	}
	buf.WriteString(fmt.Sprintf("[合并转发|%s]\n", title))
	for _, item := range r.DataList.DataItems {
		buf.WriteString(fmt.Sprintf("  %s %s\n", item.SourceName, item.SourceTime))

		// 套娃合并转发
		if item.DataType == "17" && item.RecordXML != nil {
			content := item.RecordXML.RecordInfo.String(item.DataTitle, host)
			if content != "" {
				for _, line := range strings.Split(content, "\n") {
					buf.WriteString(fmt.Sprintf("  %s\n", line))
				}
			}
			continue
		}

		switch item.DataFmt {
		case "pic", "jpg":
			buf.WriteString(fmt.Sprintf("  ![图片](http://%s/image/%s)\n", host, item.FullMD5))
		default:
			for _, line := range strings.Split(item.DataDesc, "\n") {
				buf.WriteString(fmt.Sprintf("  %s\n", line))
			}
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

// PatMsg 拍一拍消息结构
type PatMsg struct {
	ChatUser  string  `xml:"chatUser"`  // 被拍的用户
	RecordNum int     `xml:"recordNum"` // 记录数量
	Records   Records `xml:"records"`   // 拍一拍记录
}

// Records 拍一拍记录集合
type Records struct {
	Record []PatRecord `xml:"record"` // 拍一拍记录列表
}

// PatRecord 单条拍一拍记录
type PatRecord struct {
	FromUser   string `xml:"fromUser"`   // 发起拍一拍的用户
	PattedUser string `xml:"pattedUser"` // 被拍的用户
	Templete   string `xml:"templete"`   // 模板文本
	CreateTime int64  `xml:"createTime"` // 创建时间
	SvrId      string `xml:"svrId"`      // 服务器ID
	ReadStatus int    `xml:"readStatus"` // 已读状态
}

// WCPayInfo 微信支付信息
type WCPayInfo struct {
	PaySubType        int    `xml:"paysubtype"`        // 支付子类型
	FeeDesc           string `xml:"feedesc"`           // 金额描述，如"￥200000.00"
	TranscationID     string `xml:"transcationid"`     // 交易ID
	TransferID        string `xml:"transferid"`        // 转账ID
	InvalidTime       string `xml:"invalidtime"`       // 失效时间
	BeginTransferTime string `xml:"begintransfertime"` // 开始转账时间
	EffectiveDate     string `xml:"effectivedate"`     // 生效日期
	PayMemo           string `xml:"pay_memo"`          // 支付备注
	ReceiverUsername  string `xml:"receiver_username"` // 接收方用户名
	PayerUsername     string `xml:"payer_username"`    // 支付方用户名
}

// FinderFeed 视频号信息
type FinderFeed struct {
	ObjectID            string          `xml:"objectId"`
	FeedType            string          `xml:"feedType"`
	Nickname            string          `xml:"nickname"`
	Avatar              string          `xml:"avatar"`
	Desc                string          `xml:"desc"`
	MediaCount          string          `xml:"mediaCount"`
	ObjectNonceID       string          `xml:"objectNonceId"`
	LiveID              string          `xml:"liveId"`
	Username            string          `xml:"username"`
	AuthIconURL         string          `xml:"authIconUrl"`
	AuthIconType        int             `xml:"authIconType"`
	ContactJumpInfoStr  string          `xml:"contactJumpInfoStr"`
	SourceCommentScene  int             `xml:"sourceCommentScene"`
	MediaList           FinderMediaList `xml:"mediaList"`
	MegaVideo           FinderMegaVideo `xml:"megaVideo"`
	BizUsername         string          `xml:"bizUsername"`
	BizNickname         string          `xml:"bizNickname"`
	BizAvatar           string          `xml:"bizAvatar"`
	BizUsernameV2       string          `xml:"bizUsernameV2"`
	BizAuthIconURL      string          `xml:"bizAuthIconUrl"`
	BizAuthIconType     int             `xml:"bizAuthIconType"`
	EcSource            string          `xml:"ecSource"`
	LastGMsgID          string          `xml:"lastGMsgID"`
	ShareBypData        string          `xml:"shareBypData"`
	IsDebug             int             `xml:"isDebug"`
	ContentType         int             `xml:"content_type"`
	FinderForwardSource string          `xml:"finderForwardSource"`
}

type FinderMediaList struct {
	Media []FinderMedia `xml:"media"`
}

type FinderMedia struct {
	ThumbURL          string `xml:"thumbUrl"`
	FullCoverURL      string `xml:"fullCoverUrl"`
	VideoPlayDuration string `xml:"videoPlayDuration"`
	URL               string `xml:"url"`
	CoverURL          string `xml:"coverUrl"`
	Height            string `xml:"height"`
	MediaType         string `xml:"mediaType"`
	FullClipInset     string `xml:"fullClipInset"`
	Width             string `xml:"width"`
}

type FinderMegaVideo struct {
	ObjectID      string `xml:"objectId"`
	ObjectNonceID string `xml:"objectNonceId"`
}

type SysMsg struct {
	Type              string             `xml:"type,attr"`
	DelChatRoomMember *DelChatRoomMember `xml:"delchatroommember,omitempty"`
	SysMsgTemplate    *SysMsgTemplate    `xml:"sysmsgtemplate,omitempty"`
}

// 第一种消息类型：删除群成员/二维码邀请
type DelChatRoomMember struct {
	Plain string `xml:"plain"`
	Text  string `xml:"text"`
	Link  QRLink `xml:"link"`
}

type QRLink struct {
	Scene      string       `xml:"scene"`
	Text       string       `xml:"text"`
	MemberList QRMemberList `xml:"memberlist"`
	QRCode     string       `xml:"qrcode"`
}

type QRMemberList struct {
	Usernames []UsernameItem `xml:"username"`
}

type UsernameItem struct {
	Value string `xml:",chardata"`
}

// 第二种消息类型：系统消息模板
type SysMsgTemplate struct {
	ContentTemplate ContentTemplate `xml:"content_template"`
}

type ContentTemplate struct {
	Type     string   `xml:"type,attr"`
	Plain    string   `xml:"plain"`
	Template string   `xml:"template"`
	LinkList LinkList `xml:"link_list"`
}

type LinkList struct {
	Links []Link `xml:"link"`
}

type Link struct {
	Name       string     `xml:"name,attr"`
	Type       string     `xml:"type,attr"`
	MemberList MemberList `xml:"memberlist"`
	Separator  string     `xml:"separator,omitempty"`
	Title      string     `xml:"title,omitempty"`
}

type MemberList struct {
	Members []Member `xml:"member"`
}

type Member struct {
	Username string `xml:"username"`
	Nickname string `xml:"nickname"`
}

func (s *SysMsg) String() string {
	if s.Type == "delchatroommember" {
		return s.DelChatRoomMemberString()
	}
	return s.SysMsgTemplateString()
}

func (s *SysMsg) DelChatRoomMemberString() string {
	if s.DelChatRoomMember == nil {
		return ""
	}
	return s.DelChatRoomMember.Plain
}

func (s *SysMsg) SysMsgTemplateString() string {
	if s.SysMsgTemplate == nil {
		return ""
	}

	template := s.SysMsgTemplate.ContentTemplate.Template
	links := s.SysMsgTemplate.ContentTemplate.LinkList.Links

	// 创建一个映射，用于存储占位符名称和对应的替换内容
	replacements := make(map[string]string)

	// 遍历所有链接，为每个占位符准备替换内容
	for _, link := range links {
		var replacement string

		// 根据链接类型和成员信息生成替换内容
		switch link.Type {
		case "link_profile":
			// 使用自定义分隔符，如果未指定则默认使用"、"
			separator := link.Separator
			if separator == "" {
				separator = "、"
			}

			// 处理成员信息，格式为 nickname(username)
			var memberTexts []string
			for _, member := range link.MemberList.Members {
				if member.Nickname != "" {
					memberText := member.Nickname
					if member.Username != "" {
						memberText += "(" + member.Username + ")"
					}
					memberTexts = append(memberTexts, memberText)
				}
			}

			// 使用指定的分隔符连接所有成员文本
			replacement = strings.Join(memberTexts, separator)

		// 可以根据需要添加其他链接类型的处理逻辑
		default:
			if link.Title != "" {
				replacement = link.Title
			} else {
				replacement = ""
			}
		}

		// 将占位符名称和替换内容存入映射
		replacements["$"+link.Name+"$"] = replacement
	}

	// 使用正则表达式查找并替换所有占位符
	re := regexp.MustCompile(`\$([^$]+)\$`)
	result := re.ReplaceAllStringFunc(template, func(match string) string {
		if replacement, ok := replacements[match]; ok {
			return replacement
		}
		// 如果找不到对应的替换内容，保留原占位符
		return match
	})

	return result
}
