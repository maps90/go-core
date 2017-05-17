package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

func Logger(name string) echo.MiddlewareFunc {
	l := log.StandardLogger()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			remoteAddr := req.RemoteAddr
			if ip := req.Header.Get(echo.HeaderXRealIP); ip != "" {
				remoteAddr = ip
			} else if ip = req.Header.Get(echo.HeaderXForwardedFor); ip != "" {
				remoteAddr = ip
			} else {
				remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
			}

			entry := l.WithFields(log.Fields{
				"request": req.RequestURI,
				"method":  req.Method,
				"remote":  remoteAddr,
			})

			if reqID := req.Header.Get("X-Request-Id"); reqID != "" {
				entry = entry.WithField("request_id", reqID)
			}

			entry.Info("started handling request")

			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			latency := stop.Sub(start)

			entry.WithFields(log.Fields{
				"size":        res.Size,
				"status":      res.Status,
				"text_status": http.StatusText(res.Status),
				"took":        strconv.FormatInt(int64(latency), 10),
				fmt.Sprintf("#%s.latency", name): stop.Sub(start).String(),
			}).Info("completed handling request")

			return nil
		}
	}
}
