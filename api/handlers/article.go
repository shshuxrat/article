package handlers

import (
	"article/models"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetArticles(c *gin.Context) {

	s := c.Query("search")

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		log.Println(err)
		c.JSON(400, err.Error())
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Println(err)
		c.JSON(400, err.Error())
	}

	if s != "" {
		resp, err := h.strg.Article().Search(models.Query{Offset: offset, Limit: limit, Search: s})
		if err != nil {
			log.Println(err)
			c.JSON(400, err.Error())
			return
		}
		c.JSON(200, resp)

	} else {

		resp, err := h.strg.Article().GetList(models.Query{Offset: offset, Limit: limit})

		if err != nil {
			log.Println(err)
			c.JSON(400, err.Error())
			return
		}

		c.JSON(200, resp)

	}

}

func (h *Handler) CreateArticle(c *gin.Context) {
	var article models.Article

	if err := c.BindJSON(&article); err != nil {
		log.Println(err)
		c.JSON(422, err.Error())
		return
	}

	rowsAffected, err := h.strg.Article().Create(article)

	if err != nil {
		log.Println(err)
		c.JSON(422, err.Error())
		return
	}

	if rowsAffected == 0 {
		c.JSON(400, "already exist id")
		return
	}

	c.JSON(201, "success")

}

func (h *Handler) GetArticlesById(c *gin.Context) {

	i := c.Query("id")

	if i != "" {
		var id int
		var errId error
		if id, errId = strconv.Atoi(i); errId != nil {
			log.Println(errId)
			c.JSON(500, errId.Error())
			return
		}

		a, err := h.strg.Article().GetByID(id)

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

func (h *Handler) UpdateArticle(c *gin.Context) {
	var article models.Article
	if err := c.BindJSON(&article); err != nil {
		log.Println(err)
		c.JSON(422, err.Error())
		return
	}

	afferctedRaw, err := h.strg.Article().Update(article)
	if err != nil {
		log.Println(err)
		c.JSON(422, err.Error)
	}

	log.Println("affected rows : ", afferctedRaw)
	c.JSON(201, "success")

}

func (h *Handler) DeleteArticle(c *gin.Context) {
	var id int
	if err := c.BindJSON(&id); err != nil {
		log.Println(err)
		c.JSON(422, err.Error())
		return
	}
	afferctedRaw, err := h.strg.Article().Delete(id)
	if err != nil {
		log.Println(err)
		c.JSON(422, err.Error)
	}

	log.Println("affected rows : ", afferctedRaw)

	c.JSON(200, "deleted")
}
