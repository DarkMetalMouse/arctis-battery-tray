package main

import (
	"fmt"
	"time"

	hid "github.com/DarkMetalMouse/hid"
	"github.com/getlantern/systray"
)

var dev hid.Device = nil

func getBatteryLevel() (int, error) {
	if dev == nil {
		return -1, fmt.Errorf("HID device uninitialized")
	}
	var err error
	if err = dev.Write([]byte{0x06, 0x18}); err != nil {
		return -1, err
	}

	var response []byte = []byte{0, 0, 0}
	for response[1] != 0x18 { // Sometimes we can get the response for the wrong request
		if response, err = dev.Read(); err != nil {
			return -1, err

		}
	}
	level := int(response[2])
	if level > 100 {
		fmt.Println(response)
		level = 100
	}

	return level, err
}

func main() {

	var err error

	c := hid.FastFindDevices(0x1038, 0x12AD)

	for {
		devInfo, ok := <-c
		if !ok {
			break
		}
		if devInfo.OutputReportLength == 31 {
			if dev, err = devInfo.Open(); err != nil {
				fmt.Printf("Open failed: %s\n", err)
				return
			}
			defer dev.Close()
			break
		}
	}

	onExit := func() {
		fmt.Println("Tray exited")
	}

	systray.Run(onReady, onExit)

}

func onReady() {
	systray.SetTemplateIcon(icons[0], icons[0])
	systray.SetTitle("Arctis 7 Battery Level")
	systray.SetTooltip("Arctis 7 Battery Level")
	quit := systray.AddMenuItem("Quit", "Quit the app")

	go func() {
		for {
			select {

			// quit
			case <-quit.ClickedCh:
				systray.Quit()
				return
			case <-time.After(time.Second):
				if level, err := getBatteryLevel(); err == nil {
					systray.SetTemplateIcon(icons[level], icons[level])
				} else {
					systray.SetTemplateIcon(icons[0], icons[0])

				}

			}
		}
	}()

}
