// Package mpvalidator validator
// User: 姜伟
// Time: 2020-02-19 06:31:26
package mpvalidator

const (
    // FieldIgnoreSign 字段-忽略接口签名
    FieldIgnoreSign = "_sign"
)

// Filter Filter
type Filter struct {
    Field     string                 `json:"field"`      // 字段名
    Desc      string                 `json:"desc"`       // 字段描述
    DataType  string                 `json:"data_type"`  // 字段数据类型
    DataRules map[string]interface{} `json:"data_rules"` // 字段数据规则
}

// Filters Filters
type Filters []Filter
