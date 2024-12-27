package interfaces

// Bot defines the methods needed by commands
type Bot interface {
	GetUserAliases() map[string][]string
}
