package mpserver

import (
    "sync"
    "time"

    "github.com/a07061625/gompf/mpf"
)

type listenerTCPConf struct {
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
    onceTCPConf sync.Once
    insTCPConf  *listenerTCPConf
)

func init() {
    insTCPConf = &listenerTCPConf{}
}

func newTCPConf() *listenerTCPConf {
    onceTCPConf.Do(func() {
        conf := mpf.NewConfig().GetConfig("listener")
        confPrefix := "tcp." + mpf.EnvProjectKey() + "."
        insTCPConf.reusePort = conf.GetBool(confPrefix + "ReusePort")
        insTCPConf.deferAccept = conf.GetBool(confPrefix + "DeferAccept")
        insTCPConf.fastOpen = conf.GetBool(confPrefix + "FastOpen")
        insTCPConf.backlog = conf.GetInt(confPrefix + "Backlog")
        insTCPConf.keepAlive = conf.GetBool(confPrefix + "KeepAlive")
        insTCPConf.keepAlivePeriod = time.Duration(conf.GetInt64(confPrefix+"KeepAlivePeriod")) * time.Second
        insTCPConf.noDelay = conf.GetBool(confPrefix + "NoDelay")
        insTCPConf.linger = conf.GetInt(confPrefix + "Linger")

        timeout := time.Duration(conf.GetInt64(confPrefix+"Deadline")) * time.Second
        timeoutRead := time.Duration(conf.GetInt64(confPrefix+"ReadDeadline")) * time.Second
        timeoutWrite := time.Duration(conf.GetInt64(confPrefix+"WriteDeadline")) * time.Second
        insTCPConf.deadline = time.Unix(int64(timeout), 0)
        insTCPConf.readDeadline = time.Unix(int64(timeoutRead), 0)
        insTCPConf.writeDeadline = time.Unix(int64(timeoutWrite), 0)
    })

    return insTCPConf
}
