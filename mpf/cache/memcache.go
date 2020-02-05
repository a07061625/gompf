/**
 * memcache缓存
 * User: 姜伟
 * Date: 2019/12/24 0024
 * Time: 11:53
 */
package cache

import (
    "sync"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/bradfitz/gomemcache/memcache"
)

type cacheMem struct {
    conn        *memcache.Client
    connTime    int
    refreshTime int
    idleTime    int
}

func (c *cacheMem) connect() {
    conf := mpf.NewConfig().GetConfig("caches")
    servers := conf.GetStringSlice("memcache." + mpf.EnvProjectKey() + ".servers")
    c.conn = memcache.New(servers...)
    pingErr := c.conn.Ping()
    if pingErr != nil {
        c.conn = nil
        panic(mperr.NewCacheMem(errorcode.CacheMemCacheConnect, "memcache连接失败", pingErr))
    }

    c.connTime = time.Now().Second()
    c.idleTime = conf.GetInt("memcache." + mpf.EnvProjectKey() + ".idle")
    c.refreshTime = c.connTime + c.idleTime
}

func (c *cacheMem) Reconnect() {
    nowTime := time.Now().Second()
    if c.conn == nil {
        c.connect()
    } else if c.refreshTime < nowTime {
        pingErr := c.conn.Ping()
        if pingErr != nil {
            c.connect()
        }
        c.refreshTime = nowTime + c.idleTime
    }
}

func (c *cacheMem) GetConn() *memcache.Client {
    return c.conn
}

var (
    onceMem sync.Once
    insMem  *cacheMem
)

func init() {
    insMem = &cacheMem{nil, 0, 0, 0}
}

func NewMem() *cacheMem {
    onceMem.Do(func() {
        insMem.connect()
    })

    return insMem
}
