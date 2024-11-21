package routers

import (
	"kong-anime-go/internal/api/anime"
	"kong-anime-go/internal/api/ping"
	animesrv "kong-anime-go/internal/services/anime"
	pingsrv "kong-anime-go/internal/services/ping"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Ping
	pingSrv := pingsrv.NewPingService()
	pingHandler := ping.NewPingHandler(pingSrv)

	// Anime
	animeSrv := animesrv.NewAnimeService()
	animeHandler := anime.NewAnimeHandler(animeSrv)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/hello", pingHandler.GetHello)
		v1.GET("/ping", pingHandler.GetPing)

		// Anime
		v1.GET("/anime/test", animeHandler.GetTest)
	}

	return router
}
