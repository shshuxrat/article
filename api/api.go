package api

import (
	"article/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetUpAPI(r *gin.Engine, h handlers.Handler) {

	r.POST("/articles", h.CreateArticle)
	r.GET("/articles", h.GetArticles) //search ant getAll
	r.GET("/getid", h.GetArticlesById)
	r.PUT("/update", h.UpdateArticle)
	r.DELETE("delete", h.DeleteArticle)

	r.POST("/author", h.CreateAuthor)
	r.GET("/authors", h.GetAuthors) //search ant getAll
	r.GET("/getidp", h.GetAuthorById)
	r.PUT("updatep", h.UpdateAuthor)

}
