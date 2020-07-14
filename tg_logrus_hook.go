package tghook

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// Hook structure
type Hook struct {
	Client   *tgbotapi.BotAPI
	ClientID int64
	levels   []logrus.Level
}

// NewHook init
func NewHook(apiKey string, clientID int64, levels []logrus.Level) (*Hook, error) {
	client, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return nil, err
	}

	return &Hook{
		Client:   client,
		ClientID: clientID,
		levels:   levels,
	}, nil
}

// Fire routine
func (hook *Hook) Fire(logEntry *logrus.Entry) error {
	var notifyErr string

	if err, ok := logEntry.Data["error"].(error); ok {
		notifyErr = err.Error()
	} else {
		notifyErr = logEntry.Message
	}

	msg := tgbotapi.MessageConfig{}
	msg.ChatID = hook.ClientID

	msg.Text = fmt.Sprintf(
		"%s: %s",
		strings.ToUpper(logEntry.Level.String()),
		notifyErr,
	)

	logEntry.Logger.Debug(msg)
	result, err := hook.Client.Send(msg)
	logEntry.Logger.Debug(result)

	if err != nil {
		logEntry.Logger.WithFields(logrus.Fields{
			"source": "telegram",
			"error":  err,
		}).Warn("Failed to send error to Telegram")
	}

	return nil
}

// Levels setting
func (hook *Hook) Levels() []logrus.Level {
	return hook.levels
}
