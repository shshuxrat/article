package storage

import (
	"article/models"
	"errors"
)

func NewArticleStorage() ArticleStorage {
	return ArticleStorage{
		data: make(map[int]models.Article),
	}
}

type ArticleStorage struct {
	data map[int]models.Article
}

var ErrorAlreadyEXists = errors.New("already exists")
var ErrorNotCreatedYet = errors.New("this entity is notcreated yet")
var ErrorNotFound = errors.New("not found")

func (storage *ArticleStorage) Add(entity models.Article) error {
	if _, ok := storage.data[entity.ID]; ok {
		return ErrorAlreadyEXists
	}
	storage.data[entity.ID] = entity
	return nil
}

func (storage *ArticleStorage) GetByID(ID int) (models.Article, error) {
	var resp models.Article
	var ok bool
	if resp, ok = storage.data[ID]; !ok {
		return resp, ErrorNotFound
	}
	return resp, nil
}

func (storage *ArticleStorage) GetList() []models.Article {
	var resp []models.Article
	var i int = 0
	for _, val := range storage.data {
		resp = append(resp, val)
		i++
	}
	return resp
}

func (storage *ArticleStorage) Search(str string) []models.Article {
	var resp []models.Article
	for _, value := range storage.data {
		if value.Author.Firstname == str || value.Author.Lastname == str {
			resp = append(resp, value)
		}
	}
	return resp
}

func (storage *ArticleStorage) Update(entity models.Article) error {
	resp := storage.data[entity.ID]

	if resp.ID != entity.ID {
		return ErrorNotCreatedYet
	}

	storage.data[entity.ID] = entity

	return nil
}

func (storage *ArticleStorage) Delete(ID int) error {
	resp := storage.data[ID]

	if resp.ID != ID {
		return ErrorNotFound
	}

	delete(storage.data, ID)

	return nil
}
