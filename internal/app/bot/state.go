package bot

type UserState int

const (
	Idle                      UserState = 0
	InRegistrationFlow        UserState = 1
	InSendMessageFlow         UserState = 2
	InCreateMasterRequestFlow UserState = 3
)
