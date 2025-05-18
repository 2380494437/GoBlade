package routes

import (
	"github.com/gin-gonic/gin"
	"GoBlade/routes/v1"
	"GoBlade/routes/v2"
)

func RegisterRoutes(r *gin.Engine) {
	apiV1 := r.Group("/api/v1")
	v1.Register(apiV1)

	apiV2 := r.Group("/api/v2")
	v2.Register(apiV2)
}
