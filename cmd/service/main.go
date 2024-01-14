package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"income_expense_bot/internal/clients/telegram"
	"income_expense_bot/internal/lib"
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
	consumeBatchSize := mustFetchConsumeUpdatesBatchSize(envs)
	logLevel := mustFetchLogLevel(envs)

	ctx := context.Background()
	logger := setupZeroLog(logLevel, time.RFC3339)

	logger.Info().
		Str("tg_token", tgBotToken).
		Int("batch_size", consumeBatchSize).
		Msg("прочитали переменные окружения")

	outgoingMessage := telegram.OutgoingMessage{
		ChatId: 88700971,
		Text:   "add income",
		ReplyMarkup: telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					telegram.InlineKeyboardButton{
						Text:         "❌ delete",
						CallbackData: lib.Pointer("delete callback data"),
					},
					telegram.InlineKeyboardButton{
						Text:         "🔄 update",
						CallbackData: lib.Pointer("delete callback data"),
					},
				},
			},
		},
	}

	tgClient := telegram.New("api.telegram.org", tgBotToken)
	err = tgClient.SendMessage(ctx, outgoingMessage)
	if err != nil {
		logger.Fatal().Msg(fmt.Sprintf("fatal error %+w", err))
	}

}

func setupZeroLog(logLevel zerolog.Level, timeFormat string) zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = time.RFC3339

	return logger
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
