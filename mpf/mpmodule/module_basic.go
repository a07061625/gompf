// Package mpmodule request
// User: 姜伟
// Time: 2020-02-19 05:03:55
package mpmodule

import "github.com/a07061625/gompf/mpf"

const (
    // NodeTypeHTTP 节点类型－http
    NodeTypeHTTP = "http"
    // NodeTypeRPC 节点类型－rpc
    NodeTypeRPC = "rpc"
)

// ModuleBasic 基础模块
type ModuleBasic struct {
    tag  string
    name string
}

// GetModuleTag 获取模块标识
func (m *ModuleBasic) GetModuleTag() string {
    return m.tag
}

// GetModuleName 获取模块名称
func (m *ModuleBasic) GetModuleName() string {
    return m.name
}

// NewBasic 实例化
func NewBasic(tag string) ModuleBasic {
    m := ModuleBasic{}
    m.tag = tag
    m.name = mpf.EnvProjectTag() + tag
    return m
}
