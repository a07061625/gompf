/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/5 0005
 * Time: 12:33
 */
package account

const (
    MaterialTypeImage = "image"
    MaterialTypeVoice = "voice"
    MaterialTypeVideo = "video"
    MaterialTypeThumb = "thumb"

    MessageTypeMpNews = "mpnews"
    MessageTypeText   = "text"
    MessageTypeVoice  = "voice"
    MessageTypeMusic  = "music"
    MessageTypeImage  = "image"
    MessageTypeVideo  = "video"
    MessageTypeWxCard = "wxcard"
)

var (
    MaterialTypes map[string]string
    MessageTypes  map[string]string
)

func init() {
    MaterialTypes = make(map[string]string)
    MaterialTypes[MaterialTypeImage] = "图片"
    MaterialTypes[MaterialTypeVoice] = "语音"
    MaterialTypes[MaterialTypeVideo] = "视频"
    MaterialTypes[MaterialTypeThumb] = "缩略图"

    MessageTypes = make(map[string]string)
    MessageTypes[MessageTypeMpNews] = "图文"
    MessageTypes[MessageTypeText] = "文本"
    MessageTypes[MessageTypeVoice] = "语音"
    MessageTypes[MessageTypeMusic] = "音乐"
    MessageTypes[MessageTypeImage] = "图片"
    MessageTypes[MessageTypeVideo] = "视频"
    MessageTypes[MessageTypeWxCard] = "卡券"
}
