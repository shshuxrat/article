package storage

import (
	"article/models"
)

type StorageI interface {
	Article() ArticleRepoI
	Author() AuthorRepoI
}

type ArticleRepoI interface {
	Create(entity models.ArticleCreateModel) (int64, error)
	GetList(query models.Query) ([]models.Article, error)
	GetByID(Id int) (models.Article, error)

	Update(entity models.ArticleUpdateModel) (int64, error)
	Delete(ID int) (int64, error)
}

type AuthorRepoI interface {
	Create(entity models.PersonCreateModel) (int64, error)
	GetList(query models.Query) ([]models.Person, error)
	GetByID(Id int) (models.Person, error)

	Update(entity models.PersonUpdateModel) (int64, error)
	Delete(ID int) (int64, error)
}
