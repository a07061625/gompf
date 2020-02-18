// Package mpcache redis
// User: 姜伟
// Time: 2020-02-19 06:16:30
package mpcache

import (
    "errors"
    "strconv"
    "sync"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/go-redis/redis/v7"
)

type cacheRedis struct {
    conn        *redis.Client
    db          int
    connTime    int64
    clientName  string
    refreshTime int64
    idleTime    int64
}

func (c *cacheRedis) connect() {
    conf := mpf.NewConfig().GetConfig("caches")
    c.db = conf.GetInt("redis." + mpf.EnvProjectKey() + ".db")
    c.conn = redis.NewClient(&redis.Options{
        Addr:     conf.GetString("redis." + mpf.EnvProjectKey() + ".address"),
        Password: conf.GetString("redis." + mpf.EnvProjectKey() + ".password"),
        DB:       c.db,
    })

    pingErr := c.conn.Ping().Err()
    if pingErr != nil {
        c.conn = nil
        panic(mperr.NewCacheRedis(errorcode.CacheRedisConnect, "redis连接失败", pingErr))
    }

    c.connTime = time.Now().Unix()
    c.idleTime = conf.GetInt64("redis." + mpf.EnvProjectKey() + ".idle")
    c.refreshTime = c.connTime + c.idleTime
    clientKey := strconv.FormatInt(c.connTime, 10) + mpf.ToolCreateNonceStr(8, "numlower")
    c.clientName = mpf.HashCrc32(clientKey, "")
    c.conn.Do("client", "setname", c.clientName)
}

func (c *cacheRedis) Reconnect() {
    nowTime := time.Now().Unix()
    if c.conn == nil {
        c.connect()
    } else if c.refreshTime < nowTime {
        pingErr := c.conn.Ping().Err()
        if pingErr != nil {
            c.conn.Close()
            c.connect()
        }
        c.refreshTime = nowTime + c.idleTime
    }
}

func (c *cacheRedis) GetConn() *redis.Client {
    return c.conn
}

func (c *cacheRedis) GetClientName() string {
    return c.clientName
}

func (c *cacheRedis) GetDb() int {
    return c.db
}

func (c *cacheRedis) DoHmSet(data []string) (interface{}, error) {
    if len(data) > 0 {
        dataList := make([]interface{}, 0)
        dataList = append(dataList, "hmset")
        for _, v := range data {
            dataList = append(dataList, v)
        }
        return c.conn.Do(dataList...).Result()
    }
    return nil, errors.New("数据不能为空")
}

var (
    onceRedis sync.Once
    insRedis  *cacheRedis
)

func init() {
    insRedis = &cacheRedis{nil, -1, 0, "", 0, 0}
}

// NewRedis NewRedis
func NewRedis() *cacheRedis {
    onceRedis.Do(func() {
        insRedis.connect()
    })

    return insRedis
}
