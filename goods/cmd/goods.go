package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	_ "rubik/server/goods/init"
	"rubik/server/goods/routers"
)

func main()  {
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())
	routers.SetRouter(app, "/api/v1")
	app.Run(iris.Addr(":4006"), iris.WithoutPathCorrection)
}