#  Pipeline a.k.a middleware in Go

Just a playground with some interesting concepts like pipelines aka middleware, handleFuncs, request validations etc. Check it out.


## Pipeline
```go
pipe := pipeline.NewPipeline(1)

yield := pipe.Through(
    func(value interface{}, next pipeline.Handler) {
        next(value.(int)+2, nil)
    },
    func(value interface{}, next pipeline.Handler) {
        next(value.(int)*2, nil)
    },
).Return()

fmt.Printf("yield: %v\n", yield) // yield: 6
```

## Validation middleware 
### Validatable interface
```go
type Validatable interface {
	validate() error
}
```
### Generic http middleware to validate http requests
```go
func Validate(model Validatable, handler func(c *Request)) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
        // ...
    }
}
```

And a coupe of more go goodies. Should you notice a bug or any thing that can be improved send in a PR. Hope it helps someone. 😎