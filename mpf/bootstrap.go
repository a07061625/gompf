// Package mpf bootstrap
// User: 姜伟
// Time: 2020-02-19 04:59:21
package mpf

import (
    "log"
    "os"
    "regexp"
    "strconv"
    "strings"
    "sync"

    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mplog"
)

type bootstrap struct {
    dirRoot       string // 项目根目录
    dirConfigs    string // 配置目录
    dirLogs       string // 日志目录
    envType       string // 环境类型
    projectTag    string // 项目标识
    projectModule string // 项目模块
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
    }
    return trueDir
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

func (bs *bootstrap) SetEnvType(envType string) {
    bs.envType = envType
}

func (bs *bootstrap) CheckEnvType() string {
    if (bs.envType != EnvTypeDev) && (bs.envType != EnvTypeProduct) {
        log.Fatalln("环境类型不支持")
    }

    return bs.envType
}

func (bs *bootstrap) SetProjectTag(projectTag string) {
    bs.projectTag = projectTag
}

func (bs *bootstrap) CheckProjectTag() string {
    match, _ := regexp.MatchString(`^[0-9a-z]{3}$`, bs.projectTag)
    if !match {
        log.Fatalln("项目标识不合法")
    }

    return bs.projectTag
}

func (bs *bootstrap) SetProjectModule(projectModule string) {
    bs.projectModule = projectModule
}

func (bs *bootstrap) CheckProjectModule() string {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]+$`, bs.projectModule)
    if !match {
        log.Fatalln("项目模块不合法")
    }
    return bs.projectModule
}

// NewBootstrap 实例化
func NewBootstrap() *bootstrap {
    return &bootstrap{"", "", "", "", "", ""}
}

var (
    onceBoot sync.Once
)

// LoadBoot 初始化加载
func LoadBoot(bs *bootstrap) {
    onceBoot.Do(func() {
        // 配置相关
        insConfig.dirConfigs = bs.CheckDirConfigs()

        insEnv.envType = bs.CheckEnvType()
        insEnv.projectTag = bs.CheckProjectTag()
        insEnv.projectModule = bs.CheckProjectModule()
        insEnv.projectKey = bs.CheckEnvType() + bs.CheckProjectTag()
        insEnv.projectKeyModule = bs.CheckProjectTag() + bs.CheckProjectModule()
        os.Setenv(GoEnvEnvType, insEnv.envType)
        os.Setenv(GoEnvProjectTag, insEnv.projectTag)
        os.Setenv(GoEnvProjectModule, insEnv.projectModule)
        os.Setenv(GoEnvProjectKey, insEnv.projectKey)
        os.Setenv(GoEnvProjectKeyModule, insEnv.projectKeyModule)

        // 环境相关
        serverConfig := NewConfig().GetConfig("server")
        serverConfigPrefix := EnvType() + "." + EnvProjectKeyModule() + "."
        serverHost := serverConfig.GetString(serverConfigPrefix + "host")
        serverPort := serverConfig.GetInt(serverConfigPrefix + "port")
        serverType := serverConfig.GetString(serverConfigPrefix + "type")
        if (serverPort <= 1024) || (serverPort > 65535) {
            log.Fatalln("服务端口不合法")
        }
        if (serverType != EnvServerTypeAPI) && (serverType != EnvServerTypeRPC) {
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
