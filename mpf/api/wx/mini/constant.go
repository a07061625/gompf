/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/5 0005
 * Time: 12:33
 */
package mini

const (
    MessageCustomTypeText        = "text"
    MessageCustomTypeImage       = "image"
    MessageCustomTypeLink        = "link"
    MessageCustomTypeProgramPage = "miniprogrampage"
)

var (
    MessageCustomTypes map[string]string
)

func init() {
    MessageCustomTypes = make(map[string]string)
    MessageCustomTypes[MessageCustomTypeText] = "文本"
    MessageCustomTypes[MessageCustomTypeImage] = "图片"
    MessageCustomTypes[MessageCustomTypeLink] = "图文链接"
    MessageCustomTypes[MessageCustomTypeProgramPage] = "小程序卡片"
}
