package models

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// User model for database
type User struct {
	Username string
	Email    string
	Password string
	Todos    []string
}

// RegisterUser creates the user and stores it on the database
func RegisterUser(username string, email string, password string) (string, error) {
	result := User{}
	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	// Establish Database connection
	session, err := mgo.Dial("localhost")
	if err != nil {
		return "", err
	}
	defer session.Close()
	c := session.DB("choctodo").C("users")

	// Check if username already exists
	// if User exists return error
	err = c.Find(bson.M{"username": username}).One(&result)
	if err == nil {
		userExists := errors.New("Username already exists")
		return "", userExists
	}

	// Create a user and add to the database
	err = c.Insert(&User{username, email, string(hashedPassword), []string{}})
	if err != nil {
		return "", err
	}
	return "Registeration Successful", nil
}

// Authenticate checks the database to see if the username and password are correct
func Authenticate(username string, password string) (interface{}, error) {
	result := User{}

	// Establish database connection
	session, err := mgo.Dial("localhost")
	if err != nil {
		// log.Fatal(err)
		return nil, err
	}
	defer session.Close()
	c := session.DB("choctodo").C("users")

	// Query the database with the user and compare the password with the hashed password
	err = c.Find(bson.M{"username": username}).One(&result)
	if err != nil {
		// log.Fatal(err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))
	if err != nil {
		// log.Fatal(err)
		return nil, err
	}

	return result, nil
	// compare username password to entered password
	// if password is true return user and session
}

// ReadTodos finds the username and returns only the todos
func ReadTodos(username string) ([]string, error) {
	result := User{}
	// Establish database connection
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer session.Close()
	c := session.DB("choctodo").C("users")

	err = c.Find(bson.M{"username": username}).One(&result)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return result.Todos, nil
}

// SaveTodos saves todos to the database
func SaveTodos(username string, update []string) (string, error) {

	// Establish database connection
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Print(err)
	}
	defer session.Close()
	c := session.DB("choctodo").C("users")

	query := bson.M{"username": username}
	changeQuery := bson.M{"$set": bson.M{"todos": update}}

	err = c.Update(query, changeQuery)
	if err != nil {
		return "", err
	}
	return "Todos has been updated", nil
}
