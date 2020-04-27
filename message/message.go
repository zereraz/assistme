package message

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sahilm/fuzzy"
	"github.com/zereraz/assistme/db"
	l "github.com/zereraz/assistme/log"
	"github.com/zereraz/assistme/user"
)

func ListenToCommands(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.InlineQuery != nil {
			handleInline(bot, update.InlineQuery)
		}

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "") //update.Message.Text)
		if update.Message.IsCommand() {
			l.Log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg.Text = fmt.Sprintf("Got command %s", update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(msg); err != nil {
				l.Log.Fatalf("Error while sending message: %v", err)
			}
		} else {
			continue
		}
	}
}

func handleInline(bot *tgbotapi.BotAPI, inlineQuery *tgbotapi.InlineQuery) error {
	_, err := db.GetDb()

	if err != nil {
		return err
	}

	username := inlineQuery.From.UserName

	currQuery := inlineQuery.Query

	// Fetch User object if exists
	u, err := user.FetchUser(username)

	if err == user.ErrUserNotFound {
		name := strings.TrimSpace(fmt.Sprintf("%s %s", inlineQuery.From.FirstName, inlineQuery.From.LastName))
		u, err = user.NewUser(name, username, int64(inlineQuery.From.ID), nil)
		if err != nil {
			return err
		}
		err = u.AddToDb()
		if err != nil {
			return err
		}
	}

	var categories []string

	for _, c := range u.Categories {
		categories = append(categories, c.Name)
	}

	matches := fuzzy.Find(currQuery, categories)

	if currQuery == "" {
		for _, c := range categories {
			matches = append(matches, fuzzy.Match{Str: c})
		}
	}

	// first check if this is a known category
	articles := make([]interface{}, len(matches))
	for i, match := range matches {
		article := tgbotapi.NewInlineQueryResultArticle(inlineQuery.ID+strconv.Itoa(i), match.Str, match.Str)
		article.Description = match.Str
		articles[i] = article
	}

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     10,
		Results:       articles,
	}

	if _, err := bot.Request(inlineConf); err != nil {
		return err
	}
	return nil
}
