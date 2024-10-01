package gameplay

type PreSongSelectAddresses struct {
	Status        int64 `sig:"48 83 F8 04 73 1E"`
	SettingsClass int64 `sig:"83 E0 20 85 C0 7E 2F"`
}

type PreSongSelectData struct {
	Status uint32 `memory:"[Status - 0x4]"`
}

type staticAddresses struct {
	PreSongSelectAddresses
	Base        int64 `sig:"F8 01 74 04 83 65"`
	MenuMods    int64 `sig:"C8 FF ?? ?? ?? ?? ?? 81 0D ?? ?? ?? ?? 00 08 00 00"`
	PlayTime    int64 `sig:"5E 5F 5D C3 A1 ?? ?? ?? ?? 89 ?? 04"`
	ChatChecker int64 `sig:"0A D7 23 3C 00 00 ?? 01"`
	SkinData    int64 `sig:"75 21 8B 1D"`
	Rulesets    int64 `sig:"7D 15 A1 ?? ?? ?? ?? 85 C0"`
	ChatArea    int64 `sig:"33 47 9D FF 5B 7F FF FF"`
}

func (staticAddresses) Ruleset() string {
	return "[[Rulesets - 0xB] + 0x4]"
}

func (staticAddresses) Beatmap() string {
	return "[Base - 0xC]"
}

func (PreSongSelectAddresses) Settings() string {
	return "[SettingsClass + 0x8]"
}

func (staticAddresses) PlayContainer() string {
	return "[[[[PlayContainerBase + 0x7] + 0x4] + 0xC4] + 0x4]"
}

func (staticAddresses) Leaderboard() string {
	return "[[[LeaderboardBase+0x1] + 0x4] + 0x7C] + 0x24"
}

type menuD struct {
	PreSongSelectData
	MenuGameMode       int32   `memory:"[Base - 0x33]"`
	Plays              int32   `memory:"[Base - 0x33] + 0xC"`
	Artist             string  `memory:"[[Beatmap] + 0x18]"`
	ArtistOriginal     string  `memory:"[[Beatmap] + 0x1C]"`
	Title              string  `memory:"[[Beatmap] + 0x24]"`
	TitleOriginal      string  `memory:"[[Beatmap] + 0x28]"`
	AR                 float32 `memory:"[Beatmap] + 0x2C"`
	CS                 float32 `memory:"[Beatmap] + 0x30"`
	HP                 float32 `memory:"[Beatmap] + 0x34"`
	OD                 float32 `memory:"[Beatmap] + 0x38"`
	StarRatingStruct   uint32  `memory:"[Beatmap] + 0x8C"`
	AudioFilename      string  `memory:"[[Beatmap] + 0x64]"`
	BackgroundFilename string  `memory:"[[Beatmap] + 0x68]"`
	Folder             string  `memory:"[[Beatmap] + 0x78]"`
	Creator            string  `memory:"[[Beatmap] + 0x7C]"`
	Name               string  `memory:"[[Beatmap] + 0x80]"`
	Path               string  `memory:"[[Beatmap] + 0x90]"`
	Difficulty         string  `memory:"[[Beatmap] + 0xAC]"`
	MapID              int32   `memory:"[Beatmap] + 0xC8"`
	SetID              int32   `memory:"[Beatmap] + 0xCC"`
	RankedStatus       int32   `memory:"[Beatmap] + 0x12C"` // unknown, unsubmitted, pending/wip/graveyard, unused, ranked, approved, qualified
	MD5                string  `memory:"[[Beatmap] + 0x6C]"`
	ObjectCount        int32   `memory:"[Beatmap] + 0xFC"`
}

type gameplayD struct {
	Retries             int32   `memory:"[Base - 0x33] + 0x8"`
	PlayerName          string  `memory:"[[[Ruleset + 0x68] + 0x38] + 0x28]"`
	ModsXor1            int32   `memory:"[[[Ruleset + 0x68] + 0x38] + 0x1C] + 0xC"`
	ModsXor2            int32   `memory:"[[[Ruleset + 0x68] + 0x38] + 0x1C] + 0x8"`
	HitErrors           []int32 `memory:"[[[Ruleset + 0x68] + 0x38] + 0x38]"`
	Mode                int32   `memory:"[[Ruleset + 0x68] + 0x38] + 0x64"`
	MaxCombo            int16   `memory:"[[Ruleset + 0x68] + 0x38] + 0x68"`
	ScoreV2             int32   `memory:"Ruleset + 0x100"`
	Hit100              int16   `memory:"[[Ruleset + 0x68] + 0x38] + 0x88"`
	Hit300              int16   `memory:"[[Ruleset + 0x68] + 0x38] + 0x8A"`
	Hit50               int16   `memory:"[[Ruleset + 0x68] + 0x38] + 0x8C"`
	HitGeki             int16   `memory:"[[Ruleset + 0x68] + 0x38] + 0x8E"`
	HitKatu             int16   `memory:"[[Ruleset + 0x68] + 0x38] + 0x90"`
	HitMiss             int16   `memory:"[[Ruleset + 0x68] + 0x38] + 0x92"`
	Combo               int16   `memory:"[[Ruleset + 0x68] + 0x38] + 0x94"`
	PlayerHPSmooth      float64 `memory:"[[Ruleset + 0x68] + 0x40] + 0x14"`
	PlayerHP            float64 `memory:"[[Ruleset + 0x68] + 0x40] + 0x1C"`
	Accuracy            float64 `memory:"[[Ruleset + 0x68] + 0x48] + 0xC"`
	LeaderBoard         uint32  `memory:"[Ruleset + 0x7C] + 0x24"`
	KeyOverlayArrayAddr uint32  `memory:"[[Ruleset + 0xB0] + 0x10] + 0x4"` //has to be at the end due to memory not liking dead pointers, TODO: Fix this memory-side
	// Score               int32   `memory:"[[Ruleset + 0x68] + 0x38] + 0x78"`
}
