/**
 *  Garage Door Light
 *
 *  Copyright 2018 Nikolas Davis
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License. You may obtain a copy of the License at:
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software distributed under the License is distributed
 *  on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License
 *  for the specific language governing permissions and limitations under the License.
 *
 */
metadata {
	definition (name: "Garage Door Light", namespace: "nikdavis", author: "Nikolas Davis") {
    	capability "Actuator"
        capability "Switch"
		capability "Light"
	}


	simulator {
		// TODO: define status and reply messages here
	}

	tiles {
        standardTile("main", "device.switch", width: 3, height: 2, canChangeIcon: true, decoration: "flat") {
            state "off", label: '${name}', action: "switch.on", icon: "st.switches.light.off", backgroundColor: "#ffffff"
            state "on", label: '${name}', action: "switch.off", icon: "st.switches.light.on", backgroundColor: "#8a29f2"
        }

		main "main"
	}
}

// parse events into attributes
def parse(String description) {
	log.debug "Parsing '${description}'"

}

// handle commands
def off() {
	log.debug "Executing 'off'"
	toggleLight("off")
}

def on() {
	log.debug "Executing 'on'"
	toggleLight("on")
}

def toggleLight(value) {
	sendEvent(name: "switch", value: value)

    def result = new physicalgraph.device.HubAction(
        method: "GET",
        path: "/fire/garageDoorLight",
        headers: [
            HOST: "${getHostAddress()}:8000"
        ]
    )
}

// gets the address of the Hub
private getCallBackAddress() {
    return device.hub.getDataValue("localIP") + ":" + device.hub.getDataValue("localSrvPortTCP")
}

// gets the address of the device
private getHostAddress() {
    def ip = getDataValue("ip")
    def port = getDataValue("port")

    if (!ip || !port) {
        def parts = device.deviceNetworkId.split(":")
        if (parts.length == 2) {
            ip = parts[0]
            port = parts[1]
        } else {
            log.warn "Can't figure out ip and port for device: ${device.id}"
        }
    }

    log.debug "Using IP: $ip and port: $port for device: ${device.id}"
    return convertHexToIP(ip) + ":" + convertHexToInt(port)
}

private Integer convertHexToInt(hex) {
    return Integer.parseInt(hex,16)
}

private String convertHexToIP(hex) {
    return [convertHexToInt(hex[0..1]),convertHexToInt(hex[2..3]),convertHexToInt(hex[4..5]),convertHexToInt(hex[6..7])].join(".")
}
