package main

import "github.com/jinzhu/gorm"

type UserCredentials struct {
	Email string 		`json:"email"`
	Password string 	`json:"password"`
}

type Response struct {
	Data string 		`json:"data"`
}

type Token struct {
	Token string 		`json:"token"`
}

//User database
type User struct {
	gorm.Model
	Name      string	`json:"name"`
	Email     string	`json:"email"`
	Password  string	`json:"password"`
}

//Post database
type Post struct {
	gorm.Model
	Post      string	`json:"post"`
	User      User 		`gorm:"foreignkey:UserRefer"` // use UserRefer as foreign key
	UserRefer string
}

//Post database
type JsonResponse struct {
	// Reserved field to add some meta information to the API response
	Meta interface{} 	`json:"meta"`
	Data interface{} 	`json:"data"`
}

type JsonErrorResponse struct {
	Error *ApiError 	`json:"error"`
}

type ApiError struct {
	Status int    		`json:"status"`
	Title  string 		`json:"title"`
}
