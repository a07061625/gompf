// Package mpf json
// User: 姜伟
// Time: 2020-02-19 05:38:11
package mpf

import (
    "strings"

    "github.com/a07061625/gompf/mpf/mplog"
    jsoniter "github.com/json-iterator/go"
)

type jsonTransform struct {
    Config jsoniter.API
}

var (
    insJSON *jsonTransform
)

func init() {
    insJSON = &jsonTransform{}
    insJSON.Config = jsoniter.ConfigCompatibleWithStandardLibrary
}

// JSONMarshal 序列化
func JSONMarshal(data interface{}) string {
    res, err := insJSON.Config.Marshal(data)
    if err != nil {
        mplog.LogError(err.Error())
    }
    return string(res)
}

// JSONUnmarshal 反序列化
func JSONUnmarshal(data []byte, obj interface{}) error {
    return insJSON.Config.Unmarshal(data, obj)
}

// JSONUnmarshalMap 解析未知结构json字符串为map
func JSONUnmarshalMap(data string) (map[string]interface{}, error) {
    reader := strings.NewReader(data)
    decoder := insJSON.Config.NewDecoder(reader)
    result := make(map[string]interface{})
    err := decoder.Decode(&result)
    if err != nil {
        return nil, err
    }
    return result, nil
}
