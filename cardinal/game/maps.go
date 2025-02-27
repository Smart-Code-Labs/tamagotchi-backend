package game

// Define a struct to hold food properties including description
type FoodProperties struct {
	Price       float64
	Health      int
	Energy      int
	Description string
}

// Initialize the FoodKinds map with FoodProperties including descriptions
var FoodKinds = map[string]FoodProperties{
	"Apple":   {Price: 0.1, Health: 10, Energy: 10, Description: "Yuumy Red Food"},
	"Banana":  {Price: 0.3, Health: 5, Energy: 15, Description: "What is this Yellow Food?"},
	"Soup":    {Price: 0.5, Health: 15, Energy: 20, Description: "Spicy!!!"},
	"Carrots": {Price: 0.1, Health: 5, Energy: 25, Description: "Cheap, but powerful"},
}

// Define a struct to hold drug properties including description
type DrugProperties struct {
	Price       float64
	Value       int
	Description string
}

// Initialize the DrugKinds map with DrugProperties including descriptions
var DrugKinds = map[string]DrugProperties{
	"Vaccine": {Price: 5.0, Value: 80, Description: "A vaccine to boost your health!"},
	"Pill":    {Price: 1.0, Value: 20, Description: "A small pill to help you recover."},
	"Vitamin": {Price: 0.5, Value: 15, Description: "Essential vitamins for daily health."},
	"Mineral": {Price: 0.1, Value: 10, Description: "Important minerals to keep you strong."},
}

var BathKinds = map[string]DrugProperties{
	"Sponge": {Price: 0.1, Value: 30, Description: "Basic clean up item."},
}

// ToyKinds map initializes each toy with its properties
var ToyKinds = map[string]ToyProperties{
	"Ball":    {Name: "Ball", Description: "Yuuju!", Price: 5.0, Wellness: 15},
	"Frisbee": {Name: "Frisbee", Description: "Will be back?", Price: 1.0, Wellness: 10},
	"Rope":    {Name: "Rope", Description: "Grrrr", Price: 0.5, Wellness: 10},
	"Stick":   {Name: "Stick", Description: "Throw it! Throw it!", Price: 0.1, Wellness: 5},
}

// ToyProperties holds all necessary properties for a toy
type ToyProperties struct {
	Name        string
	Description string
	Price       float64
	Wellness    int
}

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

// Breed
var Skills = []string{"Intellect", "Force", "skilled"}
var Elements = []string{"wynd", "water", "fire", "earth"}
