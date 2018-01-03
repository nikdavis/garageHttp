### To compile and run on for Raspberry Pi Zero W
`GOOS="linux" GOARCH="arm" GOARM=6 golang build`

`scp gojiExample pi@192.168.1.22:~/gojiExample`

ssh...

`./gojiExample`

Note: I should perhaps attempt to figure out the best way to manage this service on the Pi long term, and maybe make the output a container.

### Todo
* Establish a connection over SPI, reading in the tilt sensor
* Create a GET endpoint that will share the state of the sensor
* Create PUT endpoint that will allow updating of the garage doors state (should this be blocking?)

### Resources
[SmartThings LAN connected docs](http://docs.smartthings.com/en/latest/cloud-and-lan-connected-device-types-developers-guide/building-lan-connected-device-types/index.html)
[Pi pinout](https://pinout.xyz/pinout/spi)
[Pi power reqs](https://www.raspberrypi.org/help/faqs/#powerReqs)
[Pi GPIO library](https://github.com/stianeikeland/go-rpio)
[Relay overview](http://wiki.sunfounder.cc/index.php?title=2_Channel_5V_Relay_Module)
[Garage door wiring](https://github.com/nikdavis/chamberlain-liftmaster-opener)
