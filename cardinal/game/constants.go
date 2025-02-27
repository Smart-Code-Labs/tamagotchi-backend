package game

// Constants

// Tick system
const TickRate = 1                   // Ticks per second
const TickFiveSeconds = 5 / TickRate // 5 seconds
const TickMinute = 60 / TickRate     // 1 minute
const TickHour = 3600 / TickRate     // 1 hour

const DeclineTickRate = TickFiveSeconds // Decline every minute
const ThinkTickRate = TickFiveSeconds   // Think every minute
const ActivityUpdateTickRate = TickRate

// Times
const TickEightHours = TickHour * 8 // Sleeping

// Create Pet
const (
	MaxHP           = 100
	MaxEnergy       = 100
	MaxHygiene      = 100
	MaxWellness     = 100
	InitialActivity = "None"
	InitialThink    = "..."
	MaxLevel        = int64(10)
)

const PetCost = 5

// Decline system
const HygieneThreshold = 70

// Pet Play method
const ExperienceEarn = 20
const EnergyReduce = 10
const HygieneReduce = 5
const WellnessIncrease = 5

// Pet Sleep method
const EnergyIncrease = 80

// Pet Eat method
const HealthIncrease = 20

// Pet Bath method
const HygieneIncrease = 20

// Pet Think
const ThinkSleep = "Zzz...Zzz"
const ThinkBath = "(Singing...)"
const ThinkEat = "Mmm Yummy!"
const ThinkPlay = "Love to play!"

// Pet Activity
const PetEarnMoney = 0.0001

// Player
const PlayerInitialMoney = float64(1000)
