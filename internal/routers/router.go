package routers

import (
	"kong-anime-go/internal/api/anime"
	"kong-anime-go/internal/api/category"
	"kong-anime-go/internal/api/ping"
	"kong-anime-go/internal/api/tag"
	"kong-anime-go/internal/dao"
	animesrv "kong-anime-go/internal/services/anime"
	categorysrv "kong-anime-go/internal/services/category"
	pingsrv "kong-anime-go/internal/services/ping"
	tagsrv "kong-anime-go/internal/services/tag"

	"kong-anime-go/internal/middleware" // 添加跨域中间件的导入

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// 添加跨域中间件
	router.Use(middleware.CORSMiddleware())

	// Ping
	pingSrv := pingsrv.NewPingService()
	pingHandler := ping.NewPingHandler(pingSrv)

	// Anime
	animeDAO := dao.NewAnimeDAO(db)
	categoryDAO := dao.NewCategoryDAO(db)
	tagDAO := dao.NewTagDAO(db)
	animeSrv := animesrv.NewAnimeService(animeDAO, categoryDAO, tagDAO)
	animeHandler := anime.NewAnimeHandler(animeSrv)

	// Category
	categorySrv := categorysrv.NewCategoryService(categoryDAO)
	categoryHandler := category.NewCategoryHandler(categorySrv)

	// Tag
	tagSrv := tagsrv.NewTagService(tagDAO)
	tagHandler := tag.NewTagHandler(tagSrv)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/hello", pingHandler.GetHello)
		v1.GET("/ping", pingHandler.GetPing)

		// Anime
		v1.POST("/animes", animeHandler.CreateAnime)
		v1.DELETE("/animes/:id", animeHandler.DeleteAnime)
		v1.PUT("/animes/:id", animeHandler.UpdateAnime)
		v1.GET("/animes/:id", animeHandler.GetAnimeByID)
		v1.GET("/animes", animeHandler.GetAllAnimes)
		v1.GET("/animes/search", animeHandler.GetAnimesByName)
		v1.GET("/animes/season", animeHandler.GetAnimesBySeason)
		v1.GET("/animes/category", animeHandler.GetAnimesByCategory)
		v1.GET("/animes/tag", animeHandler.GetAnimesByTag)
		v1.PATCH("/animes/:id/categories", animeHandler.AddCategoriesToAnime)
		v1.PATCH("/animes/:id/tags", animeHandler.AddTagsToAnime)
		v1.GET("/animes/seasons", animeHandler.GetAllSeasons)

		// Category
		v1.POST("/categories", categoryHandler.CreateCategory)
		v1.DELETE("/categories/:id", categoryHandler.DeleteCategory)
		v1.PUT("/categories/:id", categoryHandler.UpdateCategory)
		v1.GET("/categories/:id", categoryHandler.GetCategoryByID)
		v1.GET("/categories", categoryHandler.GetAllCategories)
		v1.GET("/categories/search", categoryHandler.GetCategoriesByName)
		v1.GET("/categories/stats", categoryHandler.GetCategoryStats)

		// Tag
		v1.POST("/tags", tagHandler.CreateTag)
		v1.DELETE("/tags/:id", tagHandler.DeleteTag)
		v1.PUT("/tags/:id", tagHandler.UpdateTag)
		v1.GET("/tags/:id", tagHandler.GetTagByID)
		v1.GET("/tags", tagHandler.GetAllTags)
		v1.GET("/tags/search", tagHandler.GetTagsByName)
		v1.GET("/tags/stats", tagHandler.GetTagStats)
	}

	return router
}
