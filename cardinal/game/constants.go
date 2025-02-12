package game

// Constants

// Health decline system
const HealthDeclineTicksPerSecond = 3600
const HygieneThresshold = 70

// Pet Play method
const EnergyReduce = 10
const HygieneReduce = 5
const WellnessIncrease = 5
const PlayDuration = 3600 // 1 * 60 x 60

// Pet Sleep method
const EnergyIncrease = 80
const SleepDuration = 28800 // 8 * 60 x 60

// Pet Eat method
const HealthIncrease = 20
const EatDuration = 3600 // 8 * 60 x 60

// Pet bath method
const HygieneIncrease = 20
const BathDuration = 3600 // 8 * 60 x 60

type Constant struct {
	Label string
	Value any
}

type WorldConstant struct {
	InstanceName  string
	InstanceTimer int
	TickRate      int
}

type BossLevelStats struct {
	Level         int64
	EnergyDefault string // decimal
	EnergyMax     string // decimal
	EnergyRefill  string // decimal
	Range         string // decimal
	Speed         string // decimal
	Defense       string // decimal
	Difficulty    string // decimal
}

var (
	AllConstantsLabel = "all"
	// ExposedConstants If you want the constant to be queryable through `query_constant`,
	// make sure to add the constant to the list of exposed constants
	// Stored as a pointer because we update these constants via
	// tx/game/set-constant
	ExposedConstants = []Constant{
		{
			Label: "world",
			Value: &WorldConstants,
		},
	}

	WorldConstants = WorldConstant{
		InstanceName:  "", // Set in SetConstantsFromEnv()
		InstanceTimer: 0,  // Set in SetConstantsFromEnv()
		TickRate:      2,  // Ticks per second
	}

	Bosses = [1]*BossLevelStats{
		&BossLevel0Stats,
	}

	BossLevel0Stats = BossLevelStats{
		Level:         0,
		EnergyDefault: "0",
		EnergyMax:     "100",
		EnergyRefill:  "30",
		Range:         "40",
		Speed:         "1",
		Defense:       "500",
		Difficulty:    "10",
	}
)
