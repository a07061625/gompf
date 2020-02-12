package mpserver

import (
    "net"
    "sync"
    "time"
)

type ConnManager struct {
    sync.WaitGroup
    Counter      int
    mux          sync.Mutex
    idleConnList map[string]net.Conn
}

func (cm *ConnManager) add(delta int) {
    cm.Counter += delta
    cm.WaitGroup.Add(delta)
}

func (cm *ConnManager) done() {
    cm.Counter--
    cm.WaitGroup.Done()
}

func (cm *ConnManager) close(t time.Duration) {
    cm.mux.Lock()
    dt := time.Now().Add(t)
    for _, c := range cm.idleConnList {
        c.SetDeadline(dt)
    }
    cm.idleConnList = nil
    cm.mux.Unlock()
    cm.WaitGroup.Wait()
    return
}

func (cm *ConnManager) rmIdleConn(key string) {
    cm.mux.Lock()
    delete(cm.idleConnList, key)
    cm.mux.Unlock()
}

func (cm *ConnManager) addIdleConn(key string, conn net.Conn) {
    cm.mux.Lock()
    cm.idleConnList[key] = conn
    cm.mux.Unlock()
}

func newConnManager() *ConnManager {
    cm := &ConnManager{}
    cm.WaitGroup = sync.WaitGroup{}
    cm.idleConnList = make(map[string]net.Conn)
    return cm
}
