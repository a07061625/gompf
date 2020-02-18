// Package mpprotocol protocol_basic
// User: 姜伟
// Time: 2020-02-19 04:56:11
package mpprotocol

import (
    "encoding/binary"
    "io"
    "strings"

    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/vmihailenco/msgpack/v4"
)

const (
    // HeaderLength 请求头长度
    HeaderLength uint32 = 4
)

// ProtocolData 协议数据类型
type ProtocolData struct {
    Command string                 // 命令
    Extend  string                 // 扩展字段
    URI     string                 // URI
    Data    map[string]interface{} // 数据
}

// NewProtocolData 实例化协议数据类型
func NewProtocolData() *ProtocolData {
    pd := &ProtocolData{}
    pd.Extend = "0000"
    pd.URI = "/"
    pd.Data = make(map[string]interface{})
    return pd
}

// Pack 数据打包
func Pack(pd *ProtocolData) []byte {
    if len(pd.Command) != 4 {
        panic(mperr.NewProtocol(errorcode.ProtocolPacket, "命令不合法", nil))
    }
    if len(pd.Extend) != 4 {
        panic(mperr.NewProtocol(errorcode.ProtocolPacket, "扩展字段不合法", nil))
    }
    if !strings.HasPrefix(pd.URI, "/") {
        panic(mperr.NewProtocol(errorcode.ProtocolPacket, "URI不合法", nil))
    }

    contentByte, err := msgpack.Marshal(pd)
    if err != nil {
        mplog.LogError("data pack error: " + err.Error())
        panic(mperr.NewProtocol(errorcode.ProtocolPacket, "URI不合法", nil))
    }
    contentLength := uint32(len(contentByte))
    buffer := make([]byte, HeaderLength+contentLength)
    binary.BigEndian.PutUint32(buffer[0:4], contentLength)
    copy(buffer[4:], contentByte)

    return buffer
}

// Unpack 数据解包
func Unpack(r io.Reader) *ProtocolData {
    // 直接先取数据,防止因部分数据格式有问题导致所有数据都不可用
    headerBuffer := make([]byte, HeaderLength)
    contentLength := binary.BigEndian.Uint32(headerBuffer[0:4])
    contentBuffer := make([]byte, contentLength)
    _, err := io.ReadFull(r, headerBuffer)
    if err != nil {
        mplog.LogError("data header unpack error: " + err.Error())
        panic(mperr.NewProtocol(errorcode.ProtocolUnPacket, "数据头解包出错", nil))
    }
    _, err = io.ReadFull(r, contentBuffer)
    if err != nil {
        mplog.LogError("data body unpack error: " + err.Error())
        panic(mperr.NewProtocol(errorcode.ProtocolUnPacket, "数据体解包出错", nil))
    }

    pd := &ProtocolData{}
    err = msgpack.Unmarshal(contentBuffer, pd)
    if err != nil {
        mplog.LogError("data body unmarshal error: " + err.Error())
        panic(mperr.NewProtocol(errorcode.ProtocolUnPacket, "数据体反压缩出错", nil))
    }

    return pd
}
