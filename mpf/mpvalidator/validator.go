/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/10 0010
 * Time: 19:43
 */
package mpvalidator

const (
    FieldIgnoreSign = "_sign" // 字段-忽略接口签名
)

type Filter struct {
    Field     string                 `json:"field"`      // 字段名
    Desc      string                 `json:"desc"`       // 字段描述
    DataType  string                 `json:"data_type"`  // 字段数据类型
    DataRules map[string]interface{} `json:"data_rules"` // 字段数据规则
}

type Filters []Filter
