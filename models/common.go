package models

type Query struct {
	Offset int    `json:"offset" default:"0"`
	Limit  int    `json:"limit" default:"10"`
	Search string `json:"search"`
}
