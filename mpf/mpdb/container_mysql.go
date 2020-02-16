/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 20-2-16
 * Time: 下午4:30
 */
package mpdb

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpdp"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/a07061625/gompf/mpf/mplog"
)

type containerMysql struct {
    mpdp.SubjectBasic
}

// 观察者通知,通过返回错误对象的Code属性来判断是否有错误
//   Code=0:没有错误 Code>0:有错误
func (container *containerMysql) NotifyObservers(data interface{}) error {
    observerNum := len(container.Observers)
    if observerNum == 0 {
        return nil
    }

    errObj := mperr.NewDbMysql(errorcode.CommonBaseSuccess, "", nil)
    defer func() {
        if r := recover(); r != nil {
            if err1, ok := r.(*mperr.ErrorCommon); ok {
                mplog.LogError("mysql observer error: " + mpf.JsonMarshal(err1))
                errObj.Code = err1.Code
                errObj.Msg = err1.Msg
            } else if err2, ok := r.(error); ok {
                mplog.LogError("mysql observer error: " + err2.Error())
                errObj.Code = errorcode.DbMysqlOperate
                errObj.Msg = err2.Error()
            } else {
                mplog.LogError("mysql observer error")
                errObj.Code = errorcode.DbMysqlOperate
                errObj.Msg = "操作错误"
            }
        }
    }()

    i := 0
    for ; i < observerNum; i++ {
        container.Observers[i].Notify(data)
    }

    return errObj
}

func NewContainerMysql() *containerMysql {
    return &containerMysql{mpdp.NewSubject()}
}
