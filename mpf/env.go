/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/23 0023
 * Time: 15:35
 */
package mpf

const (
    EnvTypeDev     = "dev"     // 环境类型-测试
    EnvTypeProduct = "product" // 环境类型-生产

    EnvServerTypeApi    = "api"    // 服务类型-API
    EnvServerTypeRpc    = "rpc"    // 服务类型-RPC
    EnvServerModeDaemon = "daemon" // 服务运行模式-守护进程
    EnvServerModeChild  = "child"  // 服务运行模式-子进程

    GoEnvReqId            = "MP_REQ_ID"
    GoEnvDirRoot          = "MP_DIR_ROOT"
    GoEnvEnvType          = "MP_ENV_TYPE"
    GoEnvProjectTag       = "MP_PROJECT_TAG"
    GoEnvProjectModule    = "MP_PROJECT_MODULE"
    GoEnvProjectKey       = "MP_PROJECT_KEY"
    GoEnvProjectKeyModule = "MP_PROJECT_KEY_MODULE"
    GoEnvServerMode       = "MP_SERVER_MODE" // 服务运行模式 daemon:守护进程 child:子进程
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

func EnvDirRoot() string {
    return insEnv.dirRoot
}

func EnvDirConfigs() string {
    return insEnv.dirConfigs
}

func EnvType() string {
    return insEnv.envType
}

func EnvProjectTag() string {
    return insEnv.projectTag
}

func EnvProjectModule() string {
    return insEnv.projectModule
}

func EnvProjectKey() string {
    return insEnv.projectKey
}

func EnvProjectKeyModule() string {
    return insEnv.projectKeyModule
}

func EnvServerHost() string {
    return insEnv.serverHost
}

func EnvServerPort() int {
    return insEnv.serverPort
}

func EnvServerDomain() string {
    return insEnv.serverDomain
}

func EnvServerType() string {
    return insEnv.serverType
}

var (
    insEnv *env
)

func init() {
    insEnv = &env{}
}
