package handlers

import (
	"fmt"
	"os"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

func RegistrationStartMessage(context runtime.Context) error {
	return context.Send("Путник, назови заветные руны и войди")
}

func RegistrationPassword(context runtime.Context, s *bot.Session) error {
	passwordPlayer := os.Getenv("BOT_USER_PASSWORD")
	passwordMaster := os.Getenv("BOT_MASTER_PASSWORD")

	password := context.MessageText()
	if password != passwordMaster && password != passwordPlayer {
		return RegistrationStartMessage(context)
	}

	s.RegistrationUser = &bot.UserData{
		ChatID:       context.ChatID(),
		TelegramName: context.FirstName(),
		PlayerName:   "",
		Faction:      &bot.Faction{},
	}
	s.RegistrationState = bot.RegistrationAwaitCodename

	switch password {
	case passwordPlayer:
		context.Send("Руны верны, Путник, войди же\nПопытай своё счастье в этом проклятом месте")
		s.RegistrationUser.Role = bot.RolePlayer
		return context.Send("Тебе понадоится новоё имя\nВыбирай с умом, никто не должен знать о твоём прошлом\nКак придумаешь напиши его на этом листке..")
	case passwordMaster:
		context.Send("Master password is correct, welcome!")
		s.RegistrationUser.Role = bot.RoleMaster
		return context.Send("Please, enter your Master Name")
	default:
		return fmt.Errorf("expected to handle password while registering user, but received this: %s", password)
	}
}

func RegistrationPlayerName(context runtime.Context, s *bot.Session) error {
	playerName := context.MessageText()

	s.RegistrationUser.PlayerName = playerName
	if s.RegistrationUser.Role == bot.RoleMaster {
		return registrationFinilize(context, s)
	}

	s.RegistrationState = bot.RegistrationAwaitFactionName
	text := `
	О, так говоришь ты из амбициозных, хочешь собрать людей?
	И как же ты назовёшь Фракцию? Какое имя скоро услышит весь город?`
	return context.Send(text)
}

func RegistrationFactionName(context runtime.Context, s *bot.Session) error {
	factionName := context.MessageText()

	s.RegistrationUser.Faction.Name = factionName
	s.RegistrationState = bot.RegistrationAwaitFactionDescription

	text := `
	Не могу понять по имени что это будет..
	Расскажи немного про свои планы, чего ты хочешь достичь?

	Твои люди будут защитниками правопорядка в городе? Будешь защищать невинных?
	Хочешь основать свой религиозный культ, и собрать верных последователей?
	Или ты друид, кто хочет вернуть зеленый цвет в этот каменный город?
	Или деньги всё о чём ты можешь думать и бесконечное количество способов их заработать?

	.. Не бери в голову что я сказал, не хочу останавливать твой полёт фантазии, ты сам(а) лучше знаешь чего хочешь

	Что это будет за фракция, расскажи мне вкратце..
	`
	return context.Send(text)
}

func RegistrationFactionDescription(context runtime.Context, s *bot.Session) error {
	factionDesc := context.MessageText()

	s.RegistrationUser.Faction.Description = factionDesc
	return registrationFinilize(context, s)
}

func registrationFinilize(context runtime.Context, s *bot.Session) error {
	id, err := db.CreateUser(s.DB, s.RegistrationUser)
	if err != nil {
		return err
	}

	_, err = db.CreateFaction(s.DB, s.RegistrationUser.Faction, id)
	if err != nil {
		return err
	}

	s.RegistrationState = bot.RegistrationFinished
	text := `
	Я услышал твою историю, ты можешь идти. Буду ждать чем обернется твоя жизнь

	Обустраивайся в этом городе, в нём хватит места для амбиций каждого

	И запомни главное: Вся магия этого мира держится на Вере людей в магию..
	`
	return context.Send(text)
}
