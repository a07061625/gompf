/**
 * json处理
 * User: 姜伟
 * Date: 2019/12/24 0024
 * Time: 9:59
 */
package mpf

import (
    "strings"

    jsoniter "github.com/json-iterator/go"
)

type jsonTransform struct {
    Config jsoniter.API
}

var (
    insJson *jsonTransform
)

func init() {
    insJson = &jsonTransform{}
    insJson.Config = jsoniter.ConfigCompatibleWithStandardLibrary
}

func JsonMarshal(data interface{}) string {
    res, err := insJson.Config.Marshal(data)
    if err != nil {
        NewLogger().Error(err.Error())
    }
    return string(res)
}

func JsonUnmarshal(data []byte, obj interface{}) error {
    return insJson.Config.Unmarshal(data, obj)
}

// 解析未知结构json字符串为map
func JsonUnmarshalMap(data string) (map[string]interface{}, error) {
    reader := strings.NewReader(data)
    decoder := insJson.Config.NewDecoder(reader)
    result := make(map[string]interface{})
    err := decoder.Decode(&result)
    if err != nil {
        return nil, err
    }
    return result, nil
}
