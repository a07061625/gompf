package mpserver

import (
    "net"
    "time"
)

func (this *listener) Accept() (conn net.Conn, err error) {
    tc, err := this.Listener.(*net.TCPListener).AcceptTCP()
    if err != nil {
        return
    }
    tc.SetKeepAlive(true)
    tc.SetKeepAlivePeriod(3 * time.Minute)
    return tc, nil
}

//获取sock文件句柄
func (this *listener) File() (uintptr, error) {
    f, err := this.Listener.(*net.TCPListener).File()
    if err != nil {
        return 0, err
    }
    return f.Fd(), nil
}
