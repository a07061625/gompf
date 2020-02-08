/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/8 0008
 * Time: 8:34
 */
package index

type indexController struct {
    common
}

func NewIndex() *indexController {
    return &indexController{newCommon()}
}
