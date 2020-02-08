/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/8 0008
 * Time: 8:36
 */
package frontend

type indexController struct {
    common
}

func NewIndex() *indexController {
    return &indexController{newCommon()}
}
