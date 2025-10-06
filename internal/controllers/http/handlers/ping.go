package healthcheck

import (
	"io"

	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog/log"
)
//add PingHandler godoc
func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := io.WriteString(c.Writer, "working as well"); err != nil {
			log.Error().Err(err).Msg("service is not working")
		}
	}
}
