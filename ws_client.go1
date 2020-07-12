package main

import (
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/context"
)

func main() {
    appClient := iris.New()
    appClient.RegisterView(iris.HTML("./views", ".html"))
    appClient.Get("/echo", func(ctx context.Context) {
        ctx.View("ws_client.html", "112233")
    })
    appClient.Run(iris.Addr(":8989"), iris.WithoutServerError(iris.ErrServerClosed))
}
