/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 15:15
 */
package mpiot

const (
    BaiDuDomainBJ     = "iot.bj.baidubce.com"
    BaiDuDomainBJMqtt = "api.mqtt.iot.bj.baidubce.com"
    BaiDuDomainGZ     = "iot.gz.baidubce.com"
    BaiDuDomainGZMqtt = "api.mqtt.iot.gz.baidubce.com"
)

var (
    BaiDuDomains map[string]int
)

func init() {
    BaiDuDomains = make(map[string]int)
    BaiDuDomains[BaiDuDomainBJ] = 1
    BaiDuDomains[BaiDuDomainBJMqtt] = 1
    BaiDuDomains[BaiDuDomainGZ] = 1
    BaiDuDomains[BaiDuDomainGZMqtt] = 1
}
