package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/go-redis/redis"
	"html/template"
)

var client *redis.Client
var templates *template.Template

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "a2a61bdc0208d11eaaabc0a6b8228ff9-2077444820.eu-west-2.elb.amazonaws.com:6379",
	})
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	comments, err := client.LRange("comments", 0, 10).Result()
	if err != nil {
		return
	}
	templates.ExecuteTemplate(w, "index.html", comments)
}
