package postgres

import (
	"article/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type articleRepo struct {
	db *sqlx.DB
}

func NewArticleRepo(db *sqlx.DB) articleRepo {
	return articleRepo{
		db: db,
	}
}

func (r articleRepo) Create(entity models.ArticleCreateModel) (int64, error) {
	resp, err := r.db.NamedExec(
		`INSERT INTO article (title, body, author_id) VALUES (:t,:b, :a_id)`,
		map[string]interface{}{
			"t":    entity.Title,
			"b":    entity.Body,
			"a_id": entity.AuthorID,
		},
	)

	if err != nil {
		return 0, err
	}

	affect, err := resp.RowsAffected()

	if err != nil {
		return 0, err
	}
	return affect, nil

}

func (r articleRepo) GetList(query models.Query) ([]models.Article, error) {
	var resp []models.Article
	var rows *sql.Rows
	var err error

	if len(query.Search) > 0 {
		rows, err = r.db.Query(
			`SELECT  ar.id , ar.title, ar.body, ar.created_at, au.id, au.firstname, au.lastname 
			FROM article AS ar JOIN author  AS au
			ON ar.author_id = au.id 
			WHERE ar.title ILIKE '%' || $3 || '%' 
			OFFSET $1 LIMIT $2`,
			query.Offset,
			query.Limit,
			query.Search,
		)
		if err != nil {
			return resp, err
		}

	} else {

		rows, err = r.db.Query(
			`SELECT  ar.id , ar.title, ar.body, ar.created_at, au.id, au.firstname, au.lastname 
			FROM article AS ar JOIN author  AS au
			ON ar.author_id = au.id 
			WHERE ar.title ILIKE '%' || $3 || '%' 
			OFFSET $1 LIMIT $2`,
			query.Offset,
			query.Limit,
			query.Search,
		)

		if err != nil {
			return resp, err
		}
	}
	defer rows.Close()

	for rows.Next() {

		var a models.Article

		err = rows.Scan(&a.ID, &a.Title, &a.Body, &a.CreatedAt, &a.Author.Id, &a.Author.Firstname, &a.Author.Lastname)

		resp = append(resp, a)

		if err != nil {
			return resp, err
		}
	}

	return resp, nil

}

func (r articleRepo) GetByID(Id int) (article models.Article, err error) {
	var a models.Article
	rows, errId := r.db.NamedQuery(
		`SELECT  ar.id , ar.title, ar.body, ar.created_at, au.id, au.firstname, au.lastname 
		FROM article AS ar
		JOIN author AS au ON ar.author_id = au.id  
		WHERE ar.id = :fn`,
		map[string]interface{}{"fn": Id},
	)

	if errId != nil {
		return a, errId
	}

	rows.Next()
	errId = rows.Scan(&a.ID, &a.Title, &a.Body, &a.CreatedAt, &a.Author.Id, &a.Author.Firstname, &a.Author.Lastname)

	if errId != nil {
		return a, errId
	}

	return a, nil
}

func (r articleRepo) Update(entity models.ArticleUpdateModel) (int64, error) {

	resp, err2 := r.db.Exec(
		"UPDATE article SET title=$1 , body=$2, author_id=$4 ,updated_at=now() where id=$3",
		entity.Title,
		entity.Body,
		entity.ID,
		entity.AuthorID,
	)

	if err2 != nil {

		return 0, err2
	}

	affect, err := resp.RowsAffected()

	if err != nil {

		return 0, err
	}

	return affect, nil

}

func (r articleRepo) Delete(Id int) (int64, error) {

	resp, err := r.db.Exec(
		`DELETE FROM article WHERE id = $1;`,
		Id,
	)

	if err != nil {

		return 0, err
	}

	affect, err := resp.RowsAffected()

	if err != nil {

		return 0, err
	}

	return affect, nil

}
