/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 16:42
 */
package tencent

const (
    PlaceSearchTypeRegion    = "region"    // 区域搜索类型-地区
    PlaceSearchTypeNearby    = "nearby"    // 区域搜索类型-圆形区域
    PlaceSearchTypeRectangle = "rectangle" // 区域搜索类型-矩形区域
)

var (
    placeSearchTypes map[string]int
)

func init() {
    placeSearchTypes = make(map[string]int)
    placeSearchTypes[PlaceSearchTypeRegion] = 1
    placeSearchTypes[PlaceSearchTypeNearby] = 1
    placeSearchTypes[PlaceSearchTypeRectangle] = 1
}
