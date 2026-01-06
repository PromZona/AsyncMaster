package bot

type UserState int

const (
	// Normal State of Being
	UserStateDefault UserState = 0

	// Registration Phase
	UserStateAwaitPassword = 100
	UserStateAwaitCodename = 101

	// Master Commands
	UserStateAwaitSavingMessage    = 200
	UserStateAwaitTitleForMesssage = 201

	// Player Commands
	UserStateAwaitResipient     = 300
	UserStateAwaitMessage       = 301
	UserStateAwaitTitleDecision = 302
	UserStateAwaitTitle         = 303
)
