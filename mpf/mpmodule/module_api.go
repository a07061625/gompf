// Package mpmodule request_api
// User: 姜伟
// Time: 2020-02-19 05:03:27
package mpmodule

// ModuleAPI api模块
type ModuleAPI struct {
    ModuleBasic
}

// NewAPI 实例化
func NewAPI(tag string) ModuleAPI {
    return ModuleAPI{NewBasic(tag)}
}
