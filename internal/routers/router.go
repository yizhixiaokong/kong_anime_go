package routers

import (
	"kong-anime-go/internal/api/anime"
	"kong-anime-go/internal/api/category"
	"kong-anime-go/internal/api/follow" // 添加追番API的导入
	"kong-anime-go/internal/api/ping"
	"kong-anime-go/internal/api/tag"
	"kong-anime-go/internal/dao"
	animesrv "kong-anime-go/internal/services/anime"
	categorysrv "kong-anime-go/internal/services/category"
	followsrv "kong-anime-go/internal/services/follow" // 添加追番服务的导入
	pingsrv "kong-anime-go/internal/services/ping"
	tagsrv "kong-anime-go/internal/services/tag"

	"kong-anime-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	// Ping
	pingSrv := pingsrv.NewService()
	pingHandler := ping.NewHandler(pingSrv)

	// Anime
	animeDAO := dao.NewAnimeDAO(db)
	categoryDAO := dao.NewCategoryDAO(db)
	tagDAO := dao.NewTagDAO(db)
	followDAO := dao.NewFollowDAO(db)
	animeSrv := animesrv.NewService(animeDAO, categoryDAO, tagDAO, followDAO)
	animeHandler := anime.NewHandler(animeSrv)

	// Category
	categorySrv := categorysrv.NewService(categoryDAO)
	categoryHandler := category.NewHandler(categorySrv)

	// Tag
	tagSrv := tagsrv.NewService(tagDAO)
	tagHandler := tag.NewHandler(tagSrv)

	// Follow
	followSrv := followsrv.NewService(followDAO, animeDAO)
	followHandler := follow.NewHandler(followSrv)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/hello", pingHandler.GetHello)
		v1.GET("/ping", pingHandler.GetPing)

		// Anime
		v1.POST("/animes", animeHandler.Create)
		v1.DELETE("/animes/:id", animeHandler.Delete)
		v1.PUT("/animes/:id", animeHandler.Update)
		v1.GET("/animes/:id", animeHandler.GetByID)
		v1.GET("/animes", animeHandler.GetAll)
		v1.GET("/animes/search", animeHandler.GetByName)
		v1.GET("/animes/season", animeHandler.GetBySeason)
		v1.GET("/animes/category", animeHandler.GetByCategory)
		v1.GET("/animes/tag", animeHandler.GetByTag)
		v1.PATCH("/animes/:id/categories", animeHandler.AddCategoriesToAnime)
		v1.PATCH("/animes/:id/tags", animeHandler.AddTagsToAnime)
		v1.GET("/animes/seasons", animeHandler.GetAllSeasons)

		// Category
		v1.POST("/categories", categoryHandler.Create)
		v1.DELETE("/categories/:id", categoryHandler.Delete)
		v1.PUT("/categories/:id", categoryHandler.Update)
		v1.GET("/categories/:id", categoryHandler.GetByID)
		v1.GET("/categories", categoryHandler.GetAll)
		v1.GET("/categories/search", categoryHandler.GetByName)
		v1.GET("/categories/stats", categoryHandler.GetStats)

		// Tag
		v1.POST("/tags", tagHandler.Create)
		v1.DELETE("/tags/:id", tagHandler.Delete)
		v1.PUT("/tags/:id", tagHandler.Update)
		v1.GET("/tags/:id", tagHandler.GetByID)
		v1.GET("/tags", tagHandler.GetAll)
		v1.GET("/tags/search", tagHandler.GetByName)
		v1.GET("/tags/stats", tagHandler.GetStats)

		// Follow
		v1.POST("/follows", followHandler.Create)
		v1.DELETE("/follows/:id", followHandler.Delete)
		v1.PUT("/follows/:id", followHandler.Update)
		v1.GET("/follows/:id", followHandler.GetByID)
		v1.GET("/follows", followHandler.GetAll)
		v1.PATCH("/follows/:id/status", followHandler.UpdateStatus)
		v1.GET("/follows/categories", followHandler.GetAllCategories)
	}

	return router
}
