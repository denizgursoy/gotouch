# Go template

[Go template](https://pkg.go.dev/text/template) library is a powerful library. Some examples can be found below. 
[repeatit.io](https://repeatit.io/) is a great place to learn/practice templating. Following examples will redirect
you to the [repeatit.io](https://repeatit.io/).

### Use a value as text
If you have a value, for example `Port`, and want to write its value as text, you should write it between delimiters `{{` `}}` with a 
leading `.` as `{{.Port}}`. See the examples:

[Port Example](shorturl.at/HKL37)

[Mail Example](shorturl.at/mwyJM)

### If, else if, else 
Go template library allows you to do conditional rendering with `if` statements . You can use following keywords in `if` statements to compare values

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

See examples:

[HTTP server example](shorturl.at/FHQW2)

[To check boolean](shorturl.at/hvxyz)

### Iterate an array
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
### Iterate a map
Values:
```yaml
endpoints: 
  "/": "Hello, World!"
  "/payment": "Payment is created"
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


