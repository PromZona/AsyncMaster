package bot

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/PromZona/AsyncMaster/internal/app/runtime"
	_ "github.com/lib/pq"
)

type BotData struct {
	DB *sql.DB

	Sessions map[int64]*Session
}

func BotInit(db *sql.DB) *BotData {
	bot := &BotData{
		DB:       db,
		Sessions: make(map[int64]*Session),
	}
	return bot
}

type RegistrationState int

const (
	RegistrationNotActive               RegistrationState = 0
	RegistrationAwaitPassword           RegistrationState = 1
	RegistrationAwaitCodename           RegistrationState = 2
	RegistrationAwaitFactionName        RegistrationState = 3
	RegistrationAwaitFactionDescription RegistrationState = 4
	RegistrationFinished                RegistrationState = 5
)

type Session struct {
	DB            *sql.DB
	PreviousRoute *Route

	// factions
	FactionUpdateUser *UserData

	// Registration
	RegistrationState RegistrationState
	RegistrationUser  *UserData

	// Message Handler
	DraftMessage     *Message
	DraftTransaction *MessageTransaction
	IsSendEveryone   bool

	// Master Request
	MasterRequest *MasterRequest
	RollRequests  []*RollRequest
	Resipients    []int64
}

func NewSession(db *sql.DB) *Session {
	return &Session{
		DB:                db,
		PreviousRoute:     nil,
		RegistrationState: RegistrationNotActive,
		RegistrationUser:  nil,
		DraftMessage:      nil,
		DraftTransaction:  nil,
		IsSendEveryone:    false,
		MasterRequest:     nil,
		RollRequests:      nil,
		Resipients:        nil,
	}
}

func (b *BotData) GetSession(chatID int64) *Session {
	s, ok := b.Sessions[chatID]
	if !ok {
		return nil
	}
	return s
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
	delete(b.Sessions, chatID)
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
	From      int64
	To        []int64

	Message *Message
}

type MasterRequestState int

const (
	MRUnasnwered      MasterRequestState = 0
	MRAnswered        MasterRequestState = 1
	MRCheckedByMaster MasterRequestState = 2
)

type MasterRequest struct {
	ID           int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	To           int64
	TextRequest  string
	TextResponse string
	State        MasterRequestState

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

type Handler func(runtime.Context, *Session) error
type HandlerID string

type Route struct {
	Callback HandlerID
	Handler  Handler

	NextPossibleCallbacks []HandlerID
}

const (
	HID_cancel HandlerID = "cancel"

	// Factions
	HID_factions_list                  HandlerID = "factions_list"
	HID_factions_update                HandlerID = "factions_update"
	HID_factions_update_player         HandlerID = "factions_update_player"
	HID_factions_update_resources      HandlerID = "factions_update_resources"
	HID_factions_update_resources_text HandlerID = "factions_update_resources_text"

	// Message Create
	HID_message_send          HandlerID = "message_send"
	HID_message_send_everyone HandlerID = "message_send_everyone"
	HID_message_player_name   HandlerID = "message_player_name"
	HID_message_text          HandlerID = "message_text"
	HID_message_title_yes     HandlerID = "message_title_yes"
	HID_message_title_no      HandlerID = "message_title_no"
	HID_message_title_text    HandlerID = "message_title_text"

	// Message Get
	HID_message_list HandlerID = "message_list"
	HID_message_get  HandlerID = "message_get"

	// Master Request Create
	HID_mr_send          HandlerID = "master_request_send"
	HID_mr_send_everyone HandlerID = "master_request_send_everyone"
	HID_mr_resipient     HandlerID = "master_request_resipient"
	HID_mr_text          HandlerID = "master_request_text"
	HID_mr_dice_yes      HandlerID = "master_request_dice_yes"
	HID_mr_dice_no       HandlerID = "master_request_dice_no"
	HID_mr_dice_text     HandlerID = "master_request_dice_text"

	// Master Request Answer
	HID_mr_answer      HandlerID = "master_request_answer"
	HID_mr_answer_text HandlerID = "master_request_answer_text"
	HID_mr_answer_roll HandlerID = "master_request_answer_roll"

	// Master Request Get
	HID_mr_player_get       HandlerID = "master_request_player_get"
	HID_mr_master_get       HandlerID = "master_request_master_get"
	HID_mr_master_mark_read HandlerID = "master_request_master_mark_read"
)
