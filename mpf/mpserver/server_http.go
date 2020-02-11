/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/11 0011
 * Time: 1:15
 */
package mpserver

type serverHttp struct {
    serverBase
}

func (s *serverHttp) init() {
    s.initBase()
}

func (s *serverHttp) Restart() {
    s.restartBase()
}

func (s *serverHttp) Start() {
    s.startBase()
}

func (s *serverHttp) Stop() {
    s.stopBase()
}
