package postgres

import (
	"article/models"
	"time"

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

func (r articleRepo) Create(entity models.Article) error {
	var err error
	var next_id int
	var not_found bool = true
	t := time.Now()

	//there search person if exists or not if person exists next_id = person.id else next_id = max of all ids

	rows, err := r.db.Query("SELECT a.id, a.firstname, a.lastname FROM author AS a")

	if err != nil {
		return err
	}
	next_id = 0
	for rows.Next() {
		var a models.Article

		err = rows.Scan(&a.ID, &a.Author.Firstname, &a.Author.Lastname)
		if err != nil {
			return err
		}

		if a.Author.Firstname == entity.Author.Firstname && a.Author.Lastname == entity.Author.Lastname {
			next_id = a.ID
			not_found = false // there we found person already created
			break
		} else {
			if next_id < a.ID {
				next_id = a.ID

			}
		}

	}

	//end search

	if not_found {
		next_id++
		_, err2 := r.db.NamedExec(
			`INSERT INTO author(id,firstname,lastname,created_at) VALUES (:i,:fn,:ln,:t_at)`,
			map[string]interface{}{
				"i":    next_id,
				"fn":   entity.Author.Firstname,
				"ln":   entity.Author.Lastname,
				"t_at": t,
			},
		)

		if err2 != nil {
			return err2
		}

	}
	t = time.Now()

	_, err5 := r.db.NamedExec(
		`INSERT INTO article (title, body,author_id,created_at) VALUES (:t,:b, :a_id,:t_at)`,
		map[string]interface{}{
			"t":    entity.Title,
			"b":    entity.Body,
			"a_id": next_id,
			"t_at": t,
		},
	)

	if err5 != nil {
		return err5
	}

	return nil

}

func (r articleRepo) GetList(query models.Query) ([]models.Article, error) {
	var resp []models.Article

	rows, err := r.db.Query(
		"SELECT  article.id , article.title, article.body, article.created_at, author.firstname, author.lastname FROM article JOIN author ON article.author_id = author.id  OFFSET $1 LIMIT $2",
		query.Offset,
		query.Limit,
	)

	if err != nil {
		return resp, err
	}
	for rows.Next() {

		var a models.Article

		err = rows.Scan(&a.ID, &a.Title, &a.Body, &a.CreatedAt, &a.Author.Firstname, &a.Author.Lastname)

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
		`SELECT  article.id , article.title, article.body, article.created_at, author.firstname, author.lastname FROM article JOIN author ON article.author_id = author.id  WHERE article.id = :fn`,
		map[string]interface{}{"fn": Id},
	)

	if errId != nil {
		return a, errId
	}

	rows.Next()
	errId = rows.Scan(&a.ID, &a.Title, &a.Body, &a.CreatedAt, &a.Author.Firstname, &a.Author.Lastname)

	if errId != nil {
		return a, errId
	}

	return a, nil
}

func (r articleRepo) Search(query models.Query) ([]models.Article, error) {

	var resp []models.Article
	rows, err := r.db.Query(
		"SELECT  article.id , article.title, article.body, article.created_at, author.firstname, author.lastname FROM article JOIN author ON article.author_id = author.id  OFFSET $1 LIMIT $2",
		query.Offset,
		query.Limit,
	)

	if err != nil {
		return resp, err
	}
	for rows.Next() {
		var a models.Article

		err1 := rows.Scan(&a.ID, &a.Title, &a.Body, &a.CreatedAt, &a.Author.Firstname, &a.Author.Lastname)

		if err1 != nil {
			return resp, err
		}

		if a.Author.Firstname == query.Search || a.Author.Lastname == query.Search {
			resp = append(resp, a)
		}
	}

	return resp, nil

}

func (r articleRepo) Update(entity models.Article) (int64, error) {

	resp, err2 := r.db.Exec(
		"UPDATE article SET title=$1 , body=$2, updated_at=now() where id=$3",
		entity.Title,
		entity.Body,
		entity.ID,
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
