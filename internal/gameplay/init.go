package gameplay

import (
	"buttplugosu/pkg/logging"
	"buttplugosu/pkg/memory"
	"regexp"
)

var osuProcessRegex = regexp.MustCompile(`.*osu!\.exe.*`)
var patterns staticAddresses
var menuData menuD
var gameplayData gameplayD

func initBase() error {
	var err error

	// find osu process
	processes, err = memory.FindProcess(osuProcessRegex, "osu!lazer", "osu!framework")
	if err != nil {
		return err
	}

	process = processes[0]

	logging.Global.Info().
		Int("pid", process.Pid()).
		Msg("Found process")

	// resolve song select pattern
	err = memory.ResolvePatterns(process, &patterns.PreSongSelectAddresses)
	if err != nil {
		logging.Global.
			Err(err).
			Msg("Resolving patterns failed")
		return err
	}

	// read pre song select data
	if err = memory.Read(process, &patterns.PreSongSelectAddresses, &menuData.PreSongSelectData); err != nil {
		logging.Global.
			Err(err).
			Msg("Reading failed")
		return err
	}

	// resolve all other patterns
	logging.Global.Info().
		Msg("Resolving patterns")

	err = memory.ResolvePatterns(process, &patterns)
	if err != nil {
		return err
	}

	logging.Global.Info().
		Msg("Resolved patterns")

	DynamicAddresses.IsReady = true
	return nil
}
