package main

import (
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
)

func init() {
    dir, err := os.Getwd()
    if err != nil {
        log.Fatalln("get root dir error")
    }

    dirRoot := strings.Replace(dir, "\\", "/", -1)
    configDir := dirRoot + "/configs"
    mpf.LoadConfig(configDir)
    serviceConfig := mpf.NewConfig().GetConfig("service")
    mpf.LoadEnv(serviceConfig)
    projectConfig := mpf.NewConfig().GetConfig("project")
    project.LoadProject(projectConfig)
    logDir := dirRoot + "/logs"
    logConfig := mpf.NewConfig().GetConfig("log")
    mpf.LoadLog(logDir, logConfig)
}

func main() {
    defer mpf.ToolHandleError()()
    // headers := make(map[string]interface{})
    // configs := make(map[string]interface{})
    // data, err := tool.SendHttpReq("https://www.baidu.com", headers, configs)
    // if err != nil {
    //     log2.NewZap().Error(err.Error())
    // } else {
    //     fmt.Println(string(data))
    // }

    // info1 := make(map[string]string)
    // info1["test1"] = "1111"
    // info2 := make(map[string]string)
    // info2["test2"] = "2222"
    // infos := make(map[string]interface{})
    // infos["1111"] = info1
    // infos["2222"] = info2
    // infos2 := make(map[string]interface{})
    // infos2["3333"] = infos
    // // infos3 := make([]map[string]string, 0)
    // // infos3 = append(infos3, info1)
    // // infos3 = append(infos3, info2)
    // // infos4 := make([]map[string]interface{}, 0)
    // // infos4 = append(infos4, infos)
    // data := mxj.Map{}
    // data["info1"] = "0000"
    // data["info2"] = "9999"
    // data["info0"] = info1
    // data["infos"] = infos
    // data["infos2"] = infos2
    // // data["infos3"] = infos3
    // // data["infos4"] = infos4
    // xmlStr, _ := data.Xml("xml")
    // fmt.Println(string(xmlStr))
    // xmlStr := "<AccessControlPolicy><Owner><ID>string</ID><DisplayName>string</DisplayName></Owner><AccessControlList><Grant><Grantee xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:type=\"Group\"><URI>string</URI></Grantee><Permission>Enum</Permission></Grant><Grant><Grantee xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:type=\"CanonicalUser\"><ID>string</ID><DisplayName>string</DisplayName></Grantee><Permission>Enum</Permission></Grant></AccessControlList></AccessControlPolicy>"
    // mv, _ := mxj.NewMapXml([]byte(xmlStr))
    // fmt.Println(mpf.JsonMarshal(mv))
    os.Setenv("aaaa", "cccc")
    str := project.RedisPrefix(project.RedisPrefixPrintFeiYinAccount)
    fmt.Println(os.Getenv("aaaa"))
    fmt.Println(str)
}
