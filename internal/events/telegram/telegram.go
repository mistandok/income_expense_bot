package telegram

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"income_expense_bot/internal/clients/telegram"
	"income_expense_bot/internal/events"
	"income_expense_bot/internal/lib"
	"income_expense_bot/internal/lib/e"
)

type Worker struct {
	logger   *zerolog.Logger
	tgClient *telegram.Client
	offset   int
	//storage storage.Storage
}

type Meta struct {
	ChatID       int
	Username     string
	CallbackInfo *CallbackInfo
}

type CallbackInfo struct {
	ID   string
	Data *string
}

var (
	ErrUnknownEvent    = errors.New("unknown event type")
	ErrUnknownMetaType = errors.New("unknown meta type")
)

func New(logger *zerolog.Logger, client *telegram.Client) *Worker {
	return &Worker{logger: logger, tgClient: client}
}

func (w *Worker) Fetch(ctx context.Context, limit int) ([]events.Event, error) {
	updates, err := w.tgClient.Updates(ctx, w.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	fetchedEvents := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		fetchedEvents = append(fetchedEvents, event(u))
	}

	w.offset = updates[len(updates)-1].ID + 1

	return fetchedEvents, nil
}

func (w *Worker) Process(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.Message:
		return w.processMessage(ctx, event)
	case events.Callback:
		return w.processCallback(ctx, event)
	default:
		return e.Wrap("can't process message", ErrUnknownEvent)
	}
}

func (w *Worker) processMessage(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	w.logger.Info().Msg(fmt.Sprintf("try process message event %+v", event))
	//if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
	//	return e.Wrap("can't process message", err)
	//}

	return w.tgClient.SendMessage(ctx, telegram.OutgoingMessage{ChatId: meta.ChatID, Text: "обработали ваше сообщение"})
}

func (w *Worker) processCallback(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	w.logger.Info().Msg(fmt.Sprintf("try process callback event %+v", event))
	return w.tgClient.AnswerCallbackQuery(ctx, telegram.OutgoingCallbackMessage{CallbackQueryId: meta.CallbackInfo.ID, Text: lib.Pointer("response on callback query"), ShowAlert: lib.Pointer(true)})
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	if updType == events.Callback {
		res.Meta = Meta{
			CallbackInfo: &CallbackInfo{
				ID:   upd.CallbackQuery.ID,
				Data: upd.CallbackQuery.Data,
			},
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.CallbackQuery != nil {
		return events.Callback
	}

	if upd.Message != nil {
		return events.Message
	}

	return events.Unknown
}
