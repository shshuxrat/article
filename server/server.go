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

		resp, err := articleRepo.GetList(models.Query{Offset: 0, Limit: 10})

		if err != nil {
			log.Println(err)
			c.JSON(400, err.Error())
			return
		}

		c.JSON(200, resp)

	}
	c.JSON(200, "success")

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

		var a models.Article
		rows, errId := db.NamedQuery(
			`SELECT  article.id , article.title, article.body, article.created_at, author.firstname, author.lastname FROM article JOIN author ON article.author_id = author.id  WHERE article.id = :fn`,
			map[string]interface{}{"fn": id},
		)

		if errId != nil {
			log.Println(errId)
			c.JSON(500, errId.Error())

		}
		rows.Next()
		errId = rows.Scan(&a.ID, &a.Title, &a.Body, &a.CreatedAt, &a.Author.Firstname, &a.Author.Lastname)
		if errId != nil {
			log.Println(errId)
			c.JSON(500, errId.Error())

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

	_, err2 := db.Exec(
		"UPDATE article SET title=$1 , body=$2, updated_at=now() where id=$3",
		article.Title,
		article.Body,
		article.ID,
	)

	if err2 != nil {
		log.Panic(err2)
		c.JSON(404, "not found")
		return
	}

	c.JSON(201, "success")

}

func deleteArticle(c *gin.Context) {
	var id int
	if err := c.BindJSON(&id); err != nil {
		log.Println(err)
		c.JSON(422, err.Error())
		return
	}

	_, errdel := db.Exec(
		`DELETE FROM article WHERE id = $1;`,
		id,
	)

	if errdel != nil {
		log.Println(errdel)
		c.JSON(404, errdel.Error())
	}
	c.JSON(200, "deleted")
}
