package api

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/fsm"
	"github.com/google/uuid"
	"tgbot/cmd/internal/app"
)

type application interface {
	CreateFolder(ctx context.Context, folder app.Folder) (*app.Folder, error)
	DeleteFolder(ctx context.Context, id uuid.UUID) error
	GetAllFolders(ctx context.Context, id int64) ([]app.Folder, error)
	OpenFolder(ctx context.Context, id uuid.UUID) ([]app.Card, error)
	CreateNameCard(ctx context.Context, card app.Card) (*app.Card, error)
	CreateTextCard(ctx context.Context, card app.Card) error
	DeleteCard(ctx context.Context, name string) error
	CreateUser(ctx context.Context, us app.User) error
	GetCard(ctx context.Context, id uuid.UUID) (*app.Card, error)
}

type api struct {
	app application
	f   *fsm.FSM
	b   *bot.Bot
}

func New(app application) *bot.Bot {
	api := &api{app,
		&fsm.FSM{},
		&bot.Bot{},
	}
	api.f = fsm.New(
		stateDefault,
		map[fsm.StateID]fsm.Callback{
			stateCreateFolder:   api.CreateFolderCallbackHandler,
			stateCreateCardName: api.CreateCardNameCallbackHandler,
			stateDeleteCard:     api.DeleteCardCallbackHandler,
		},
	)
	opts := []bot.Option{
		//bot.WithDefaultHandler(api.Start),
		bot.WithMessageTextHandler("/start", bot.MatchTypePrefix, api.Start),
		bot.WithCallbackQueryDataHandler("1", bot.MatchTypeExact, api.CreateFolder),
		bot.WithMessageTextHandler("", bot.MatchTypeContains, api.DefaultHandler),
		bot.WithCallbackQueryDataHandler("2", bot.MatchTypeExact, api.CreateCardName),
		bot.WithCallbackQueryDataHandler("3", bot.MatchTypeExact, api.BackStart),
		bot.WithCallbackQueryDataHandler("4", bot.MatchTypeExact, api.GetAllFolders),
		bot.WithCallbackQueryDataHandler("i", bot.MatchTypePrefix, api.GetFolder),
		bot.WithCallbackQueryDataHandler("5", bot.MatchTypeExact, api.DeleteFolder),
		bot.WithCallbackQueryDataHandler("6", bot.MatchTypeExact, api.GetCards),
		bot.WithCallbackQueryDataHandler("7", bot.MatchTypeExact, api.DeleteCard),
		bot.WithCallbackQueryDataHandler("8", bot.MatchTypeExact, api.Random),
		bot.WithCallbackQueryDataHandler("9", bot.MatchTypePrefix, api.getText),
		bot.WithCallbackQueryDataHandler("10", bot.MatchTypeExact, api.nextCard),
		//bot.WithCallbackQueryDataHandler("", bot.MatchTypePrefix, api.callbackHandler),
	}

	b, err := bot.New("", opts...)
	if nil != err {
		panic(err)
	}
	api.b = b
	return b
}
