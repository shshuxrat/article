package handlers

import "article/storage"

type Handler struct {
	strg storage.StorageI
}

func NewHandler(strg storage.StorageI) Handler {
	return Handler{
		strg: strg,
	}
}
