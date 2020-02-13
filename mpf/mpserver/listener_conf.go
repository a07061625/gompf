package mpserver

import (
    "sync"
    "time"

    "github.com/a07061625/gompf/mpf"
)

type listenerTcpConf struct {
    reusePort       bool
    deferAccept     bool
    fastOpen        bool
    backlog         int
    keepAlive       bool
    keepAlivePeriod time.Duration
    noDelay         bool
    linger          int
    deadline        time.Time
    readDeadline    time.Time
    writeDeadline   time.Time
}

var (
    onceTcpConf sync.Once
    insTcpConf  *listenerTcpConf
)

func init() {
    insTcpConf = &listenerTcpConf{}
}

func newTcpConf() *listenerTcpConf {
    onceTcpConf.Do(func() {
        conf := mpf.NewConfig().GetConfig("listener")
        confPrefix := "tcp." + mpf.EnvProjectKey() + "."
        insTcpConf.reusePort = conf.GetBool(confPrefix + "ReusePort")
        insTcpConf.deferAccept = conf.GetBool(confPrefix + "DeferAccept")
        insTcpConf.fastOpen = conf.GetBool(confPrefix + "FastOpen")
        insTcpConf.backlog = conf.GetInt(confPrefix + "Backlog")
        insTcpConf.keepAlive = conf.GetBool(confPrefix + "KeepAlive")
        insTcpConf.keepAlivePeriod = time.Duration(conf.GetInt64(confPrefix+"KeepAlivePeriod")) * time.Second
        insTcpConf.noDelay = conf.GetBool(confPrefix + "NoDelay")
        insTcpConf.linger = conf.GetInt(confPrefix + "Linger")

        timeout := time.Duration(conf.GetInt64(confPrefix+"Deadline")) * time.Second
        timeoutRead := time.Duration(conf.GetInt64(confPrefix+"ReadDeadline")) * time.Second
        timeoutWrite := time.Duration(conf.GetInt64(confPrefix+"WriteDeadline")) * time.Second
        insTcpConf.deadline = time.Unix(int64(timeout), 0)
        insTcpConf.readDeadline = time.Unix(int64(timeoutRead), 0)
        insTcpConf.writeDeadline = time.Unix(int64(timeoutWrite), 0)
    })

    return insTcpConf
}
