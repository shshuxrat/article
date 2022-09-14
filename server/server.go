package main

import (
	"article/models"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

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
	var resp []models.Article
	rows, err := db.Query("SELECT  article.id , article.title, article.body, article.created_at, author.firstname, author.lastname FROM article JOIN author ON article.author_id = author.id ")

	if err != nil {
		log.Panic(err)
	}

	if s != "" {

		for rows.Next() {
			var a models.Article
			err = rows.Scan(&a.ID, &a.Title, &a.Body, &a.CreatedAt, &a.Author.Firstname, &a.Author.Lastname)
			if err != nil {
				log.Panic(err)
			}

			if a.Author.Firstname == s || a.Author.Lastname == s {
				resp = append(resp, a)
			}

		}

	} else {

		for rows.Next() {
			var a models.Article
			err = rows.Scan(&a.ID, &a.Title, &a.Body, &a.CreatedAt, &a.Author.Firstname, &a.Author.Lastname)
			resp = append(resp, a)
			if err != nil {
				log.Panic(err)
			}
		}

	}
	defer rows.Close()

	c.JSON(200, resp)
}

func createArticles(c *gin.Context) {
	var article models.Article
	var next_id int
	var not_found bool = true
	if err := c.BindJSON(&article); err != nil {
		log.Println(err)
		c.JSON(422, err.Error())
		return
	}

	t := time.Now()

	//there search

	rows, err := db.Query("SELECT a.id, a.firstname, a.lastname FROM author AS a")

	if err != nil {
		log.Panic(err)
	}
	next_id = 0
	for rows.Next() {
		var a models.Article

		err = rows.Scan(&a.ID, &a.Author.Firstname, &a.Author.Lastname)
		if err != nil {
			log.Panic(err)
		}

		if a.Author.Firstname == article.Author.Firstname && a.Author.Lastname == article.Author.Lastname {
			next_id = a.ID
			not_found = false
			break
		} else {
			if next_id < a.ID {
				next_id = a.ID
			}
		}

	}

	//end search

	if not_found {
		resp2, err2 := db.NamedExec(
			`INSERT INTO author(id,firstname,lastname,created_at) VALUES (:i,:fn,:ln,:t_at)`,
			map[string]interface{}{
				"i":    next_id,
				"fn":   article.Author.Firstname,
				"ln":   article.Author.Lastname,
				"t_at": t,
			},
		)

		if err2 != nil {
			log.Panic(err2)
		}

		fmt.Printf("%#v\n", resp2)

	}
	t = time.Now()

	resp, err5 := db.NamedExec(
		`INSERT INTO article (title, body,author_id,created_at) VALUES (:t,:b, :a_id,:t_at)`,
		map[string]interface{}{
			"t":    article.Title,
			"b":    article.Body,
			"a_id": next_id,
			"t_at": t,
		},
	)
	next_id++
	if err5 != nil {
		log.Panic(err5)
	}

	fmt.Printf("%#v\n", resp)

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
