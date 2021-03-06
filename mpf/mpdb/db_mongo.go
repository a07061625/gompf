// Package mpdb db_mongo
// User: 姜伟
// Time: 2020-02-19 06:20:26
package mpdb

import (
    "context"
    "sync"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

type dbMonGo struct {
    conn        *mongo.Client
    connTime    int64
    refreshTime int64
    idleTime    int64
    dbName      string
    db          *mongo.Database
}

func (database *dbMonGo) connect() {
    conf := mpf.NewConfig().GetConfig("db")
    database.dbName = conf.GetString("mongo." + mpf.EnvProjectKey() + ".dbname")
    uri := conf.GetString("mongo."+mpf.EnvProjectKey()+".uris") + database.dbName + conf.GetString("mongo."+mpf.EnvProjectKey()+".params")
    ctx, cancel := context.WithTimeout(context.Background(), 31536000*time.Second)
    defer cancel()
    conn, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        panic(mperr.NewDbMonGo(errorcode.DbMonGoConnect, "mongo连接失败", err))
    }

    newCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    pingErr := conn.Ping(newCtx, readpref.Primary())
    if pingErr != nil {
        conn.Disconnect(newCtx)
        panic(mperr.NewDbMonGo(errorcode.DbMonGoConnect, "mongo连接失败", err))
    }

    database.conn = conn
    database.db = database.conn.Database(database.dbName)
    database.connTime = time.Now().Unix()
    database.idleTime = conf.GetInt64("mongo." + mpf.EnvProjectKey() + ".idle")
    database.refreshTime = database.connTime + database.idleTime
}

func (database *dbMonGo) Reconnect() {
    nowTime := time.Now().Unix()
    if database.conn == nil {
        database.connect()
    } else if database.refreshTime < nowTime {
        newCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        defer cancel()
        pingErr := database.conn.Ping(newCtx, readpref.Primary())
        if pingErr != nil {
            database.conn.Disconnect(newCtx)
            database.connect()
        }
        database.refreshTime = nowTime + database.idleTime
    }
}

func (database *dbMonGo) GetDb() *mongo.Database {
    return database.db
}

func (database *dbMonGo) GetDbName() string {
    return database.dbName
}

func (database *dbMonGo) GetConn() *mongo.Client {
    return database.conn
}

var (
    onceMonGo sync.Once
    insMonGo  *dbMonGo
)

func init() {
    insMonGo = &dbMonGo{nil, 0, 0, 0, "", nil}
}

// NewDbMonGo NewDbMonGo
func NewDbMonGo() *dbMonGo {
    onceMonGo.Do(func() {
        insMonGo.connect()
    })

    return insMonGo
}
