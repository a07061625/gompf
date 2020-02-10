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
    "github.com/a07061625/gompf/mpf/mplog"
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
        serverConfig := NewConfig().GetConfig("server")
        serverConfigPrefix := EnvType() + "." + EnvProjectKeyModule() + "."
        serverHost := serverConfig.GetString(serverConfigPrefix + "host")
        serverPort := serverConfig.GetInt(serverConfigPrefix + "port")
        serverType := serverConfig.GetString(serverConfigPrefix + "type")
        if (serverPort <= 1024) || (serverPort > 65535) {
            log.Fatalln("服务端口不合法")
        }
        if (serverType != EnvServerTypeApi) && (serverType != EnvServerTypeRpc) {
            log.Fatalln("服务类型不支持")
        }
        insEnv.serverHost = serverHost
        insEnv.serverPort = serverPort
        insEnv.serverType = serverType
        insEnv.serverDomain = serverHost + ":" + strconv.Itoa(serverPort)
        insEnv.dirRoot = bs.CheckDirRoot()
        insEnv.dirConfigs = bs.CheckDirConfigs()
        os.Setenv(GoEnvDirRoot, insEnv.dirRoot)

        // 项目相关
        projectConfig := NewConfig().GetConfig("project")
        project.LoadProject(projectConfig)

        // 日志相关
        logConfig := NewConfig().GetConfig("log")
        loggerFields := make(map[string]interface{})
        logExtend := make(map[string]interface{})
        logExtend["log_dir"] = bs.CheckDirLogs() + "/" + EnvProjectKey()
        logExtend["conf_prefix"] = "zap." + EnvProjectKey() + "."
        logExtend["env_type"] = EnvType()
        logExtend["project_tag"] = EnvProjectTag()
        logExtend["project_module"] = EnvProjectModule()
        logExtend["server_host"] = serverHost
        logExtend["server_port"] = strconv.Itoa(serverPort)
        mplog.Load(logConfig, loggerFields, logExtend)
    })
}
