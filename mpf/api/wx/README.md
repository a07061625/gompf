# 介绍
- account:公众号相关
- corp:企业号相关
- mini:小程序相关
- open:第三方开放平台相关
- provider:企业服务商相关
- single: 公众号和小程序共同内容

# 使用
    //在框架的bootstrap方法或者初始化方法中执行如下操作

    定义一个微信配置结构体并实现github.com/a07061625/gompf/mpf/api/wx/config.go文件中的IWxConfig接口
    实例化微信配置结构体
    调用github.com/a07061625/gompf/mpf/api/wx.LoadConfig()并将上述配置实例作为参数传入
    微信配置文件请参考github.com/a07061625/gompf/configs/wx.yaml

    定义一个微信工具结构体并实现github.com/a07061625/gompf/mpf/api/wx/util.go文件中的IWxOuter接口
    实例化微信工具结构体
    调用github.com/a07061625/gompf/mpf/api/wx.LoadUtil()并将上述工具实例作为参数传入
    完成上述操作后就可以正常使用微信接口了
