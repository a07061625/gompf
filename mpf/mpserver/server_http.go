/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/6 0006
 * Time: 10:17
 */
package mpserver

import "log"

type serverHttp struct {
    serverWeb
}

func (s *serverHttp) initServer() {
    s.initBase()
}

func NewServerHttp() *serverHttp {
    s := &serverHttp{newServerWeb()}
    s.initServer()
    if s.Runner == nil {
        log.Println("http")
    }
    return s
}
