package core

import (
	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"
	"github.com/maps90/go-core/log"
	dm "github.com/maps90/go-core/middleware"
)

type Router interface {
	Get() *echo.Echo
	Set(echo *echo.Echo)
	Setup() RouterSetup
}

type RouterSetup interface {
	Run()
	SetPort(port string) RouterSetup
	SetDebug(d bool) RouterSetup
	SetLoggerName(logName string) RouterSetup
}

type Route struct {
	port       string
	handler    *echo.Echo
	debug      bool
	loggerName string
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
	if err := echo.Start(":" + r.port); err != nil {
		log.New(log.InfoLevelLog, err.Error())
	}
}

func (r *Route) SetDebug(d bool) RouterSetup {
	r.debug = d
	return r
}

func (r *Route) SetPort(port string) RouterSetup {
	r.port = port
	return r
}

func (r *Route) SetLoggerName(logName string) RouterSetup {
	r.loggerName = logName
	return r
}

func (r *Route) useMiddleware(echo *echo.Echo) *echo.Echo {
	echo.Use(em.Recover())
	echo.Use(em.Gzip())
	if r.debug {
		echo.Use(dm.Logger(r.loggerName))
	}

	return echo
}
