package gameplay

import (
	"buttplugosu/pkg/logging"
	"buttplugosu/pkg/mem"
	"strings"
	"time"
)

func Init() {
	if err := initBase(); err != nil {
		logging.Global.Fatal().
			Err(err).
			Msg("Error occurred while initializing")
	}

	for {
		start := time.Now()

		if DynamicAddresses.IsReady {
			if err := mem.Read(
				process,
				&patterns.PreSongSelectAddresses,
				&menuData.PreSongSelectData,
			); err != nil {
				logging.Global.
					Err(err).
					Msg("Failed to read 'PreSongSelectData'")

				DynamicAddresses.IsReady = false
				continue
			}

			handleRead()
		}

		elapsed := time.Since(start)
		time.Sleep(time.Duration(1-int(elapsed.Milliseconds())) * time.Millisecond)
	}
}

func handleRead() {
	err := mem.Read(process, &patterns, &gameplayData)
	if err != nil &&
		!strings.Contains(err.Error(), "LeaderBoard") &&
		!strings.Contains(err.Error(), "KeyOverlay") {
		return
	}

	// shit way tbh but it works
	if previousHits != int(gameplayData.HitMiss) {
		vibratorQueue <- 300 * time.Millisecond

		logging.Global.Debug().
			Int16("count", gameplayData.HitMiss).
			Msg("Queueing vibration due to miss")

		previousHits = int(gameplayData.HitMiss)
	}
}
