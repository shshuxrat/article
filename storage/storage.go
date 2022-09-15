package storage

import (
	"article/models"
)

type ArticleRepoI interface {
	Create(entity models.Article) error
	GetList(query models.Query) ([]models.Article, error)
	GetByID(Id int) (models.Article, error)
	Search(query models.Query) ([]models.Article, error)
	Update(entity models.Article) (int64, error)
	Delete(ID int) (int64, error)
}
