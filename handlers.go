package main

import (
	"fmt"
	"time"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/codegangsta/negroni"
	"github.com/bitly/go-simplejson"
	"github.com/julienschmidt/httprouter"
	_"golang.org/x/crypto/bcrypt"
	_"github.com/SermoDigital/jose/jws"
	_"github.com/dgrijalva/jwt-go/request"
)

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
	if err != nil{
		panic(err)
	}
	name, _ := js.Get("Name").String()
	email, _ := js.Get("Email").String()
	password := js.Get("Password").MustString()
	//password, err := bcrypt.GenerateFromPassword([]byte(js.Get("Password").MustString()), bcrypt.DefaultCost)
	//fmt.Println(name, email, password)
	//if err != nil{
	//	panic("Password encryption went wrong.")
	//}
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
	Json_Response(response, w)

	//OriginalJson, _:= ioutil.ReadAll(r.Body)
	//r.Body.Close()
	//js, err := simplejson.NewJson([]byte(OriginalJson))
	//if err != nil{
	//	panic(err)
	//}
	//Email, _ := js.Get("Email").String()
	//Password :=js.Get("Password").MustString()

	//err = bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(Password))
	//if err != nil {
	//	fmt.Println("pw wrong")
	//	panic(err)}

}

func UserPost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprint(w, "UserPost!\n")
	negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(ProtectedHandler)),
	)
	//user := &User{}
	//writeOKResponse(w, user)
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

////Populates a model from the params in the Handler
//func populateModelFromHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params, model interface{}) error {
//	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
//	if err != nil {
//		return err
//	}
//	if err := r.Body.Close(); err != nil {
//		return err
//	}
//	if err := json.Unmarshal(body, model); err != nil {
//		return err
//	}
//	return nil
//}
