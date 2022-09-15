package main

import (
	"article/models"
	"article/storage/postgres"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB
var articleRepo postgres.ArticleRepoI

func main() {

	psqlConnString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable ",
		"localhost",
		5432,
		"postgres",
		"1506",
		"bootcamp",
	)
	var err error
	db, err = sqlx.Connect("postgres", psqlConnString)

	if err != nil {
		log.Panic(err)
	}

	articleRepo = postgres.NewArticleRepo(db)

	router := gin.Default()

	router.POST("/articles", createArticles)
	router.GET("/articles", getArticles)
	router.GET("/getid", getArticlesById)
	router.PUT("/update", updateArticle)
	router.DELETE("delete", deleteArticle)

	router.Run(":8083")
}

func getArticles(c *gin.Context) {

	s := c.Query("search")

	if s != "" {
		resp, err := articleRepo.Search(models.Query{Offset: 0, Limit: 10, Search: s})
		if err != nil {
			log.Println(err)
			c.JSON(400, err.Error())
			return
		}
		c.JSON(200, resp)

	} else {

		resp, err := articleRepo.GetList(models.Query{Offset: 0, Limit: 10, Search: ""})

		if err != nil {
			log.Println(err)
			c.JSON(400, err.Error())
			return
		}

		c.JSON(200, resp)

	}

}

func createArticles(c *gin.Context) {
	var article models.Article

	if err := c.BindJSON(&article); err != nil {
		log.Println(err)
		c.JSON(422, err.Error())
		return
	}

	err := articleRepo.Create(article)
	if err != nil {
		log.Println(err)
		c.JSON(422, err.Error())
		return
	}

	c.JSON(201, "success")

}

func getArticlesById(c *gin.Context) {

	i := c.Query("id")

	if i != "" {
		var id int
		var errId error
		if id, errId = strconv.Atoi(i); errId != nil {
			log.Println(errId)
			c.JSON(500, errId.Error())
			return
		}

		a, err := articleRepo.GetByID(id)

		if err != nil {
			log.Println(err)
			c.JSON(404, err.Error())
			return
		}

		c.JSON(200, a)
		return
	}
	c.JSON(404, "not found")
}

func updateArticle(c *gin.Context) {
	var article models.Article
	if err := c.BindJSON(&article); err != nil {
		log.Println(err)
		c.JSON(422, err.Error())
		return
	}

	afferctedRaw, err := articleRepo.Update(article)
	if err != nil {
		log.Println(err)
		c.JSON(422, err.Error)
	}

	log.Println("affected rows : ", afferctedRaw)
	c.JSON(201, "success")

}

func deleteArticle(c *gin.Context) {
	var id int
	if err := c.BindJSON(&id); err != nil {
		log.Println(err)
		c.JSON(422, err.Error())
		return
	}
	afferctedRaw, err := articleRepo.Delete(id)
	if err != nil {
		log.Println(err)
		c.JSON(422, err.Error)
	}

	log.Println("affected rows : ", afferctedRaw)

	c.JSON(200, "deleted")
}
