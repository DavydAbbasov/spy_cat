package swagger

import (
	_ "github.com/DavydAbbasov/spy-cat/docs"
	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func Swagger() gin.HandlerFunc {
	handler := httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DomID("swagger-ui"),
	)

	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
		
	}
}
