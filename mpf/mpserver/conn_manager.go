package mpserver

import (
    "net"
    "sync"
    "time"
)

type connManager struct {
    sync.WaitGroup
    Counter      int
    mux          sync.Mutex
    idleConnList map[string]net.Conn
}

func (cm *connManager) Add(delta int) {
    cm.Counter += delta
    cm.WaitGroup.Add(delta)
}

func (cm *connManager) Done() {
    cm.Counter--
    cm.WaitGroup.Done()
}

func (cm *connManager) Close(t time.Duration) {
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

func (cm *connManager) RemoveIdleConn(key string) {
    cm.mux.Lock()
    delete(cm.idleConnList, key)
    cm.mux.Unlock()
}

func (cm *connManager) AddIdleConn(key string, conn net.Conn) {
    cm.mux.Lock()
    cm.idleConnList[key] = conn
    cm.mux.Unlock()
}

func newConnManager() *connManager {
    cm := &connManager{}
    cm.WaitGroup = sync.WaitGroup{}
    cm.Counter = 0
    cm.idleConnList = make(map[string]net.Conn)
    return cm
}
