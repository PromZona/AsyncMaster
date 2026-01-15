package _example

import (
	tele "gopkg.in/telebot.v4"
)

func handleStartFlow(context tele.Context, s *Session) error {
	if s.UserState != FlowStart {
		return context.Send("This action is not available right now, finish previous action first")
	}

	s.UserState = AwaitTextMessage
	return context.Send("Example")
}

func handleTextMessage(context tele.Context, s *Session) error {
	if s.UserState != AwaitTextMessage {
		return context.Send("This action is not available right now, finish previous action first")
	}

	context.Send(context.Text())
	return finilize(context, s)
}

func finilize(context tele.Context, s *Session) error {
	s.Done = true
	return context.Send("Example finish")
}
