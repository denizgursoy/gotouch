# Go template examples

## If, else if, else 
You can use following keywords in if statements to compare values

```
eq
	Returns the boolean truth of arg1 == arg2
ne
	Returns the boolean truth of arg1 != arg2
lt
	Returns the boolean truth of arg1 < arg2
le
	Returns the boolean truth of arg1 <= arg2
gt
	Returns the boolean truth of arg1 > arg2
ge
	Returns the boolean truth of arg1 >= arg2
```

Values:
```yaml
webFramework: default
```

File:
```go
func main(){
{{ if eq  .webFramework "echo" }}
    e := echo.New()
    
    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })
    e.Logger.Fatal(e.Start(":8080"))
{{else if eq .httpLibrary "gin"}}
    r := gin.Default()
    
    r.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "Hello, World!")
    })
    r.Run(":8080")
{{ else }}
    http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
    	io.WriteString(writer, "version 1")
    })
    
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err)
    }
{{ end }}
}
```
Output:
```go
func main(){
    http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
    	io.WriteString(writer, "version 1")
    })
    
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err)
    }
}
```

## If to chek boolean
Values:
```yaml
isEcho: true
```
File:
```go
func main(){
{{ if .isEcho }}
    e := echo.New()
    
    e.GET("/", func(c echo.Context) error {
    	return c.String(http.StatusOK, "Hello, World!")
    })
    e.Logger.Fatal(e.Start(":8080"))
{{ end }}
}
```
Output:
```go
func main(){
    e := echo.New()
    
    e.GET("/", func(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
    })
    e.Logger.Fatal(e.Start(":8080"))
}
```

## Iterate an array
Values:
```yaml
endpoints: 
  - method: "GET"
    path: "/"
    status: "http.StatusOK"
    body: "Hello, World!"
  - method: "POST"
    path: "/payment"
    status: "http.StatusOK"
    body: "Payment is created"
```
File:
```go
func main(){
e := echo.New()

{{ range .endpoints }}
e.{{.method}}("{{.path}}", func (c echo.Context) error {
    return c.String({{.status}}, "{{.body}}")
})
{{ end }}
e.Logger.Fatal(e.Start(":8080"))
}
```
Output:
```go
func main(){
e := echo.New()


e.GET("/", func (c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
})

e.POST("/payment", func (c echo.Context) error {
    return c.String(http.StatusOK, "Payment is created")
})

e.Logger.Fatal(e.Start(":8080"))
}
```


## Iterate a map
Values:
```yaml
endpoints: 
   "/": "Hello, World!"
  "/payment":"Payment is created"
```
File:
```go
func main(){
e := echo.New()

{{ range $key, $value := .endpoints }}
e.GET("{{ $key }}", func (c echo.Context) error {
return c.String(http.StatusOK, "{{ $value }}")
})
{{ end }}

e.Logger.Fatal(e.Start(":8080"))
}
```
Output:
```go
func main(){
e := echo.New()


e.GET("/", func (c echo.Context) error {
return c.String(http.StatusOK, "Hello, World!")
})

e.GET("/payment", func (c echo.Context) error {
return c.String(http.StatusOK, "Payment is created")
})


e.Logger.Fatal(e.Start(":8080"))
}

```


