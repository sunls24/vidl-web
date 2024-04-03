package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	"time"
	"strings"
)

func init() {
	if gin.Mode() == gin.ReleaseMode {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		return
	}
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

var logger = zerolog.New(gin.DefaultWriter).
	Output(
		zerolog.ConsoleWriter{
			Out:        gin.DefaultWriter,
			TimeFormat: "06-01-02 15:04:05",
		},
	).
	With().
	Timestamp().
	Logger()

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()
		if path != "/" && !strings.HasPrefix(path, "/api/") {
			return
		}

		end := time.Now()
		latency := end.Sub(start)

		msg := "Request"
		if len(c.Errors) > 0 {
			msg = c.Errors.String()
		}

		var evt *zerolog.Event
		switch {
		case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
			evt = Warn()
		case c.Writer.Status() >= http.StatusInternalServerError:
			evt = Error()
		default:
			evt = Info()
		}
		evt.
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("ip", c.ClientIP()).
			Dur("latency", latency).
			Str("user_agent", c.Request.UserAgent()).
			Msg(msg)
	}
}

func Info() *zerolog.Event {
	return logger.Info()
}

func Warn() *zerolog.Event {
	return logger.Warn()
}

func Error() *zerolog.Event {
	return logger.Error()
}

func Panic() *zerolog.Event {
	return logger.Panic()
}

func Fatal() *zerolog.Event {
	return logger.Fatal()
}
