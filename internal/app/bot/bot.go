package bot

import (
	"database/sql"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	tele "gopkg.in/telebot.v4"
)

type BotData struct {
	DB *sql.DB

	UserActiveSessions    map[int64]FlowSession
	UserRegistrationCache map[int64]*UserData
}

func BotInit(db *sql.DB) *BotData {
	bot := &BotData{
		DB:                    db,
		UserActiveSessions:    make(map[int64]FlowSession),
		UserRegistrationCache: make(map[int64]*UserData),
	}
	return bot
}

type FlowSession interface {
	Name() string
	IsSupportedCallback(string) bool
	IsDone() bool
	DispatchCallback(context tele.Context, cbUnique string, cbData string) error
	DispatchText(context tele.Context) error
}

type UserData struct {
	ChatID       int64
	TelegramName string
	PlayerName   string
	Role         UserRole
	Faction      *Faction
}

func (user *UserData) Recipient() string {
	return strconv.FormatInt(int64(user.ChatID), 10)
}

func (b *BotData) ClearUserCache(chatID int64) {
	delete(b.UserActiveSessions, chatID)
	delete(b.UserRegistrationCache, chatID)
}

func (b *BotData) GetUserSession(chatID int64) FlowSession {
	session, ok := b.UserActiveSessions[chatID]
	if !ok {
		return nil
	}
	return session
}

type Message struct {
	ID        int
	Title     string
	MessageID string
	ChatID    int64 // from which chat to copy
	Text      string
}

func (msg Message) MessageSig() (string, int64) {
	return msg.MessageID, msg.ChatID
}

func (msg Message) MessageHash() string {
	return msg.MessageID + strconv.FormatInt(msg.ChatID, 10)
}

type MessageTransaction struct {
	ID        int
	CreatedAt time.Time
	From      tele.ChatID
	To        tele.ChatID

	Message *Message
}

type MasterRequest struct {
	ID           int
	CreatedAt    time.Time
	To           tele.ChatID
	TextRequest  string
	TextResponse string
	IsAnswered   bool

	RollRequests []*RollRequest
}

type RollRequest struct {
	ID         int
	CreatedAt  time.Time
	Title      string
	DiceCount  int
	DiceSides  int
	RollResult int
}

type Faction struct {
	ID          int
	Name        string
	Description string
	Resources   string
}
