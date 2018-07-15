package main

import (
	"log"
	"net/http"

	"./routes"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {

	r := mux.NewRouter()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.HandleFunc("/", routes.HomeHandler).Methods("GET")
	r.HandleFunc("/signup", routes.SignupHandler).Methods("GET")
	r.HandleFunc("/login", routes.LoginHandler).Methods("GET")
	r.HandleFunc("/authenticate", routes.AuthenticateHandler).Methods("POST")
	r.HandleFunc("/register", routes.RegisterHandler).Methods("POST")

	r.HandleFunc("/todo", routes.TodoHandler).Methods("GET")

	r.Handle("/{user}/gettodo", negroni.New(
		negroni.HandlerFunc(routes.AuthMiddleware),
		negroni.HandlerFunc(routes.GetTodo),
	))
	r.Handle("/updatetodo/{todo}", negroni.New(
		negroni.HandlerFunc(routes.AuthMiddleware),
		negroni.HandlerFunc(routes.UpdateTodoHandler),
	))

	http.Handle("/", r)

	log.Print("Listening on PORT 8000 aye")
	http.ListenAndServe(":8000", nil)
}
