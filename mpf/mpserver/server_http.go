/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/6 0006
 * Time: 10:17
 */
package mpserver

type serverHttp struct {
    serverWeb
}

func (s *serverHttp) StartServer() {
    s.baseStart()
}

func NewServerHttp() *serverHttp {
    return &serverHttp{newServerWeb()}
}
