/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/6 0006
 * Time: 10:17
 */
package mpserver

import (
    "github.com/a07061625/gompf/mpf"
)

type serverHttp struct {
    serverSimple
}

func (s *serverHttp) bootServer() {
    conf := mpf.NewConfig().GetConfig("server")
    s.bootSimple(conf)
}

func (s *serverHttp) StartServer() {
    s.bootServer()
    s.startSimple()
}

func NewServerHttp() *serverHttp {
    return &serverHttp{newServerSimple()}
}
