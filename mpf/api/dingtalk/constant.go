package dingtalk

const (
    UrlService = "https://oapi.dingtalk.com"

    AccessTokenTypeCorp              = "corp"
    AccessTokenTypeProviderAuthorize = "provider_authorize"

    MessageTypeText       = "text"
    MessageTypeImage      = "image"
    MessageTypeVoice      = "voice"
    MessageTypeFile       = "file"
    MessageTypeLink       = "link"
    MessageTypeOA         = "oa"
    MessageTypeMarkdown   = "markdown"
    MessageTypeActionCard = "action_card"

    MediaTypeImage = "image"
    MediaTypeVoice = "voice"
    MediaTypeFile  = "file"
)

var (
    MessageTypes map[string]string
    MediaTypes   map[string]string
)

func init() {
    MessageTypes = make(map[string]string)
    MessageTypes[MessageTypeText] = "文本"
    MessageTypes[MessageTypeImage] = "图片"
    MessageTypes[MessageTypeVoice] = "语音"
    MessageTypes[MessageTypeFile] = "文件"
    MessageTypes[MessageTypeLink] = "链接"
    MessageTypes[MessageTypeOA] = "OA"
    MessageTypes[MessageTypeMarkdown] = "markdown"
    MessageTypes[MessageTypeActionCard] = "卡片"

    MediaTypes = make(map[string]string)
    MediaTypes[MediaTypeImage] = "图片"
    MediaTypes[MediaTypeVoice] = "语音"
    MediaTypes[MediaTypeFile] = "文件"
}
