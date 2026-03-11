package runtime

type Context interface {
	ChatID() int64
	FirstName() string
	Text() string
	Callback() string
	Send(string, ...Keyboard) error
	SendTo(id int64, text string, k ...Keyboard) error
	Respond()
	Args() []string
	MessageID() int64
	MessageText() string
}

type Runtime interface {
	Use(Middleware)
	HandleText(Handler)
	HandleCallback(Handler)
	HandleCommand(string, Handler)
	Start() error
}

type Handler func(Context) error
type Middleware func(Handler) Handler

// btnSend := menu.Data("Send Message", sendmsgc.CBSend)
type Button struct {
	Unique string
	Text   string
	Data   string
}

type Row []Button
type Keyboard []Row
