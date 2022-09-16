package models

type Person struct {
	Id        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type PersonCreateModel struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type PersonUpdateModel struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
