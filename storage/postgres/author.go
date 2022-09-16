package postgres

import (
	"article/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type authorRepo struct {
	db *sqlx.DB
}

func NewAuthorRepo(db *sqlx.DB) authorRepo {
	return authorRepo{
		db: db,
	}
}

func (a authorRepo) Create(entity models.PersonCreateModel) (int64, error) {

	resp, err2 := a.db.NamedExec(
		`INSERT INTO author(firstname,lastname) 
		VALUES (:fn,:ln)`,
		map[string]interface{}{
			"fn": entity.Firstname,
			"ln": entity.Lastname,
		},
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

func (a authorRepo) GetList(query models.Query) ([]models.Person, error) {
	var arr []models.Person
	var rows *sql.Rows
	var err error
	if len(query.Search) > 0 {
		rows, err = a.db.Query(
			`SELECT  a.Id, a.firstname, a.lastname
			FROM author AS a
			WHERE a.firstname ILIKE '%' || $3 || '%' OR a.lastname ILIKE '%' || $3 || '%'
			OFFSET $1 LIMIT $2`,
			query.Offset,
			query.Limit,
			query.Search,
		)
		if err != nil {
			return arr, err
		}

	} else {
		rows, err = a.db.Query(
			`SELECT  a.Id, a.firstname, a.lastname
			FROM author AS a
			OFFSET $1 LIMIT $2`,
			query.Offset,
			query.Limit,
		)
		if err != nil {
			return arr, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Person
		err = rows.Scan(&p.Id, &p.Firstname, &p.Lastname)
		if err != nil {
			return arr, err
		}
		arr = append(arr, p)
	}

	return arr, nil
}
func (a authorRepo) GetByID(Id int) (models.Person, error) {
	var p models.Person

	rows, err := a.db.NamedQuery(
		`SELECT  a.id , a.firstname, a.lastname 
		FROM author AS a
		WHERE a.id = :fn`,
		map[string]interface{}{"fn": Id},
	)

	if err != nil {
		return p, err
	}

	rows.Next()
	err = rows.Scan(&p.Id, &p.Firstname, &p.Lastname)

	if err != nil {
		return p, err
	}

	return p, nil
}

func (a authorRepo) Update(entity models.PersonUpdateModel) (int64, error) {

	resp, err2 := a.db.Exec(
		"UPDATE author SET firstname=$1 , lastname=$2, updated_at=now() where id=$3",
		entity.Firstname,
		entity.Lastname,
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

func (a authorRepo) Delete(ID int) (int64, error) {
	var number int64
	return number, nil
}
