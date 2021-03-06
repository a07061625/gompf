// Package mpdb db_mysql
// User: 姜伟
// Time: 2020-02-19 06:22:40
package mpdb

import (
    "sync"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"

    // 只初始化
    _ "github.com/go-sql-driver/mysql"
    "xorm.io/xorm"
)

type dbMysql struct {
    conn        *xorm.Engine
    connTime    int64
    dbName      string
    refreshTime int64
    idleTime    int64
}

func (database *dbMysql) connect() {
    conf := mpf.NewConfig().GetConfig("db")
    database.dbName = conf.GetString("mysql." + mpf.EnvProjectKey() + ".dbname")
    dsn := conf.GetString("mysql."+mpf.EnvProjectKey()+".dsn.prefix") + database.dbName + conf.GetString("mysql."+mpf.EnvProjectKey()+".dsn.params")
    conn, err := xorm.NewEngine("mysql", dsn)
    if err != nil {
        panic(mperr.NewDbMysql(errorcode.DbMysqlConnect, "mysql连接失败", err))
    }
    conn.SetMaxIdleConns(conf.GetInt("mysql." + mpf.EnvProjectKey() + ".conn.idle"))
    conn.SetMaxOpenConns(conf.GetInt("mysql." + mpf.EnvProjectKey() + ".conn.open"))
    lifeTime := conf.GetInt("mysql." + mpf.EnvProjectKey() + ".conn.life")
    conn.SetConnMaxLifetime(time.Duration(lifeTime) * time.Second)

    database.conn = conn
    database.connTime = time.Now().Unix()
    database.idleTime = conf.GetInt64("mysql." + mpf.EnvProjectKey() + ".idle")
    database.refreshTime = database.connTime + database.idleTime
}

func (database *dbMysql) Reconnect() {
    nowTime := time.Now().Unix()
    if database.conn == nil {
        database.connect()
    } else if database.refreshTime < nowTime {
        pingErr := database.conn.Ping()
        if pingErr != nil {
            database.conn.Close()
            database.connect()
        }
        database.refreshTime = nowTime + database.idleTime
    }
}

func (database *dbMysql) GetDbName() string {
    return database.dbName
}

func (database *dbMysql) GetConn() *xorm.Engine {
    return database.conn
}

var (
    onceMysql sync.Once
    insMysql  *dbMysql
)

func init() {
    insMysql = &dbMysql{nil, 0, "", 0, 0}
}

// NewDbMysql NewDbMysql
func NewDbMysql() *dbMysql {
    onceMysql.Do(func() {
        insMysql.connect()
    })

    return insMysql
}
