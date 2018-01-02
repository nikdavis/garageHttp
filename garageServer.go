package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    "net"

    "github.com/zenazn/goji"
    "github.com/zenazn/goji/web"
    "github.com/hypebeast/gojistatic"
    "github.com/koron/go-ssdp"
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
    // SSDP Server
    externalIP := GetLocalIP()
    port := "8000"
    fmt.Println(externalIP)

    location := fmt.Sprintf("http://%s:%s/details.xml", externalIP, port)
    ad, err := ssdp.Advertise(
        "urn:nivvis-co:device:garageDoor:0-1",                        // send as "ST"
        "uuid:f29c575e-8ec0-4cd9-a359-1a4491bc4f79",                        // send as "USN"
        location, // send as "LOCATION"
        "go-ssdp sample",                   // send as "SERVER"
        1800)                               // send as "maxAge" in "CACHE-CONTROL"

    if err != nil {
        panic(err)
    }

    // Web server
    goji.Get("/hello/:name", hello)
    goji.Use(gojistatic.Static("public", gojistatic.StaticOptions{}))
    goji.Serve() // block

    // send/multicast "byebye" message.
    ad.Bye()
    // teminate Advertiser.
    ad.Close()
}

func GetLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}
