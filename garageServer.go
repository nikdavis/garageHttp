package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    "net"
    "time"

    "github.com/zenazn/goji"
    "github.com/zenazn/goji/web"
    "github.com/hypebeast/gojistatic"
    "github.com/koron/go-ssdp"
    "github.com/stianeikeland/go-rpio"
)

const RELAY_PIN_GARAGE_DOOR = 4
const RELAY_PIN_LIGHT = 3
var DEVICES = map[string]int{
    "garageDoor":       RELAY_PIN_GARAGE_DOOR,
    "garageDoorLight":  RELAY_PIN_LIGHT,
}
const SERVER_PORT = 8000

type Response struct {
    Success bool
    Status int
}

func fireDevice(c web.C, w http.ResponseWriter, r *http.Request) {
    response := Response{Success: true, Status: 200}
    device := c.URLParams["device"]
    var b []byte
    var err error

    if (DEVICES[device] != 0) {
      cycleRelay(DEVICES[device])
    } else {
      response.Success = false
      response.Status = 404
    }

    if !response.Success {
        w.WriteHeader(response.Status)
    }

    b, err = json.Marshal(response)

    if err != nil {
        http.Error(w, "(╯°□°）╯︵ ┻━┻   woops, we lost our cool", http.StatusInternalServerError)
        return
    }

    fmt.Fprint(w, string(b))
}

func initializeRelayPin() {
    err := rpio.Open()
    if err != nil {
      panic(err)
    }

    relayPin := rpio.Pin(RELAY_PIN_GARAGE_DOOR)
    relayPin.High()
    relayPin.Output()
    relayPin = rpio.Pin(RELAY_PIN_LIGHT)
    relayPin.High()
    relayPin.Output()
}

func cycleRelay(pinId int) {
    pin := rpio.Pin(pinId)
    
    // Relay is active low -- Sunfounder, 2 channel, active low relays
    pin.Low()
    time.Sleep(100 * time.Millisecond)
    pin.High()
}

func main() {
    initializeRelayPin()

    // SSDP setup
    ad, err := ssdp.Advertise(
        "urn:nivvis-co:device:garageDoor:0-1",
        "uuid:f29c575e-8ec0-4cd9-a359-1a4491bc4f79",
        getDeviceDetailsURL("garageDoor"),
        "go-ssdp sample",
        1800)

    if err != nil {
        panic(err)
    }

    ad, err = ssdp.Advertise(
        "urn:nivvis-co:device:garageDoorLight:0-1",
        "uuid:2dc7a051-0f66-45f0-9f7c-51bda3b9d894",
        getDeviceDetailsURL("garageDoorLight"),
        "go-ssdp sample",
        1800)

    if err != nil {
        panic(err)
    }

    // Web server
    goji.Get("/fire/:device", fireDevice)
    goji.Use(gojistatic.Static("public", gojistatic.StaticOptions{}))
    goji.Serve() // block

    // SSDP shutdown
    // send/multicast "byebye" message.
    ad.Bye()
    // teminate Advertiser.
    ad.Close()
}

func getDeviceDetailsURL(device string) string {
  localIPAddr := getLocalIP()
  return fmt.Sprintf("http://%s:%d/%sDetails.xml", localIPAddr, device, SERVER_PORT)
}

func getLocalIP() string {
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
