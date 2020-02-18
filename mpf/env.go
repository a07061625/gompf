// Package mpf env
// User: 姜伟
// Time: 2020-02-19 05:14:47
package mpf

const (
    // EnvTypeDev 环境类型-测试
    EnvTypeDev = "dev"
    // EnvTypeProduct 环境类型-生产
    EnvTypeProduct = "product"

    // EnvServerTypeAPI 服务类型-API
    EnvServerTypeAPI = "api"
    // EnvServerTypeRPC 服务类型-RPC
    EnvServerTypeRPC = "rpc"
    // EnvServerModeDaemon 服务运行模式-守护进程
    EnvServerModeDaemon = "daemon"
    // EnvServerModeChild 服务运行模式-子进程
    EnvServerModeChild = "child"

    // GoEnvReqID 请求ID
    GoEnvReqID = "MP_REQ_ID"
    // GoEnvDirRoot 项目根目录
    GoEnvDirRoot = "MP_DIR_ROOT"
    // GoEnvEnvType 环境类型
    GoEnvEnvType = "MP_ENV_TYPE"
    // GoEnvProjectTag 项目标识
    GoEnvProjectTag = "MP_PROJECT_TAG"
    // GoEnvProjectModule 项目模块
    GoEnvProjectModule = "MP_PROJECT_MODULE"
    // GoEnvProjectKey 项目代号,由环境类型+项目标识组成
    GoEnvProjectKey = "MP_PROJECT_KEY"
    // GoEnvProjectKeyModule 项目模块代号,由项目标识+项目模块组成
    GoEnvProjectKeyModule = "MP_PROJECT_KEY_MODULE"
    // GoEnvServerMode 服务运行模式 daemon:守护进程 child:子进程
    GoEnvServerMode = "MP_SERVER_MODE"
)

type env struct {
    dirRoot    string // 项目根目录
    dirConfigs string // 配置目录

    envType string // 环境类型

    projectTag       string // 项目标识
    projectModule    string // 项目模块
    projectKey       string // 项目代号
    projectKeyModule string // 项目模块代号

    serverHost   string
    serverPort   int
    serverDomain string
    serverType   string
}

// EnvDirRoot DirRoot
func EnvDirRoot() string {
    return insEnv.dirRoot
}

// EnvDirConfigs DirConfigs
func EnvDirConfigs() string {
    return insEnv.dirConfigs
}

// EnvType Type
func EnvType() string {
    return insEnv.envType
}

// EnvProjectTag ProjectTag
func EnvProjectTag() string {
    return insEnv.projectTag
}

// EnvProjectModule ProjectModule
func EnvProjectModule() string {
    return insEnv.projectModule
}

// EnvProjectKey ProjectKey
func EnvProjectKey() string {
    return insEnv.projectKey
}

// EnvProjectKeyModule ProjectKeyModule
func EnvProjectKeyModule() string {
    return insEnv.projectKeyModule
}

// EnvServerHost ServerHost
func EnvServerHost() string {
    return insEnv.serverHost
}

// EnvServerPort ServerPort
func EnvServerPort() int {
    return insEnv.serverPort
}

// EnvServerDomain ServerDomain
func EnvServerDomain() string {
    return insEnv.serverDomain
}

// EnvServerType ServerType
func EnvServerType() string {
    return insEnv.serverType
}

var (
    insEnv *env
)

func init() {
    insEnv = &env{}
}
