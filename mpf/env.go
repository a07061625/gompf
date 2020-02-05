/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/23 0023
 * Time: 15:35
 */
package mpf

import (
    "flag"
    "log"
    "os"
    "regexp"
    "sync"

    "github.com/spf13/viper"
)

const (
    EnvTypeDev     = "dev"     // 环境类型-测试
    EnvTypeProduct = "product" // 环境类型-生产

    EnvServiceTypeApi = "api" // 服务类型-API
    EnvServiceTypeRpc = "rpc" // 服务类型-RPC
)

type env struct {
    envType string // 环境类型

    projectTag       string // 项目标识
    projectModule    string // 项目模块
    projectKey       string // 项目代号
    projectKeyModule string // 项目模块代号

    serviceHost string
    servicePort uint
    serviceType string
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

func EnvServiceHost() string {
    return insEnv.serviceHost
}

func EnvServicePort() uint {
    return insEnv.servicePort
}

func EnvServiceType() string {
    return insEnv.serviceType
}

var (
    onceEnv sync.Once
    insEnv  *env

    envType       = flag.String("mpet", EnvTypeProduct, "环境类型,只能是dev或product")
    projectTag    = flag.String("mppt", "", "项目标识,由小写字母和数字组成的3位长度字符串")
    projectModule = flag.String("mppm", "", "项目模块,由字母和数字组成的字符串")
)

func init() {
    flag.Parse()

    if (*envType != EnvTypeDev) && (*envType != EnvTypeProduct) {
        log.Fatalln("环境类型不支持")
    }
    match, _ := regexp.MatchString(`^[0-9a-z]{3}$`, *projectTag)
    if !match {
        log.Fatalln("项目标识不合法")
    }
    match, _ = regexp.MatchString(`^[0-9a-zA-Z]+$`, *projectModule)
    if !match {
        log.Fatalln("项目模块不合法")
    }

    insEnv = &env{}
    insEnv.envType = *envType
    insEnv.projectTag = *projectTag
    insEnv.projectModule = *projectModule
    insEnv.projectKey = *envType + *projectTag
    insEnv.projectKeyModule = *projectTag + *projectModule
    os.Setenv("MP_ENV_TYPE", insEnv.envType)
    os.Setenv("MP_PROJECT_TAG", insEnv.projectTag)
    os.Setenv("MP_PROJECT_MODULE", insEnv.projectModule)
    os.Setenv("MP_PROJECT_KEY", insEnv.projectKey)
    os.Setenv("MP_PROJECT_KEY_MODULE", insEnv.projectKeyModule)
}

func LoadEnv(conf *viper.Viper) {
    onceEnv.Do(func() {
        serviceHost := conf.GetString(insEnv.envType + "." + insEnv.projectKeyModule + ".host")
        servicePort := conf.GetUint(insEnv.envType + "." + insEnv.projectKeyModule + ".port")
        serviceType := conf.GetString(insEnv.envType + "." + insEnv.projectKeyModule + ".type")
        if (servicePort < 1024) || (servicePort > 65535) {
            log.Fatalln("服务端口不合法")
        }
        if (serviceType != EnvServiceTypeApi) && (serviceType != EnvServiceTypeRpc) {
            log.Fatalln("服务类型不支持")
        }

        insEnv.serviceHost = serviceHost
        insEnv.servicePort = servicePort
        insEnv.serviceType = serviceType
    })
}
