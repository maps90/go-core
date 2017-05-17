package core

import (
	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"
	dm "github.com/mataharimall/digital/core/middleware"
	config "github.com/spf13/viper"
)

type Router interface {
	Get() *echo.Echo
	Set(echo *echo.Echo)
}

type RouterSetup interface {
	SetPort(port string)
	SetDebug(d bool)
	Run()
}

type Route struct {
	port    string
	handler *echo.Echo
	debug   bool
}

type RouteFactory func(e Router) (Router, error)

func NewRouter() Router {
	return &Route{
		handler: echo.New(),
	}
}

func (r *Route) Get() *echo.Echo {
	return r.handler
}

func (r *Route) Set(echo *echo.Echo) {
	r.handler = echo
}

func (r *Route) Setup() RouterSetup {
	return r
}

func (r *Route) Run() {
	echo := r.handler
	echo.Debug = r.debug

	echo = r.useMiddleware(echo)
	echo.Start(":" + r.port)
}

func (r *Route) SetDebug(d bool) {
	r.debug = d
}

func (r *Route) SetPort(port string) {
	r.port = port
}

func (r *Route) useMiddleware(echo *echo.Echo) *echo.Echo {

	echo.Use(em.Recover())
	echo.Use(em.Gzip())
	if r.debug {
		echo.Use(dm.Logger())
	}

	return echo
}
