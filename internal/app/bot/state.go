package bot

type UserState int

const (
	Idle                      UserState = 0
	InRegistrationFlow        UserState = 1
	InSendMessageFlow         UserState = 2
	InCreateMasterRequestFlow UserState = 3
)

/* old one
DELETE AFTER REFACTOR
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

)
*/
