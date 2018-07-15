package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/gorilla/context"

	"../models"
	"github.com/dgrijalva/jwt-go"
)

var views = map[string]string{
	"index":  "./views/index.html",
	"signup": "./views/signup.html",
	"login":  "./views/login.html",
	"todo":   "./views/todo.html",
}

// Secret is the secret
var Secret = []byte("chocdeveloperSecret")
var jwttoken string
var authUsername string

func toJSON(c interface{}) []byte {
	marshal, _ := json.Marshal(c)
	return marshal
}

// User defines the login user credentials
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Credentials creates the response object to send back to the client
type Credentials struct {
	Username string
	Token    string
}

// HomeHandler handles Request to the homepage, "/"
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	text, err := template.ParseFiles(views["index"])
	if err == nil {
		text.Execute(w, nil)
	}
}

// SignupHandler handles requests to "/signup"
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	text, err := template.ParseFiles(views["signup"])
	if err == nil {
		text.Execute(w, nil)
	}
}

// AuthenticateHandler handles login authentication requests, "/authenticate"
// Also creates the jwt
func AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.Write(toJSON(err))
	}

	var user User
	json.Unmarshal(body, &user)

	_, err = models.Authenticate(user.Username, user.Password)
	if err != nil {
		log.Println(err)
	} else {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"exp":      time.Now().Add(time.Hour * 72).Unix(),
		})
		tokenString, err := token.SignedString(Secret)
		if err != nil {
			log.Println(err)
		}
		w.Write(toJSON(&Credentials{Username: user.Username, Token: tokenString}))
	}
}

// LoginHandler handles Requests to the login page, "/login"
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	text, err := template.ParseFiles(views["login"])
	if err == nil {
		text.Execute(w, nil)
	}
}

// RegisterHandler handles signup registeration requests, "/register"
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostFormValue("username")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	user, err := models.RegisterUser(username, email, password)
	if err != nil {
		log.Print(err)
	} else {
		log.Print(user)
	}
}

// AuthMiddleware authenticates the JWT token
func AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	var UserCred Credentials
	json.Unmarshal(body, &UserCred)

	token, err := jwt.Parse(UserCred.Token, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})
	if err == nil && token.Valid {
		// Create context to send username to other handlers after authentictaion
		context.Set(r, "username", UserCred.Username)
		log.Print(UserCred.Username)
		next(w, r)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized, You should sign in First")
	}
}

// GetTodo get the users todos from the database
func GetTodo(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	username := context.Get(r, "username")

	todos, err := models.ReadTodos(username.(string))
	if err != nil {
		log.Fatal(err)
	}
	w.Write(toJSON(todos))
}

// TodoHandler handles get requests to the todo route
func TodoHandler(w http.ResponseWriter, r *http.Request) {
	text, err := template.ParseFiles(views["todo"])
	if err == nil {
		text.Execute(w, nil)
	}
}

// UpdateTodoHandler updates the value of the users todo
func UpdateTodoHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	log.Print("update")

	tasks := strings.Split(vars["todo"], ",")

	w.WriteHeader(http.StatusOK)

	username := context.Get(r, "username")

	log.Print(username)
	log.Print(tasks)

	status, err := models.SaveTodos(username.(string), tasks)
	jsonErr, _ := json.Marshal(err)
	jsonStatus, _ := json.Marshal(status)
	if err != nil {
		w.Write(jsonErr)
	}
	w.Write(jsonStatus)
}
