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

func (s *serverHttp) ReStart() {
}

func (s *serverHttp) Start() {
    s.bootstrap()
}

func (s *serverHttp) Stop() {
}
