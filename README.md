# 简介
go语言多功能框架(go multi-purpose frame),集成多种常用功能,包括微信,支付宝,短信,物流等众多实用功能

# 安装
    go get github.com/a07061625/gompf/mpf
    git clone https://github.com/a07061625/gompf GOPATH/pkg/src/github.com/a07061625/gompf

# 使用
## 初始化
    bs := mpf.NewBootstrap()
    bs.SetDirRoot(dirRoot)
    bs.SetDirConfigs(dirConfigs)
    bs.SetDirLogs(dirLogs)
    github.com/a07061625/gompf/mpf.LoadBoot(bs)
    // 进行上述设置后就可以愉快的使用了

## 命令行参数
### 必填
- -mppt: 项目标识,小写字母和数字组成的3位长度字符串
- -mppm: 项目模块,字母和数字组成的字符串
- -mpot: 操作类型,start:启动服务 stop:停止服务 restart:重启服务

### 选填
- -mpet: 环境类型,默认为product,dev:测试 product:生产

# 框架详解
## 代码结构
- configs: 配置文件
- mpf: 源码

https://www.cnblogs.com/ascii0x03/p/8781643.html

## 自定义常量
- MP_DIR_ROOT: 项目根目录
- MP_ENV_TYPE: 环境类型
- MP_PROJECT_TAG: 项目标识
- MP_PROJECT_MODULE: 项目模块
- MP_PROJECT_KEY: 项目代号,由环境类型+项目标识组成
- MP_PROJECT_KEY_MODULE: 项目模块代号,由项目标识+项目模块组成
