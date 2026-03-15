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

type MockRuntime struct {
	middlewares    []Middleware
	userManager    MockUserManager
	handlerManager HandlerManager
	messageManager MessageManager
}

type HandlerManager struct {
	HandleCallback Handler
	HandleText     Handler
	HandleCommand  Handler
}

// MockUserManager We are mocking telegram users.
// The most important part of the mock is creating ChatID, that will be used
type MockUserManager struct {
	Users        []MockUser
	UnusedChatID int64
}

type MessageManager struct {
	Messages        []Message
	UnusedMessageID int64
}

type Message struct {
	ID     int64
	Text   string
	ChatID int64
	Name   string
}

type MockUser struct {
	ChatID       int64
	TelegramName string
	PlayerName   string
}

type Command struct {
	Command   CommandType
	Args      [5]string
	ArgsCount int

	// is argument in quotes: /user John "Create Sword" -> {false, false, true}
	IsTextArgs [5]bool
}

type CommandType int

const (
	CTServer CommandType = 0
	CTUser   CommandType = 1
	CTExit   CommandType = 2
)

func (c *MockContext) ChatID() int64 {
	return c.chatID
}
func (c *MockContext) FirstName() string {
	return c.firstName
}
func (c *MockContext) Callback() string {
	return c.callback
}

func (c *MockContext) Send(text string, k ...Keyboard) error {
	fmt.Printf("[BotReply]: %s\n", text)
	for _, keyboard := range k {
		for _, row := range keyboard {
			for _, btn := range row {
				fmt.Printf("%s [%s | %s]\n", btn.Text, btn.Unique, btn.Data)
			}
		}
	}
	return nil
}
func (c *MockContext) SendTo(id int64, text string, k ...Keyboard) error {
	fmt.Print("SendTo received")
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

func NewMockRuntime() *MockRuntime {
	return &MockRuntime{
		middlewares:    make([]Middleware, 0),
		userManager:    MockUserManager{},
		handlerManager: HandlerManager{},
		messageManager: MessageManager{},
	}
}

func (mr *MockRuntime) Use(m Middleware) {
	mr.middlewares = append(mr.middlewares, m)
}
func (mr *MockRuntime) HandleText(h Handler) {
	chain := mr.apply(h)
	mr.handlerManager.HandleText = chain
}
func (mr *MockRuntime) HandleCallback(h Handler) {
	chain := mr.apply(h)
	mr.handlerManager.HandleCallback = chain
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
				fmt.Printf("%s\n", err)
			}
		case CTUser:
			err := processUserCommand(command, &mr.userManager, &mr.handlerManager, &mr.messageManager)
			if err != nil {
				fmt.Printf("%s\n", err)
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

func (m *MockUserManager) createUser(telegramName string, playerName string) {
	m.Users = append(m.Users, MockUser{
		ChatID:       m.UnusedChatID,
		TelegramName: telegramName,
		PlayerName:   playerName,
	})
	m.UnusedChatID++
}

func (m *MockUserManager) getUser(name string) *MockUser {
	for _, u := range m.Users {
		if u.PlayerName == name {
			return &u
		}
		if u.TelegramName == name {
			return &u
		}
	}
	return nil
}

func (m *MessageManager) createMessage(text string, chatID int64, name string) int64 {
	msg := Message{
		ID:     m.UnusedMessageID,
		Text:   text,
		ChatID: chatID,
		Name:   name,
	}
	m.Messages = append(m.Messages, msg)
	m.UnusedMessageID++
	return msg.ID
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
			if isInsideQuotes {
				result.IsTextArgs[wordsCount-1] = true
			}
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
	default:
		return fmt.Errorf("server unknown command, %s", serverCommand)
	}

	return nil
}

func processUserCommand(command Command, userManager *MockUserManager, handlers *HandlerManager, messages *MessageManager) error {
	if command.ArgsCount < 2 {
		return fmt.Errorf("expected user command args")
	}

	name := command.Args[0]
	user := userManager.getUser(name)
	if user == nil {
		return fmt.Errorf("user %s does not exist", name)
	}

	isText := command.IsTextArgs[1]
	if isText {
		// TODO: Handle Commands through text. When we have command in text
		// /user John "/elevate"
		text := command.Args[1]
		id := messages.createMessage(text, user.ChatID, user.TelegramName)
		err := handlers.HandleText(&MockContext{
			chatID:      user.ChatID,
			firstName:   user.PlayerName,
			callback:    "",
			args:        []string{},
			messageID:   id,
			messageText: text,
		})
		return err
	}

	cb := command.Args[1]
	err := handlers.HandleCallback(&MockContext{
		chatID:      user.ChatID,
		firstName:   user.PlayerName,
		callback:    cb,
		args:        []string{},
		messageID:   0,
		messageText: "",
	})
	return err
}
