package mpapp

func (app *appBasic) Tel() {
    //serverEvents := websocket.Namespaces{
    //    "default": websocket.Events{
    //        websocket.OnNamespaceConnected: func(nsConn *websocket.NSConn, msg websocket.Message) error {
    //            // with `websocket.GetContext` you can retrieve the Iris' `Context`.
    //            ctx := websocket.GetContext(nsConn.Conn)
    //
    //            log.Printf("[%s] connected to namespace [%s] with IP [%s]",
    //                nsConn, msg.Namespace,
    //                ctx.RemoteAddr())
    //            return nil
    //        },
    //        websocket.OnNamespaceDisconnect: func(nsConn *websocket.NSConn, msg websocket.Message) error {
    //            log.Printf("[%s] disconnected from namespace [%s]", nsConn, msg.Namespace)
    //            return nil
    //        },
    //        "chat": func(nsConn *websocket.NSConn, msg websocket.Message) error {
    //            // room.String() returns -> NSConn.String() returns -> Conn.String() returns -> Conn.ID()
    //            log.Printf("[%s] sent: %s", nsConn, string(msg.Body))
    //
    //            // Write message back to the client message owner with:
    //            // nsConn.Emit("chat", msg)
    //            // Write message to all except this client with:
    //            nsConn.Conn.Server().Broadcast(nsConn, msg)
    //            return nil
    //        },
    //    },
    //}
    //ws := websocket.New(websocket.DefaultGorillaUpgrader, serverEvents)
}
