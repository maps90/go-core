package middleware

import (
	"fmt"
	"runtime"

	"github.com/labstack/echo"
	"github.com/getsentry/raven-go"
)

func AppRecover(env string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					var err error
					switch r := r.(type) {
					case error:
						err = r
					default:
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, 4 << 10)
					length := runtime.Stack(stack, true)
					c.Logger().Printf("[%s] %s %s\n", "PANIC RECOVER", err, stack[:length])
					raven.CaptureError(err, map[string]string{
						"env": env,
						"error": err.Error(),
					})
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
