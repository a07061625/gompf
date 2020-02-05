/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/27 0027
 * Time: 19:47
 */
package mpprint

type utilPrint struct {
}

var (
    insUtil *utilPrint
)

func init() {
    insUtil = &utilPrint{}
}

func NewUtilPrint() *utilPrint {
    return insUtil
}
