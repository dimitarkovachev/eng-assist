package transcription

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

func GetDeviceIDByName(name string) int {
	// Initialize SDL2 Audio subsystem
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		log.Fatalf("Failed to initialize SDL: %v", err)
	}
	defer sdl.Quit()

	// Get number of capture devices (true = capture devices, false = playback devices)
	numDevices := sdl.GetNumAudioDevices(true)

	if numDevices < 0 {
		log.Fatalf("Failed to get audio devices: %v", sdl.GetError())
	}

	fmt.Printf("init: found %d capture devices:\n", numDevices)

	deviceID := 999
	// List all capture devices
	for i := 0; i < numDevices; i++ {
		deviceName := sdl.GetAudioDeviceName(i, true) // true = capture
		if deviceName != "" {
			fmt.Printf("init:    - Capture device #%d: '%s'\n", i, deviceName)
		} else {
			fmt.Printf("init:    - Capture device #%d: (unknown)\n", i)
		}
		if deviceName == name {
			deviceID = i
		}
	}

	fmt.Printf("init: found deviceID: %d\n", deviceID)
	return deviceID
}
