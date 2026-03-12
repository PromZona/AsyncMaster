package runtime

import (
	"bufio"
	"fmt"
	"os"
	// "strings"
)

type MockContext struct {
	chatID      int64
	firstName   string
	callback    string
	args        []string
	messageID   int64
	messageText string
}

func (c *MockContext) ChatID() int64 {
	return c.chatID
}
func (c *MockContext) FirstName() string {
	return c.firstName
}
func (c *MockContext) Callback() string {
	return c.callback
}

func (c *MockContext) Send(string, ...Keyboard) error {
	return nil
}
func (c *MockContext) SendTo(id int64, text string, k ...Keyboard) error {
	return nil
}

func (c *MockContext) Respond() {
	// Respond is a telegram thing to send a small nice message while doing callback
	// We do not use them for now, so will leave it empty
}

func (c *MockContext) Args() []string {
	return c.args
}
func (c *MockContext) MessageID() int64 {
	return c.messageID
}
func (c *MockContext) MessageText() string {
	return c.messageText
}

type MockRuntime struct {
	middlewares []Middleware
	mock        Mock
	userManager MockUserManager
}

func NewMockRuntime() *MockRuntime {
	return &MockRuntime{
		middlewares: make([]Middleware, 0),
		mock:        Mock{},
		userManager: MockUserManager{},
	}
}

func (mr *MockRuntime) Use(m Middleware) {
	mr.middlewares = append(mr.middlewares, m)
}
func (mr *MockRuntime) HandleText(h Handler) {
	mr.apply(h)

}
func (mr *MockRuntime) HandleCallback(Handler) {

}
func (mr *MockRuntime) HandleCommand(string, Handler) {

}

func (mr *MockRuntime) apply(h Handler) Handler {
	for i := len(mr.middlewares) - 1; i >= 0; i-- {
		h = mr.middlewares[i](h)
	}
	return h
}

func (mr *MockRuntime) Start() error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		command, err := parseCommandLine(input)
		if err != nil {
			return err
		}

		isExit := false
		switch command.Command {
		case CTServer:
			err := processServerCommand(command, &mr.userManager)
			if err != nil {
				return err
			}
		case CTUser:
			err := processUserCommand(command)
			if err != nil {
				return err
			}

		case CTExit:
			isExit = true
		}
		if isExit {
			break
		}
	}

	if scanner.Err() != nil {
		fmt.Printf("Error met")
	}

	return nil
}

type CommandType int

const (
	CTServer CommandType = 0
	CTUser   CommandType = 1
	CTExit   CommandType = 2
)

type Command struct {
	Command   CommandType
	Args      [5]string
	ArgsCount int
}

func parseCommandLine(input string) (Command, error) {
	var result Command

	var word [1024]rune
	var wordSize = 0
	var wordsCount = 0
	var isInsideQuotes = false
	procWord := func() error {
		if wordsCount == 0 {
			command := string(word[:wordSize])
			switch command {
			case "/user":
				result.Command = CTUser
			case "/server":
				result.Command = CTServer
			default:
				return fmt.Errorf("met unexpected command %s", command)
			}
		} else {
			arg := string(word[:wordSize])
			result.Args[wordsCount-1] = arg
		}
		return nil
	}

	for _, c := range input {
		switch c {
		case '"':
			isInsideQuotes = !isInsideQuotes
		case ' ':
			if isInsideQuotes {
				word[wordSize] = c
				wordSize++
				continue
			}
			if err := procWord(); err != nil {
				return Command{}, nil
			}

			wordSize = 0
			wordsCount++
			if wordsCount > 4 {
				return Command{}, fmt.Errorf("more then 5 args")
			}
		default:
			word[wordSize] = c
			wordSize++
			continue
		}
	}

	// End of the command process. Need to process last word, if exists
	if wordSize > 0 {
		if err := procWord(); err != nil {
			return Command{}, nil
		}
	}

	result.ArgsCount = wordsCount

	return result, nil
}

func processServerCommand(command Command, userManager *MockUserManager) error {
	if command.ArgsCount == 0 {
		return fmt.Errorf("expected server command args")
	}

	serverCommand := command.Args[0]
	switch serverCommand {
	case "create":
		if command.ArgsCount != 3 {
			return fmt.Errorf("expected 3 arguments, but met %d", command.ArgsCount)
		}
		userManager.createUser(command.Args[1], command.Args[2])
		fmt.Printf("User created\n")
	}

	return nil
}

func processUserCommand(command Command) error {
	return nil
}

type MockUserManager struct {
	Users        []MockUser
	UnusedChatID int64
}

type MockUser struct {
	ChatID       int64
	TelegramName string
	PlayerName   string
}

func (m *MockUserManager) createUser(telegramName string, playerName string) {
	m.Users = append(m.Users, MockUser{
		ChatID:       m.UnusedChatID,
		TelegramName: telegramName,
		PlayerName:   playerName,
	})
	m.UnusedChatID++
}

type Mock struct {
}
