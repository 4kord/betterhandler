# BetterHandler
BetterHandler is a handler with some useful features that implements the standard net/http handler interface.

## Feautures

- Writing String into ResponseWriter
- Writing JSON into ResponseWriter
- Writing XML into ResponseWriter
---
- Parsing request body into a struct depending on the Content-Type header. Supported: "application/json", "application/xml", "multipart/form-data".
---
- Getting base url
---
- Set cookie
- Get cookie
- Get cookie value
- Expire a client cookie / all cookies

## Installation

```sh
go get -u github.com/4kord/betterhandler
```

## Quickstart
```go
package main

import (
    "net/http"
    "github.com/4kord/betterhandler"
)

func main() {
    http.ListenAndServe(":3000", betterhandler.BH(func(c *betterhandler.Ctx) {
        c.rw.WriteHeader(http.StatusOK)

        c.String("Hello, World!")
    }))
}
```

## Examples
##### Writing String into ResponseWriter
```golang
func handler(c *betterhandler.Ctx) {
    c.String("Hello, World!")
}
```

##### Writing JSON into ResponseWriter
```golang
type Json struct {
    Key1 string  `json:"key1"`
    Key2 int     `json:"key2"`
    Key3 float64 `json:"key3"`
}

func handler(c *betterhandler.Ctx) {
    c.JSON(Json{
        Key1: "value1",
        Key2: 10,
        Key3: 3.14,
    })
}

func handler2(c *betterhandler.Ctx) {
    c.JSON(betterhandler.Map{
        "Key1": "value1",
        "Key2": 10,
        "Key3": 3.14,
    })
}
```

##### Writing XML into ResponseWriter
```golang
type Xml struct {
    Key1 string  `xml:"key1"`
    Key2 int     `xml:"key2"`
    Key3 float64 `xml:"key3"`
}

func handler(c *betterhandler.Ctx) {
    c.XML(Xml{
        Key1: "value1",
        Key2: 10,
        Key3: 3.14,
    })
}

func handler2(c *betterhandler.Ctx) {
    c.XML(betterhandler.Map{
        "Key1": "value1",
        "Key2": 10,
        "Key3": 3.14,
    })
}
```

##### Parsing request body into a struct
###### STRUCT TAGS: application/json - json, application/xml - xml, multipart/form-data - form
```go
type V struct {
    Key1 string                  `form:"key1"`
    Key2 int                     `form:"key2"`
    Key3 float64                 `form:"key3"`
    Key4 []*Multipart.FileHeader `form:"key4"`
}

func handler(c *betterhandler.Ctx) {
    var myStruct V
    
    c.BodyParser(&myStruct)
}
```

##### Getting base url
```go
func handler(c *betterhandler.Ctx) {
    fmt.Println(c.BaseUrl())
}
```

##### Set cookie
```go
func handler(c *betterhandler.Ctx) {
    cookie := &http.Cookie{
        Name: "myCookie",
        Value: "value",
        HttpOnly: true,
    }

    c.SetCookie(cookie)
}
```

##### Get cookie
```go
func handler(c *betterhandler.Ctx) {
    cookie, err := c.GetCookie("myCookie")
    if err != nil {
        if errors.Is(err, http.ErrNoCookie) {
            fmt.Println("Cookie not found")
            return
        }
        fmt.Println("Unexpected error")
    }
}
```

##### Expire a client cookie / all cookies
```go
func handler(c *betterhandler.Ctx) {
    c.ClearCookie("myCookie") // Expire myCookie
}

func handler2(c *betterhandler.Ctx) {
    c.ClearCookie("myCookie2") // Expire myCookie, myCookie2
}

func handler3(c *betterhandler.Ctx) {
    c.ClearCookie() // Expire all cookies
}
```
