/**
 * 框架常量
 * User: 姜伟
 * Date: 2019/12/24 0024
 * Time: 15:39
 */
package frame

const (
    MWEventGlobalPrefix    = 1 // 全局事件-前置处理
    MWEventGlobalSuffix    = 2 // 全局事件-后置处理
    MWEventMvcModule       = 1 // mvc事件-模块处理
    MWEventMvcController   = 2 // mvc事件-控制器处理
    MWEventMvcActionPrefix = 3 // mvc事件-前置动作处理
    MWEventMvcActionSuffix = 4 // mvc事件-后置动作处理
)
