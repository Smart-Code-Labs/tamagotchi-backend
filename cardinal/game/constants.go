package game

// Constants

// Tick system
const TickRate = 1               // Ticks per second
const TickMinute = 60 / TickRate // 1 minute
const TickHour = 3600 / TickRate // 1 hour

const DeclineTickRate = TickMinute // Decline every minute
const ThinkTickRate = TickMinute   // Think every minute
const ActivityUpdateTickRate = TickRate

// Times
const TickEightHours = TickHour * 8 // Sleeping

// Create Pet
const (
	InitialHP = 100
	InitialE  = 100
	InitialHy = 100
	InitialWn = 100
)

// Decline system
const HygieneThresshold = 70

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

// pet Thinking
// Define a custom type to hold the min and max values.  This makes it clearer
// what the constant represents and allows you to easily add more related
// values later if needed (e.g., a name, a unit, etc.).
type Range struct {
	Min int
	Max int
}

type Message struct {
	Text string
}

var HealthMessages = map[Range]Message{
	{Min: 30, Max: 40}: {Text: "I don't feel well."},
	{Min: 0, Max: 30}:  {Text: "Im gona dye!!!"},
}

var HygieneMessages = map[Range]Message{
	{Min: 50, Max: 60}: {Text: "My whole body itches"},
	{Min: 0, Max: 40}:  {Text: "OMG!!! Im really dirty! Someone please Bath me. Please!"},
}

var WellnessMessages = map[Range]Message{
	{Min: 30, Max: 70}: {Text: "I fell a little depress today."},
	{Min: 0, Max: 30}:  {Text: "Dont know what to do..."},
}

var EnergyMessages = map[Range]Message{
	{Min: 70, Max: 80}: {Text: "Im bored. I would kill to go outside."},
	{Min: 0, Max: 60}:  {Text: "Im bored... to death? Play with me!"},
}
