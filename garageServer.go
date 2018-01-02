package main

import (
    "fmt"
    "encoding/json"
    "net/http"

    "github.com/zenazn/goji"
    "github.com/zenazn/goji/web"
    "github.com/hypebeast/gojistatic"
)

type HelloWorld struct {
    Hello string
}

type Error struct {
    Status int
    Error string
}

type Payload interface {

}

type Response struct {
    Payload Payload
    Error   *Error
}

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
    response := Response{Payload: &HelloWorld{Hello: c.URLParams["name"]}}
    var b []byte
    var err error

    if response.Error != nil {
        w.WriteHeader(response.Error.Status)
        b, _ = json.Marshal(*response.Error)
    } else {
        b, _ = json.Marshal(response.Payload)
    }

    if err != nil {
        http.Error(w, "(╯°□°）╯︵ ┻━┻   woops, we lost our cool", http.StatusInternalServerError)
        return
    }

    fmt.Fprint(w, string(b))
}

func main() {
    goji.Get("/hello/:name", hello)
    goji.Use(gojistatic.Static("public", gojistatic.StaticOptions{}))
    goji.Serve()
}
