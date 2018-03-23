package main

import (
	"log"
	"fmt"
	"time"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/bitly/go-simplejson"
	"github.com/julienschmidt/httprouter"
	"github.com/dgrijalva/jwt-go/request"
)


// A map to store the Posts with the ID as the key acts as the storage in lieu of an actual database.
var postDB = make(map[string]*Post)

const SecretKey  = "Vera is testing Ya"

func PostShow(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	fmt.Fprint(w, "Welcome to MockTwitter!\n")
	var Post []*Post
	posts := db.Find(&Post)
	for _, posts := range postDB {
		Post = append(Post, posts)
	}
	println(posts)
	writeOKResponse(w, &Post)
}

func UserCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "UserCreate!\n")
	OriginalJson, _:= ioutil.ReadAll(r.Body)
	r.Body.Close()
	js, err := simplejson.NewJson([]byte(OriginalJson))
	fatal(err)
	name, _ := js.Get("Name").String()
	email, _ := js.Get("Email").String()
	password := js.Get("Password").MustString()
	db.Create(&User{Name: name, Email: email, Password: string(password)})
}

func UserLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	fmt.Fprint(w, "UserLogin!\n")

	var Credential UserCredentials
	err := json.NewDecoder(r.Body).Decode(&Credential)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	var User User
	encodePW := db.Where("Email = ?", Credential.Email).First(&User)
	if encodePW.Error!=nil {
		panic(err)
	}
	if strings.ToLower(Credential.Email) != User.Email {
		if Credential.Password != User.Password {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error logging in")
			fmt.Fprint(w, "Invalid credentials")
			return
		}
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(2)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error extracting the key")
		fatal(err)
	}

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		fatal(err)
	}

	response := Token{tokenString}
	writeOKResponse(w, response)
}

func UserPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "UserPost!\n")

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})
	if err == nil {
		if token.Valid {
			response := Response{"Gained access to protected resource"}
			OriginalJson, _:= ioutil.ReadAll(r.Body)
			r.Body.Close()
			js, err := simplejson.NewJson([]byte(OriginalJson))
			fatal(err)
			post, _ := js.Get("Post").String()
			userid, _ := js.Get("User_refer").String()
			db.Create(&Post{Post: post, UserRefer: userid})
			writeOKResponse(w, response)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, err)
	}
}

// Writes the response as a standard JSON response with StatusOK
func writeOKResponse(w http.ResponseWriter, m interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&JsonResponse{Data: m}); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

// Writes the error response as a Standard API JSON response with a response code
func writeErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: errorCode, Title: errorMsg}})
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
