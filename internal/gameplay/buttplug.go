package gameplay

import (
	"buttplugosu/pkg/logging"
	"github.com/pidurentry/buttplug-go"
	"github.com/pidurentry/buttplug-go/device"
	"time"
)

var (
	Client        buttplug.DeviceManager
	vibratorQueue = make(chan time.Duration)
)

// handleVibrationQueue handles the vibration queue
func handleVibrationQueue() {
	var speed = device.Speed{Speed: 1.0} // Why is this a structure?

	for v := range vibratorQueue {
		if Client == nil || len(Client.Vibrators()) <= 0 {
			continue
		}

		// go through all vibrators and start them
		for _, x := range Client.Vibrators() {
			_ = x.Vibrate(speed)
		}

		time.Sleep(v)

		// go through all vibrators and stop them
		for _, d := range Client.Vibrators() {
			_ = d.Stop()
		}
	}
}

func HandlePlug() {
	go handleVibrationQueue()

	client, err := buttplug.Dial("ws://127.0.0.1:12345/buttplug")
	if err != nil {
		logging.Global.Fatal().
			Err(err).
			Msg("Connection failed")
		return
	}

	// create new handler and handshake
	handler := buttplug.NewHandler(client)
	info, err := handler.Handshake("osu")
	if err != nil {
		logging.Global.Fatal().
			Err(err).
			Msg("Handshake failed")
		return
	}

	logging.Global.Info().
		Msg("Connected to WebSocket")

	// handle pings
	go func() {
		ticker := info.MaxPingTime().Ticker()

		for {
			<-ticker.C

			if !handler.Ping() {
				continue
			}
		}
	}()

	logging.Global.Info().
		Msg("Scanning for devices")

	// create new device manager and scans for...toys
	Client = buttplug.NewDeviceManager(handler)
	Client.Scan(5 * time.Second).Wait()

	if len(Client.Devices()) <= 0 {
		logging.Global.Fatal().
			Msg("No devices found")
	}

	logging.Global.Info().
		Int("count", len(Client.Devices())).
		Msg("Found devices")
}
