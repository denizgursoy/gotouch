# Go template

[Go template](https://pkg.go.dev/text/template) library is a powerful library. Some examples can be found below. 
[repeatit.io](https://repeatit.io/) is a great place to learn/practice templating. Following examples will redirect
you to the [repeatit.io](https://repeatit.io/).

### Use a value as text
If you have a value, for example `Port`, and want to write its value as text, you should write it between delimiters `{{` `}}` with a 
leading `.` as `{{.Port}}`. See the examples:

[Port Example](https://repeatit.io/#/share/eyJ0ZW1wbGF0ZSI6InBhY2thZ2UgbWFpblxuXG5pbXBvcnQgKFxuICAgXCJpb1wiXG4gICBcImxvZ1wiXG4gICBcIm5ldC9odHRwXCJcbilcblxuZnVuYyBtYWluKCkge1xuICAgaHR0cC5IYW5kbGVGdW5jKFwiL1wiLCBnZXRSb290KVxuICAgbG9nLkZhdGFsbG4oaHR0cC5MaXN0ZW5BbmRTZXJ2ZShcIjp7eyAuUG9ydCB9fVwiLCBuaWwpKVxufVxuXG5mdW5jIGdldFJvb3QodyBodHRwLlJlc3BvbnNlV3JpdGVyLCByICpodHRwLlJlcXVlc3QpIHtcbiAgIGlvLldyaXRlU3RyaW5nKHcsIFwiU2VydmVyIGdvdCB0aGUgcmVxdWVzdFxcblwiKVxufSIsImlucHV0IjoiUG9ydDogODA4MCIsImNvbmZpZyI6eyJ0ZW1wbGF0ZSI6InRleHQiLCJmdWxsU2NyZWVuSFRNTCI6ZmFsc2UsImZ1bmN0aW9ucyI6WyJzcHJpZyJdLCJvcHRpb25zIjpbImxpdmUiXSwiaW5wdXRUeXBlIjoieWFtbCJ9fQ==)

[Mail Example](https://repeatit.io/#/share/eyJ0ZW1wbGF0ZSI6IkhlbGxvIHt7LlJlY2VpdmVyfX1cblxuSSB3YW50ZWQgdG8gaW5mb3JtIHlvdSB0aGF0IEkgbGVhcm4ge3suVG9waWN9fVxuXG5CZXN0IHJlZ2FyZHMuLi5cbnt7LlNlbmRlcn19IiwiaW5wdXQiOiJSZWNlaXZlcjogTXIuIFNtaXRoXG5TZW5kZXI6IE1ycy4gU21pdGhcblRvcGljOiBHbyBUZW1wbGF0ZSBMaWJyYXJ5IiwiY29uZmlnIjp7InRlbXBsYXRlIjoidGV4dCIsImZ1bGxTY3JlZW5IVE1MIjpmYWxzZSwiZnVuY3Rpb25zIjpbInNwcmlnIl0sIm9wdGlvbnMiOlsibGl2ZSJdLCJpbnB1dFR5cGUiOiJ5YW1sIn19)

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

[HTTP server example](https://repeatit.io/#/share/eyJ0ZW1wbGF0ZSI6ImZ1bmMgbWFpbigpe1xue3sgaWYgZXEgIC5XZWJGcmFtZXdvcmsgXCJlY2hvXCIgfX1cbiAgICBlIDo9IGVjaG8uTmV3KClcbiAgICBcbiAgICBlLkdFVChcIi9cIiwgZnVuYyhjIGVjaG8uQ29udGV4dCkgZXJyb3Ige1xuICAgICAgICByZXR1cm4gYy5TdHJpbmcoaHR0cC5TdGF0dXNPSywgXCJIZWxsbywgV29ybGQhXCIpXG4gICAgfSlcbiAgICBlLkxvZ2dlci5GYXRhbChlLlN0YXJ0KFwiOjgwODBcIikpXG57e2Vsc2UgaWYgZXEgLldlYkZyYW1ld29yayBcImdpblwifX1cbiAgICByIDo9IGdpbi5EZWZhdWx0KClcbiAgICBcbiAgICByLkdFVChcIi9cIiwgZnVuYyhjICpnaW4uQ29udGV4dCkge1xuICAgICAgICBjLlN0cmluZyhodHRwLlN0YXR1c09LLCBcIkhlbGxvLCBXb3JsZCFcIilcbiAgICB9KVxuICAgIHIuUnVuKFwiOjgwODBcIilcbnt7IGVsc2UgfX1cbiAgICBodHRwLkhhbmRsZUZ1bmMoXCIvXCIsIGZ1bmMod3JpdGVyIGh0dHAuUmVzcG9uc2VXcml0ZXIsIHJlcXVlc3QgKmh0dHAuUmVxdWVzdCkge1xuICAgIFx0aW8uV3JpdGVTdHJpbmcod3JpdGVyLCBcInZlcnNpb24gMVwiKVxuICAgIH0pXG4gICAgXG4gICAgZXJyIDo9IGh0dHAuTGlzdGVuQW5kU2VydmUoXCI6ODA4MFwiLCBuaWwpXG4gICAgaWYgZXJyICE9IG5pbCB7XG4gICAgICAgIGxvZy5GYXRhbChlcnIpXG4gICAgfVxue3sgZW5kIH19XG59IiwiaW5wdXQiOiJXZWJGcmFtZXdvcms6IGRlZmF1bHQiLCJjb25maWciOnsidGVtcGxhdGUiOiJ0ZXh0IiwiZnVsbFNjcmVlbkhUTUwiOmZhbHNlLCJmdW5jdGlvbnMiOlsic3ByaWciXSwib3B0aW9ucyI6WyJsaXZlIl0sImlucHV0VHlwZSI6InlhbWwifX0=)

[To check boolean](https://repeatit.io/#/share/eyJ0ZW1wbGF0ZSI6ImZ1bmMgbWFpbigpe1xue3sgaWYgLklzRWNobyB9fVxuICAgIGUgOj0gZWNoby5OZXcoKVxuICAgIFxuICAgIGUuR0VUKFwiL1wiLCBmdW5jKGMgZWNoby5Db250ZXh0KSBlcnJvciB7XG4gICAgXHRyZXR1cm4gYy5TdHJpbmcoaHR0cC5TdGF0dXNPSywgXCJIZWxsbywgV29ybGQhXCIpXG4gICAgfSlcbiAgICBlLkxvZ2dlci5GYXRhbChlLlN0YXJ0KFwiOjgwODBcIikpXG57eyBlbmQgfX1cbn0iLCJpbnB1dCI6IklzRWNobzogdHJ1ZSIsImNvbmZpZyI6eyJ0ZW1wbGF0ZSI6InRleHQiLCJmdWxsU2NyZWVuSFRNTCI6ZmFsc2UsImZ1bmN0aW9ucyI6WyJzcHJpZyJdLCJvcHRpb25zIjpbImxpdmUiXSwiaW5wdXRUeXBlIjoieWFtbCJ9fQ==)

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


