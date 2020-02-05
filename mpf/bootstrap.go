/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/6 0006
 * Time: 0:59
 */
package mpf

import (
    "log"
    "os"
    "strconv"
    "strings"
    "sync"

    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/robfig/cron"
)

type bootstrap struct {
    dirRoot    string // 项目根目录
    dirConfigs string // 配置目录
    dirLogs    string // 日志目录
}

func (bs *bootstrap) formatDir(dir string) string {
    trueDir := strings.Replace(dir, "\\", "/", -1)
    f, err := os.Stat(trueDir)
    if err != nil {
        log.Fatalln(dir + " invalid")
    }
    if !f.IsDir() {
        log.Fatalln(dir + " is not dir")
    }

    if (len(trueDir) > 1) && strings.HasSuffix(trueDir, "/") {
        return strings.TrimSuffix(trueDir, "/")
    } else {
        return trueDir
    }
}

func (bs *bootstrap) SetDirRoot(dirRoot string) {
    bs.dirRoot = bs.formatDir(dirRoot)
}

func (bs *bootstrap) CheckDirRoot() string {
    if len(bs.dirRoot) == 0 {
        log.Fatalln("root dir must be set")
    }
    return bs.dirRoot
}

func (bs *bootstrap) SetDirConfigs(dirConfigs string) {
    bs.dirConfigs = bs.formatDir(dirConfigs)
}

func (bs *bootstrap) CheckDirConfigs() string {
    if len(bs.dirConfigs) == 0 {
        log.Fatalln("configs dir must be set")
    }
    return bs.dirConfigs
}

func (bs *bootstrap) SetDirLogs(dirLogs string) {
    bs.dirLogs = bs.formatDir(dirLogs)
}

func (bs *bootstrap) CheckDirLogs() string {
    if len(bs.dirLogs) == 0 {
        log.Fatalln("logs dir must be set")
    }
    return bs.dirLogs
}

func NewBootstrap() *bootstrap {
    return &bootstrap{"", "", ""}
}

var (
    onceBoot sync.Once
)

func LoadBoot(bs *bootstrap) {
    onceBoot.Do(func() {
        // 配置相关
        insConfig.dirConfigs = bs.CheckDirConfigs()

        // 环境相关
        serviceConfig := NewConfig().GetConfig("service")
        serviceConfigPrefix := EnvType() + "." + EnvProjectKeyModule() + "."
        serviceHost := serviceConfig.GetString(serviceConfigPrefix + "host")
        servicePort := serviceConfig.GetUint(serviceConfigPrefix + "port")
        serviceType := serviceConfig.GetString(serviceConfigPrefix + "type")
        if (servicePort <= 1024) || (servicePort > 65535) {
            log.Fatalln("服务端口不合法")
        }
        if (serviceType != EnvServiceTypeApi) && (serviceType != EnvServiceTypeRpc) {
            log.Fatalln("服务类型不支持")
        }
        insEnv.serviceHost = serviceHost
        insEnv.servicePort = servicePort
        insEnv.serviceType = serviceType
        insEnv.dirRoot = bs.CheckDirRoot()
        os.Setenv("MP_DIR_ROOT", insEnv.dirRoot)

        // 项目相关
        projectConfig := NewConfig().GetConfig("project")
        project.LoadProject(projectConfig)

        // 日志相关
        insLog.logDir = bs.CheckDirLogs() + "/" + EnvProjectKey()
        err := os.MkdirAll(insLog.logDir, os.ModePerm)
        if err != nil {
            log.Fatalln("log dir create fail:" + err.Error())
        }
        logConfig := NewConfig().GetConfig("log")
        logConfigPrefix := "zap." + EnvProjectKey() + "."
        insLog.SetLogAccess(logConfig.GetString(logConfigPrefix + "access"))
        insLog.SetLogError(logConfig.GetString(logConfigPrefix + "error"))
        insLog.SetLogSuffix(logConfig.GetString(logConfigPrefix + "suffix"))

        c := cron.New()
        c.AddFunc(logConfig.GetString(logConfigPrefix+"cron.access"), insLog.ChangeAccessLog)
        c.AddFunc(logConfig.GetString(logConfigPrefix+"cron.error"), insLog.ChangeErrorLog)
        c.Start()
        insLog.SetCron(c)
        insLog.ChangeAccessLog()
        insLog.ChangeErrorLog()

        fields := logConfig.GetStringMapString(logConfigPrefix + "fields")
        fields["env_type"] = EnvType()
        fields["env_project_tag"] = EnvProjectTag()
        fields["env_project_module"] = EnvProjectModule()
        fields["env_service_host"] = serviceHost
        fields["env_service_port"] = strconv.Itoa(int(servicePort))
        insLog.createLogger(fields)
    })
}
