package main

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/sipki-tech/database/connectors"
	"log"
	"os"
	"os/signal"
	"tgbot/cmd/internal/adapters/repo"
	"tgbot/cmd/internal/api"
	"tgbot/cmd/internal/app"
)

// Send any text message to the bot after the bot has been started

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	reponew, err := repo.New(ctx, repo.Config{
		Postgres: connectors.Raw{
			Query: "postgres://user_svc:user_pass@localhost:5432/user_db?sslmode=disable",
		},
		MigrateDir: "./cmd/migrate",
		Driver:     "postgres",
	})
	if err != nil {
		panic(err)
	}
	appnew := app.New(reponew)
	b := api.New(appnew)
	log.Print("starting server")
	b.Start(ctx)
}

// Квиз бот список вопросов из бд, у каждого впроса варианты ответов в бд в другой бд, один вариант должен быть
// помечен правильным (тру фолс) регистрация в боте при старте выбор категории вопросов
// начисление баллов за правильное кол-во ответов,
// фио номер телефона возраст
