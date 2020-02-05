/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 16:42
 */
package gaode

var (
    inputTipDataTypes   map[string]int
    coordConvertSysList map[string]int
)

func init() {
    inputTipDataTypes = make(map[string]int)
    inputTipDataTypes["all"] = 1     // 返回数据类型-所有数据类型
    inputTipDataTypes["poi"] = 1     // 返回数据类型-POI数据
    inputTipDataTypes["bus"] = 1     // 返回数据类型-公交站点数据
    inputTipDataTypes["busline"] = 1 // 返回数据类型-公交线路数据

    coordConvertSysList = make(map[string]int)
    coordConvertSysList["gps"] = 1
    coordConvertSysList["mapbar"] = 1
    coordConvertSysList["baidu"] = 1
    coordConvertSysList["autonavi"] = 1
}
