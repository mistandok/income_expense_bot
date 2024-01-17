package event_consumer

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"income_expense_bot/internal/events"
	"time"
)

type EventFetcher interface {
	Fetch(ctx context.Context, limit int) ([]events.Event, error)
}

type EventProcessor interface {
	Process(ctx context.Context, e events.Event) error
}

type Consumer struct {
	fetcher   EventFetcher
	processor EventProcessor
	batchSize int
	logger    *zerolog.Logger
}

func New(logger *zerolog.Logger, fetcher EventFetcher, processor EventProcessor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
		logger:    logger,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		// TODO: учесть возможность ретрая.
		fetchedEvents, err := c.fetcher.Fetch(ctx, c.batchSize)
		if err != nil {
			c.logger.Err(err).Msg("fetch events error")

			continue
		}

		if len(fetchedEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(ctx, fetchedEvents); err != nil {
			c.logger.Err(err).Msg("handle events error")

			continue
		}
	}
}

/*
1. Потеря событий: ретраи, возвращение в хранилище, фоллбэк, подтверждение
2. обработка всей пачки: останавливаться после первой же ошибки, либо ввести счетчик ошибок
3. Параллельная обработка событий. Потребуется структура WaitGroup.
*/
func (c *Consumer) handleEvents(ctx context.Context, events []events.Event) error {
	for _, event := range events {
		c.logger.Info().Msg(fmt.Sprintf("got new event: %+v", event))

		// TODO: можно добавить механизм ретрая, или добавить механизм добавления во временное хранилище (фоллбэк, к примеру в оперативной памяти программы, или в редисе).
		// обработка событий идет одно за одним, но можно их обрабатывать асинхронно.
		if err := c.processor.Process(ctx, event); err != nil {
			c.logger.Err(err).Msg("can't handle event")

			continue
		}
	}

	return nil
}
