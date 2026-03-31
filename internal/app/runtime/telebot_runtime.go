package runtime

import (
	tele "gopkg.in/telebot.v4"
)

// ---
// TeleContext
// ---
type TeleContext struct {
	ctx tele.Context
}

func (t *TeleContext) ChatID() int64 {
	return t.ctx.Chat().ID
}

func (t *TeleContext) FirstName() string {
	return t.ctx.Sender().FirstName
}

func (t *TeleContext) Callback() string {
	if t.ctx.Callback() == nil {
		return ""
	}
	return t.ctx.Callback().Data
}

func (t *TeleContext) Send(msg string, k ...Keyboard) error {
	menu := keyboardToMarkup(k...)
	return t.ctx.Send(msg, menu)
}

func (t *TeleContext) SendTo(id int64, text string, k ...Keyboard) error {
	menu := keyboardToMarkup(k...)
	_, err := t.ctx.Bot().Send(tele.ChatID(id), text, menu)
	return err
}

func (t *TeleContext) Respond() {
	t.ctx.Respond()
}

func (t *TeleContext) Args() []string {
	return t.ctx.Args()
}

func (t *TeleContext) MessageID() int64 {
	return int64(t.ctx.Message().ID)
}

func (t *TeleContext) MessageText() string {
	return t.ctx.Message().Text
}

func keyboardToMarkup(k ...Keyboard) *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}
	for _, keyboard := range k {
		var rows []tele.Row
		for _, row := range keyboard {
			var btns []tele.Btn
			for _, btn := range row {
				btns = append(btns, menu.Data(btn.Text, btn.Unique, btn.Data))
			}
			rows = append(rows, menu.Row(btns...))
		}
		menu.Inline(rows...)
	}

	return menu
}

// ---
// TelebotRuntime
// ---
type TelebotRuntime struct {
	Bot         *tele.Bot
	middlewares []Middleware
}

func (t *TelebotRuntime) HandleText(h Handler) {
	chain := t.apply(h)

	t.Bot.Handle(tele.OnText, func(ctx tele.Context) error {
		cont := &TeleContext{ctx: ctx}
		return chain(cont)
	})
}

func (t *TelebotRuntime) HandleCallback(h Handler) {
	chain := t.apply(h)

	t.Bot.Handle(tele.OnCallback, func(ctx tele.Context) error {
		cont := &TeleContext{ctx: ctx}
		return chain(cont)
	})
}

func (t *TelebotRuntime) HandleCommand(command string, h Handler) {
	chain := t.apply(h)

	t.Bot.Handle(command, func(ctx tele.Context) error {
		cont := &TeleContext{ctx: ctx}
		return chain(cont)
	})
}

func (t *TelebotRuntime) Start() error {
	t.Bot.Start()
	return nil
}

func (t *TelebotRuntime) Use(m Middleware) {
	t.middlewares = append(t.middlewares, m)
}

func (t *TelebotRuntime) apply(h Handler) Handler {
	for i := len(t.middlewares) - 1; i >= 0; i-- {
		h = t.middlewares[i](h)
	}
	return h
}
