/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/25 0025
 * Time: 23:23
 */
package mpmap

const (
    // 百度
    BaiDuCheckTypeServerIp = "server-ip" // 校验类型-服务端ip
    BaiDuCheckTypeServerSn = "server-sn" // 校验类型-服务端签名
    BaiDuCheckTypeBrowse   = "browse"    // 校验类型-浏览器

    // 腾讯
    TencentGetTypeServer = "server" // 获取类型-服务端
    TencentGetTypeMobile = "mobile" // 获取类型-移动端
    TencentGetTypeBrowse = "browse" // 获取类型-网页端
)

var (
    BdCheckTypes    map[string]int
    TencentGetTypes map[string]int
)

func init() {
    BdCheckTypes = make(map[string]int)
    BdCheckTypes[BaiDuCheckTypeServerIp] = 1
    BdCheckTypes[BaiDuCheckTypeServerSn] = 1
    BdCheckTypes[BaiDuCheckTypeBrowse] = 1

    TencentGetTypes = make(map[string]int)
    TencentGetTypes[TencentGetTypeServer] = 1
    TencentGetTypes[TencentGetTypeMobile] = 1
    TencentGetTypes[TencentGetTypeBrowse] = 1
}
