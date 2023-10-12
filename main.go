package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type ArticleList struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts = []ArticleList{}
var showPost = ArticleList{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html",
		"templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	//Select all articles
	res, err := db.Query("Select * from `articles` limit 5")
	if err != nil {
		panic(err)
	}
	posts = []ArticleList{}
	for res.Next() {
		var post ArticleList
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)

	}

	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html",
		"templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || title == " " || anons == "" || anons == " " || full_text == "" || full_text == " " {
		fmt.Fprintf(w, "Incorrect data")
	} else {

		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
		if err != nil {
			panic(err)
		}

		defer db.Close()

		//Add an article to MySQL DB
		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`, `anons`, `full_text`)"+
			"VALUES ('%s', '%s', '%s')", title, anons, full_text))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func show_post(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/show.html", "templates/header.html",
		"templates/footer.html")
	if err != nil {
		log.Println("No template", err)
	}
	vars := mux.Vars(r)

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//Select an article
	res, err := db.Query(fmt.Sprintf("Select * from `articles` where `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}
	showPost = ArticleList{}
	for res.Next() {
		var post ArticleList
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		showPost = post

	}

	t.ExecuteTemplate(w, "show", showPost)
}

func contacts(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/contacts.html", "templates/header.html",
		"templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "contacts", nil)
}

func handleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")
	rtr.HandleFunc("/contacts/", contacts).Methods("GET")

	http.Handle("/", rtr)

	if er := http.ListenAndServe(":8080", nil); er != nil {
		log.Println(er)
	}
	rtr.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
}

func main() {
	handleFunc()
}
