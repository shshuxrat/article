package postgres

import (
	"article/models"
	"time"

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

func (a authorRepo) Create(entity models.Person) (int64, error) {
	var err error
	var not_found bool = true
	t := time.Now()

	//there search  author if exists or not
	rows, err := a.db.Query("SELECT a.id  FROM author AS a")

	if err != nil {
		return 0, err
	}

	for rows.Next() {
		var a int
		err = rows.Scan(&a)
		if err != nil {
			return 0, err
		}
		if a == entity.Id {
			not_found = false // there we found person already created
			break
		}
	}
	//end search

	if not_found {
		resp, err2 := a.db.NamedExec(
			`INSERT INTO author(id,firstname,lastname,created_at) VALUES (:i,:fn,:ln,:t_at)`,
			map[string]interface{}{
				"i":    entity.Id,
				"fn":   entity.Firstname,
				"ln":   entity.Lastname,
				"t_at": t,
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

	return 0, nil
}

func (a authorRepo) GetList(query models.Query) ([]models.Person, error) {
	var arr []models.Person
	rows, err := a.db.Query(
		`SELECT  a.Id, a.firstname, a.lastname
		FROM author AS a  
		OFFSET $1 LIMIT $2`,
		query.Offset,
		query.Limit,
	)
	if err != nil {
		return arr, err
	}

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
func (a authorRepo) Search(query models.Query) ([]models.Person, error) {
	var arr []models.Person

	rows, err := a.db.Query(
		`SELECT  a.Id , a.firstname, a.lastname 
		FROM author AS a
		OFFSET $1 LIMIT $2`,
		query.Offset,
		query.Limit,
	)

	if err != nil {
		return arr, err
	}
	for rows.Next() {
		var a models.Person

		err = rows.Scan(&a.Id, &a.Firstname, &a.Lastname)

		if err != nil {
			return arr, err
		}

		if a.Firstname == query.Search || a.Lastname == query.Search {
			arr = append(arr, a)
		}
	}

	return arr, nil
}
func (a authorRepo) Update(entity models.Person) (int64, error) {

	resp, err2 := a.db.Exec(
		"UPDATE author SET firstname=$1 , lastname=$2, updated_at=now() where id=$3",
		entity.Firstname,
		entity.Lastname,
		entity.Id,
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
