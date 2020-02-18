// Package mpdb container_mongo
// User: 姜伟
// Time: 2020-02-19 06:19:36
package mpdb

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpdp"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/a07061625/gompf/mpf/mplog"
)

type containerMonGo struct {
    mpdp.SubjectBasic
}

// 观察者通知,通过返回错误对象的Code属性来判断是否有错误
//   Code=0:没有错误 Code>0:有错误
func (container *containerMonGo) NotifyObservers(data interface{}) error {
    observerNum := len(container.Observers)
    if observerNum == 0 {
        return nil
    }

    errObj := mperr.NewDbMonGo(errorcode.CommonBaseSuccess, "", nil)
    defer func() {
        if r := recover(); r != nil {
            if err1, ok := r.(*mperr.ErrorCommon); ok {
                mplog.LogError("mongo observer error: " + mpf.JSONMarshal(err1))
                errObj.Code = err1.Code
                errObj.Msg = err1.Msg
            } else if err2, ok := r.(error); ok {
                mplog.LogError("mongo observer error: " + err2.Error())
                errObj.Code = errorcode.DbMonGoOperate
                errObj.Msg = err2.Error()
            } else {
                mplog.LogError("mongo observer error")
                errObj.Code = errorcode.DbMonGoOperate
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

// NewContainerMonGo NewContainerMonGo
func NewContainerMonGo() *containerMonGo {
    return &containerMonGo{mpdp.NewSubject()}
}
