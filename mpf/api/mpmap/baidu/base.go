/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 14:56
 */
package baidu

var (
    coordTranslateTypes   map[string]int
    placeSearchCoordTypes map[string]int
    placeSearchTypes      map[string]int
)

func init() {
    coordTranslateTypes = make(map[string]int)
    coordTranslateTypes["1"] = 1 // 坐标类型-GPS角度
    coordTranslateTypes["2"] = 1 // 坐标类型-GPS米制
    coordTranslateTypes["3"] = 1 // 坐标类型-google
    coordTranslateTypes["4"] = 1 // 坐标类型-google米制
    coordTranslateTypes["5"] = 1 // 坐标类型-百度
    coordTranslateTypes["6"] = 1 // 坐标类型-百度米制
    coordTranslateTypes["7"] = 1 // 坐标类型-mapbar
    coordTranslateTypes["8"] = 1 // 坐标类型-51

    placeSearchCoordTypes = make(map[string]int)
    placeSearchCoordTypes["1"] = 1 // 坐标类型-GPS
    placeSearchCoordTypes["2"] = 1 // 坐标类型-国测局
    placeSearchCoordTypes["3"] = 1 // 坐标类型-百度
    placeSearchCoordTypes["4"] = 1 // 坐标类型-百度墨卡托米制

    placeSearchTypes = make(map[string]int)
    placeSearchTypes["region"] = 1    // 区域搜索类型-地区
    placeSearchTypes["nearby"] = 1    // 区域搜索类型-圆形区域
    placeSearchTypes["rectangle"] = 1 // 区域搜索类型-矩形区域
}
