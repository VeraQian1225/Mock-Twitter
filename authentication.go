package main

import (
	_ "time"
	_"strconv"
	"io/ioutil"
	_ "encoding/hex"
	_ "crypto/sha512"
	_"github.com/satori/go.uuid"
	_"github.com/SermoDigital/jose/jws"
	_"github.com/SermoDigital/jose/crypto"
	_"github.com/julienschmidt/httprouter"
	"github.com/dgrijalva/jwt-go/request"
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/bitly/go-simplejson"
)


const SecretKey  = "Vera is testing Haha"


func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

	if err == nil {
		if token.Valid {
			next(w, r)
			OriginalJson, _:= ioutil.ReadAll(r.Body)
			r.Body.Close()
			js, err := simplejson.NewJson([]byte(OriginalJson))
			if err != nil{
				panic(err)
			}
			post, _ := js.Get("Post").String()
			userid, _ := js.Get("User_refer").String()
			db.Create(&Post{Post: post, UserRefer: userid})
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, err)
	}
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{"Gained access to protected resource"}
	Json_Response(response, w)
}

func Json_Response(response interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}