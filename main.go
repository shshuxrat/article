package main

import (
	"article/models"
	"article/storage"
	"fmt"
	"time"
)

func main() {
	articleStorage := storage.NewArticleStorage()

	var a1 models.Article
	a1.ID = 1
	a1.Title = "Lorem"
	a1.Body = "Lorem ipsum"
	var p models.Person = models.Person{
		Firstname: "Adam",
		Lastname:  "Doe",
	}
	a1.Author = p
	t := time.Now()
	a1.CreatedAt = &t

	err := articleStorage.Add(a1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(a1, " is added")
	}

	var a2 models.Article
	a2.ID = 2
	a2.Title = "John"
	a2.Body = "John Junior"
	var p2 models.Person = models.Person{
		Firstname: "John",
		Lastname:  "Terry",
	}
	a2.Author = p2
	t2 := time.Now()
	a2.CreatedAt = &t2

	err2 := articleStorage.Add(a2)
	if err2 != nil {
		fmt.Println(err2)
	} else {
		fmt.Println(a2, " is added")
	}

	var a3 models.Article
	a3.ID = 3
	a3.Title = "Adam"
	a3.Body = "Adam Junior"
	var p3 models.Person = models.Person{
		Firstname: "Adam",
		Lastname:  "John",
	}
	a3.Author = p3
	t3 := time.Now()
	a3.CreatedAt = &t3

	err3 := articleStorage.Add(a3)
	if err3 != nil {
		fmt.Println(err3)
	} else {
		fmt.Println(a3, " is added")
	}

	/*resp, err4 := articleStorage.GetByID(3)

	if err4 != nil {
		fmt.Println(err4)
	} else {
		fmt.Println(resp, " is found ")
	}*/

	//fmt.Println(articleStorage.GetList())

	//fmt.Println(articleStorage.Search("John"))
	//fmt.Println(articleStorage.Search("Junior"))

	/* a3.Title = "new title"
	err5 := articleStorage.Update(a3)
	if err5 != nil {
		fmt.Println(err5)
	} else {
		fmt.Println(a3, "  updated successfully")
	} */
	var id_del int = 3
	err4 := articleStorage.Delete(id_del)
	if err4 != nil {
		fmt.Println(err4)
	} else {
		fmt.Println("id =", id_del, "is deleted")
	}

}
