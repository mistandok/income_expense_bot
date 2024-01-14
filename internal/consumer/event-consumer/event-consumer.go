package event_consumer

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"income_expense_bot/internal/events"
	"income_expense_bot/internal/lib/e"
	"log"

	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
	logger    *zerolog.Logger
}

func New(logger *zerolog.Logger, fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
		logger:    logger,
	}
}

func (c Consumer) Start(ctx context.Context) error {
	for {
		// TODO: учесть возможность ретрая.
		gotEvents, err := c.fetcher.Fetch(ctx, c.batchSize)
		if err != nil {
			c.logger.Err(e.Wrap("consumer", err))

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(ctx, gotEvents); err != nil {
			c.logger.Err(e.Wrap("handle events", err))

			continue
		}
	}
}

/*
1. Потеря событий: ретраи, возвращение в хранилище, фоллбэк, подтверждение
2. обработка всей пачки: останавливаться после первой же ошибки, либо ввести счетчик ошибок
3. Параллельная обработка событий. Потребуется структура WaitGroup.
*/
func (c Consumer) handleEvents(ctx context.Context, events []events.Event) error {
	for _, event := range events {
		c.logger.Info().Msg(fmt.Sprintf("got new event: %+v", event))
		log.Printf("got new event: %s", event.Text)

		// TODO: можно добавить механизм ретрая, или добавить механизм добавления во временное хранилище (фоллбэк, к примеру в оперативной памяти программы, или в редисе).
		// обработка событий идет одно за одним, но можно их обрабатывать асинхронно.
		if err := c.processor.Process(ctx, event); err != nil {
			c.logger.Err(e.Wrap("can't handle event", err))

			continue
		}
	}

	return nil
}
