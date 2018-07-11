package main

import (
	"github.com/hydah/golib/httpsvr"
	"github.com/hydah/golib/logger"
)

type Controller struct {
	httpsvr.Controller
}

func (c *Controller) Get(ctx *httpsvr.Context) {
	ctx.Text("OK", 200)
}

func main() {
	s := httpsvr.NewHTTPServer()
	s.Init()
	s.UseLogger()
	s.Get("/", &Controller{})
	s.ServeFile("/static/t.txt", "/tmp/t.txt")

	logger.Debug("http server is ready")

	s.Run(":8888")
}
