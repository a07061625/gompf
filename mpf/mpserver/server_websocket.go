package mpserver

import (
    "net/http"
    "strconv"
    "sync"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/a07061625/gompf/mpf/mpresponse"
    "github.com/gorilla/websocket"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/context"
)

type HandlerFunc func(messageType int, messageData map[string]interface{}) (bool, *mpresponse.ResultAPI)

type connWebSocket struct {
    Id   string
    Conn *websocket.Conn
}

type IServerWebSocket interface {
    Handler() context.Handler
}

type serverWebSocket struct {
    upGrader    websocket.Upgrader
    connections map[string]*connWebSocket
    connAdd     sync.Mutex
    connClear   sync.Mutex
    handlers    map[string]HandlerFunc
}

func (s *serverWebSocket) ClearConnection() int {
    s.connClear.Lock()
    defer s.connClear.Unlock()

    clearNum := len(s.connections)
    for k, v := range s.connections {
        err := v.Conn.Close()
        if err != nil {
            mplog.LogError("web socket connection " + k + " close error: " + err.Error())
        }
    }
    s.connections = make(map[string]*connWebSocket)

    return clearNum
}

func (s *serverWebSocket) RemoveConnection(key string) {
    conn, ok := s.connections[key]
    if !ok {
        return
    }

    err := conn.Conn.Close()
    if err != nil {
        mplog.LogError("web socket connection " + key + " close error: " + err.Error())
    }
    delete(s.connections, key)
}

func (s *serverWebSocket) AddConnection(key string, conn *websocket.Conn, expireTime int64) {
    s.connAdd.Lock()
    defer s.connAdd.Unlock()
    _, ok := s.connections[key]
    if ok {
        s.RemoveConnection(key)
    }

    nowTime := time.Now().Unix()
    connection := &connWebSocket{}
    connection.Id = mpf.ToolCreateNonceStr(8, "numlower") + strconv.FormatInt(nowTime, 10)
    connection.Conn = conn
    s.connections[key] = connection
}

func (s *serverWebSocket) SetUpGrader(upGrader websocket.Upgrader) {
    s.upGrader = upGrader
}

func (s *serverWebSocket) SetHandlers(handlers map[string]HandlerFunc) {
    s.handlers = handlers
}

func (s *serverWebSocket) handleUpGrader(ctx context.Context) (string, *websocket.Conn) {
    conn, err := s.upGrader.Upgrade(ctx.ResponseWriter(), ctx.Request(), ctx.ResponseWriter().Header())
    if err != nil {
        mplog.LogError("web socket create connection error: " + err.Error())
        ctx.StatusCode(iris.StatusServiceUnavailable)
        return "", nil
    }
    conn.EnableWriteCompression(true)
    conn.SetCompressionLevel(3)

    webSocketKey := ctx.Request().Header.Get("Sec-WebSocket-Key")
    s.AddConnection(webSocketKey, conn, 60)
    return webSocketKey, conn
}

func (s *serverWebSocket) handleMessage(messageType int, messageData []byte, connKey string) (bool, *mpresponse.ResultAPI) {
    result := mpresponse.NewResultAPI()

    handlerData, err := mpf.JsonUnmarshalMap(string(messageData))
    if err != nil {
        result.Code = errorcode.CommonBaseWebSocket
        result.Msg = "解析数据出错"
        return false, result
    }
    msgEvent, ok := handlerData["event"]
    if !ok {
        result.Code = errorcode.CommonBaseWebSocket
        result.Msg = "事件类型不存在"
        return false, result
    }

    switch msgEvent.(type) {
    case string:
        handlerFunc, ok := s.handlers[msgEvent.(string)]
        if ok {
            return handlerFunc(messageType, handlerData)
        } else {
            result.Code = errorcode.CommonBaseWebSocket
            result.Msg = "事件类型不支持"
        }
    default:
        result.Code = errorcode.CommonBaseWebSocket
        result.Msg = "事件类型数据格式错误"
    }
    return false, result
}

func (s *serverWebSocket) Handler() context.Handler {
    return func(ctx context.Context) {
        key, conn := s.handleUpGrader(ctx)
        if conn == nil {
            return
        }
        defer s.RemoveConnection(key)

        go func() {
            closeFlag := false
            result := mpresponse.NewResultAPI()
            for {
                result.Refresh()
                messageType, messageData, err := conn.ReadMessage()
                if err != nil {
                    mplog.LogError("web socket connection read message error: " + err.Error())
                    result.Code = errorcode.CommonBaseWebSocket
                    result.Msg = "读取消息出错"
                } else {
                    closeFlag, result = s.handleMessage(messageType, messageData, key)
                }
                conn.WriteMessage(messageType, []byte(mpf.JsonMarshal(result)))
                if closeFlag {
                    s.RemoveConnection(key)
                    break
                }
            }
        }()
    }
}

func NewServerWebSocket() *serverWebSocket {
    s := &serverWebSocket{}
    s.upGrader = websocket.Upgrader{
        HandshakeTimeout:  3 * time.Second,
        ReadBufferSize:    4096,
        WriteBufferSize:   4096,
        EnableCompression: true,
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
        Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
            mplog.LogError("web socket error: " + reason.Error() + ", status: " + strconv.Itoa(status))

            result := mpresponse.NewResultAPI()
            result.Code = errorcode.CommonBaseWebSocket
            result.Msg = reason.Error()
            w.Write([]byte(mpf.JsonMarshal(result)))
        },
    }
    s.connections = make(map[string]*connWebSocket)
    s.connAdd = sync.Mutex{}
    s.connClear = sync.Mutex{}
    s.handlers = make(map[string]HandlerFunc)

    return s
}
