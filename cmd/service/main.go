package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"income_expense_bot/internal/clients/telegram"
	event_consumer "income_expense_bot/internal/consumer/event-consumer"
	telegram_worker "income_expense_bot/internal/events/telegram"
	"income_expense_bot/internal/lib/e"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	envs, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("не удалось прочитать переменные окружения")
	}
	tgBotToken := mustFetchTgBotToken(envs)
	tgBotHost := mustFetchTgBotHost(envs)
	consumeBatchSize := mustFetchConsumeUpdatesBatchSize(envs)
	logLevel := mustFetchLogLevel(envs)

	// инициализируем все объекты, необходимые для работы
	ctx := context.Background()
	logger := setupZeroLog(logLevel, time.RFC3339)
	tgClient := telegram.New(tgBotHost, tgBotToken)
	worker := telegram_worker.New(logger, tgClient)
	consumer := event_consumer.New(logger, worker, worker, consumeBatchSize)

	// стартуем бота
	logger.Info().Msg("service started")
	if err := consumer.Start(ctx); err != nil {
		logger.Fatal().Err(e.Wrap("service is stopped", err))
	}

}

func setupZeroLog(logLevel zerolog.Level, timeFormat string) *zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = time.RFC3339

	return &logger
}

func mustFetchLogLevel(envs map[string]string) zerolog.Level {
	name := "LOG_LEVEL"
	logLevelStr, ok := envs[name]
	if !ok {
		log.Fatal(fmt.Sprintf("не задан %s", name))
	}
	logLevelInt, err := strconv.Atoi(logLevelStr)
	if err != nil {
		log.Fatal(fmt.Sprintf("некорректное значение для %s", logLevelStr))
	}
	return zerolog.Level(logLevelInt)
}

func mustFetchTgBotToken(envs map[string]string) string {
	name := "TG_BOT_TOKEN"
	tgBotToken, ok := envs[name]
	if !ok {
		log.Fatal(fmt.Sprintf("не задан %s", name))
	}
	return tgBotToken
}

func mustFetchTgBotHost(envs map[string]string) string {
	name := "TG_BOT_HOST"
	tgBotHost, ok := envs[name]
	if !ok {
		log.Fatal(fmt.Sprintf("не задан %s", name))
	}
	return tgBotHost
}

func mustFetchConsumeUpdatesBatchSize(envs map[string]string) int {
	name := "CONSUME_UPDATES_BATCH_SIZE"
	batchSizeStr, ok := envs[name]
	if !ok {
		log.Fatal(fmt.Sprintf("не задан %s", name))
	}
	batchSizeInt, err := strconv.Atoi(batchSizeStr)
	if err != nil {
		log.Fatal(fmt.Sprintf("некорректное значение для %s", batchSizeStr))
	}
	return batchSizeInt
}
