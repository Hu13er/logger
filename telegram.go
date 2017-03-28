package logger

import (
	"sync"

	telegram "gopkg.in/telegram-bot-api.v4"
)

type telegSteam struct {
	bot     *telegram.BotAPI
	buf     []byte
	maxSize int
	chatID  int64
	mutex   sync.Mutex
}

func newTelegSteam(token string, chatID int64, maxSize int) (*telegSteam, error) {
	bot, err := telegram.NewBotAPI(token)
	return &telegSteam{bot: bot, buf: make([]byte, 0), maxSize: maxSize, chatID: chatID, mutex: sync.Mutex{}}, err
}

func (t *telegSteam) Write(buf []byte) (n int, err error) {

	t.mutex.Lock()
	defer t.mutex.Unlock()

	for i := 0; n < len(buf); i++ {
		min := t.maxSize
		if min > len(buf) {
			min = len(buf) - n
		}

		chunk := buf[i*t.maxSize : min]
		t.buf = append(t.buf, chunk...)
		if er := t.flush(); er != nil {
			err = er
			return
		}
		n += len(chunk)

	}
	return
}

func (t *telegSteam) flush() error {
	msg := telegram.NewMessage(t.chatID, string(t.buf))
	_, err := t.bot.Send(msg)
	t.buf = t.buf[0:0]
	return err
}
