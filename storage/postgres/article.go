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

func (r articleRepo) Create(entity models.Article) (int64, error) {
	var err error
	var not_found bool = true
	t := time.Now()

	//there search  author if exists or not
	rows, err := r.db.Query("SELECT a.id  FROM author AS a")

	if err != nil {
		return 0, err
	}

	for rows.Next() {
		var a int
		err = rows.Scan(&a)
		if err != nil {
			return 0, err
		}
		if a == entity.Author.Id {
			not_found = false // there we found person already created
			break
		}
	}
	defer rows.Close()
	//end search

	if not_found {
		_, err2 := r.db.NamedExec(
			`INSERT INTO author(id,firstname,lastname,created_at) VALUES (:i,:fn,:ln,:t_at)`,
			map[string]interface{}{
				"i":    entity.Author.Id,
				"fn":   entity.Author.Firstname,
				"ln":   entity.Author.Lastname,
				"t_at": t,
			},
		)

		if err2 != nil {
			return 0, err2
		}

	}
	t = time.Now()
	//there search  article if exists or not
	rows, err = r.db.Query("SELECT a.id FROM article AS a")
	not_found = true
	if err != nil {
		return 0, err
	}

	for rows.Next() {
		var a int
		err = rows.Scan(&a)
		if err != nil {
			return 0, err
		}
		if a == entity.ID {
			not_found = false // there we found person already created
			break
		}
	}
	defer rows.Close()
	//end search

	if not_found {
		resp, err5 := r.db.NamedExec(
			`INSERT INTO article (title, body,author_id,created_at) VALUES (:t,:b, :a_id,:t_at)`,
			map[string]interface{}{
				"t":    entity.Title,
				"b":    entity.Body,
				"a_id": entity.Author.Id,
				"t_at": t,
			},
		)

		if err5 != nil {
			return 0, err5
		}

		affect, err := resp.RowsAffected()

		if err != nil {
			return 0, err
		}
		return affect, nil
	}

	return 0, nil

}

func (r articleRepo) GetList(query models.Query) ([]models.Article, error) {
	var resp []models.Article

	rows, err := r.db.Query(
		`SELECT  article.id , article.title, article.body, article.created_at, author.id, author.firstname, author.lastname 
		FROM article JOIN author ON article.author_id = author.id  
		OFFSET $1 LIMIT $2`,
		query.Offset,
		query.Limit,
	)

	if err != nil {
		return resp, err
	}
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

func (r articleRepo) Search(query models.Query) ([]models.Article, error) {

	var resp []models.Article
	rows, err := r.db.Query(
		`SELECT  ar.id , ar.title, ar.body, ar.created_at, au.id, au.firstname, au.lastname 
		FROM article AS ar
		JOIN author AS au ON ar.author_id = au.id  
		OFFSET $1 LIMIT $2`,
		query.Offset,
		query.Limit,
	)

	if err != nil {
		return resp, err
	}
	for rows.Next() {
		var a models.Article

		err1 := rows.Scan(&a.ID, &a.Title, &a.Body, &a.CreatedAt, &a.Author.Id, &a.Author.Firstname, &a.Author.Lastname)

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
