package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-vals/pipeline"
	"net/http"
	"os"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

var emails = []string{"a@b.c"}

// Sample custom validator
func unused(value interface{}) error {
	for _, email := range emails {
		if email == value {
			return errors.New("Email used up")
		}
	}
	return nil
}

type Validatable interface {
	validate() error
}

type AdminUser struct {
	Name  string
	Email string
	Role  string
}

type User struct {
	Name  string
	Email string
}

func (u *User) validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Email,
			validation.Required,
			validation.Length(3, 20),
			validation.By(unused),
			is.Email,
			validation.In("admin@home.org", "admin@example.com").Error("You cannot just use any email you like!"),
		),
	)
}

func Validate(model Validatable, handler func(c *Request)) func(rw http.ResponseWriter, r *http.Request) {

	return func(rw http.ResponseWriter, r *http.Request) {

		request := &Request{rw, r}

		if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
			rw.WriteHeader(http.StatusNotAcceptable)
			request.Json(Json{
				"message": "Your request is not valid.",
				"errors":  err.Error(),
			})
			return
		}

		errs := model.validate()

		if errs != nil {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusNotAcceptable)
			request.Json(Json{
				"message": "Some of the data is not valid, please check what you have provided.",
				"errors":  errs,
			})
		} else {
			handler(request)
		}
	}

}

type Counter int

func (c *Counter) Inc() *Counter {
	*c += 1
	return c
}

type Json = map[string]interface{}

type Request struct {
	http.ResponseWriter
	R *http.Request
}

func (r *Request) Json(resp interface{}) error {
	return json.NewEncoder(r).Encode(resp)
}

var JSON = func(value interface{}, next pipeline.Handler) {
	value.(Request).Header().Add("Content-Type", "application/json")
	next(value, nil)
}

var ValidateStruct = func(u interface{}) pipeline.Handler {
	return func(value interface{}, next pipeline.Handler) {
		if err := json.NewDecoder(value.(Request).R.Body).Decode(u); err != nil {
			println(err.Error())
		}

		err := u.(*User).validate()

		if err != nil {
			value.(Request).WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(value.(Request)).Encode(err)
		} else {
			next(value, nil)
		}

	}
}

func main() {

	http.HandleFunc("/api", func(rw http.ResponseWriter, r *http.Request) {
		request := &Request{rw, r}
		middlewares := pipeline.NewPipeline(request)

		user := User{}
		middlewares.Through(
			JSON,
			ValidateStruct(&user),
			func(value interface{}, next pipeline.Handler) {
				token, ok := value.(Request).R.Header["Authorization"]

				if ok && token[0] == "secure-token" {
					next(value, nil)
					return
				}

				request.Json(Json{
					"message": "You are not authenticated",
				})
				return

			},
			func(value interface{}, next pipeline.Handler) {
				request.Json("Welcome to admin area..")
			})

	})

	http.HandleFunc("/", Validate(&User{}, func(c *Request) {
		c.Json(Json{
			"message": "Welcome to a go-dump!",
		})
	}))

	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "3000"
	}

	fmt.Printf("Listening on port %s", port)
	http.ListenAndServe(":"+port, nil)
}
