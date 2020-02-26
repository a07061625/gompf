// Package mpmodule request_rpc
// User: 姜伟
// Time: 2020-02-19 05:03:42
package mpmodule

import (
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpcache"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/frame"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
)

// ModuleRPC rpc模块
type ModuleRPC struct {
    ModuleBasic
    fuseState                string
    fuseKey                  string
    fuseTimeCloseKeep        time.Duration
    fuseTimeOpenKeep         time.Duration
    fuseReqNumFail           int
    fuseReqNumFailNow        int
    fuseReqNumHalfSuccess    int
    fuseReqNumHalfSuccessNow int
}

// checkFuseState 检测熔断器状态
// @return uint 0:正常 >0:不正常
func (m *ModuleRPC) checkFuseState() uint {
    if m.fuseState == frame.FuseStateOpen {
        _, ok := mpcache.NewLocal().Get(m.fuseKey)
        if ok {
            return errorcode.CommonBaseFuse
        }
        m.fuseState = frame.FuseStateHalfOpen
        m.fuseReqNumFailNow = 0
        m.fuseReqNumHalfSuccessNow = 0
    }
    return 0
}

// refreshFuseState 更新熔断器状态
func (m *ModuleRPC) refreshFuseState(reqResult bool) {
    if reqResult {
        if m.fuseState == frame.FuseStateHalfOpen {
            m.fuseReqNumHalfSuccessNow++
            if m.fuseReqNumHalfSuccessNow >= m.fuseReqNumHalfSuccess {
                m.fuseState = frame.FuseStateClosed
                m.fuseReqNumFailNow = 0
                m.fuseReqNumHalfSuccessNow = 0
            }
        }
    } else if m.fuseState == frame.FuseStateClosed {
        _, ok := mpcache.NewLocal().Get(m.fuseKey)
        if !ok {
            mpcache.NewLocal().Set(m.fuseKey, 1, m.fuseTimeCloseKeep)
            m.fuseReqNumFailNow = 0
            m.fuseReqNumHalfSuccessNow = 0
        }
        m.fuseReqNumFailNow++
        if m.fuseReqNumFailNow >= m.fuseReqNumFail {
            m.fuseState = frame.FuseStateOpen
            m.fuseReqNumFailNow = 0
            m.fuseReqNumHalfSuccessNow = 0
        }
    } else if m.fuseState == frame.FuseStateHalfOpen {
        mpcache.NewLocal().Set(m.fuseKey, 1, m.fuseTimeOpenKeep)
        m.fuseState = frame.FuseStateOpen
        m.fuseReqNumFailNow = 0
        m.fuseReqNumHalfSuccessNow = 0
    }
}

// NewRPC rpc模块
func NewRPC(tag string) ModuleRPC {
    conf := mpf.NewConfig().GetConfig("server")
    confPrefix := mpf.EnvType() + "." + mpf.EnvProjectKeyModule() + "."
    m := ModuleRPC{NewBasic(tag), "", "", 0, 0, 0, 0, 0, 0}
    m.fuseState = frame.FuseStateClosed
    m.fuseKey = mpf.HashCrc32(project.LocalPrefixFuse+m.tag, "")
    m.fuseTimeCloseKeep = time.Duration(conf.GetInt(confPrefix+"time.closekeep")) * time.Second
    m.fuseTimeOpenKeep = time.Duration(conf.GetInt(confPrefix+"time.openkeep")) * time.Second
    m.fuseReqNumFail = conf.GetInt(confPrefix + "reqnum.fail")
    m.fuseReqNumHalfSuccess = conf.GetInt(confPrefix + "reqnum.halfsuccess")
    return m
}
