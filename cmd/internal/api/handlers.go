package api

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/fsm"

	"github.com/google/uuid"
	"log"
	"math/rand"
	"strings"
	"tgbot/cmd/internal/app"
)

const (
	stateDefault        fsm.StateID = "default"
	stateCreateFolder   fsm.StateID = "createFolder"
	stateCreateCardName fsm.StateID = "createCardName"
	stateCreateTextCard fsm.StateID = "createTextCard"
	stateDeleteCard     fsm.StateID = "deleteCard"
)

const startMessage = "Описание бота"

var RandomID = make(map[int64]Str)

type Str struct {
	Counter  int
	CardList []uuid.UUID
}

//func (a *api) callbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
//	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
//		CallbackQueryID: update.CallbackQuery.ID,
//		ShowAlert:       false,
//	})
//	b.SendMessage(ctx, &bot.SendMessageParams{
//		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
//		Text:   "ХУЙ ЗНАЕТ ЧТО ЭТО",
//	})
//}

func (a *api) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	a.app.CreateUser(ctx, app.User{TgUserID: chatID, Name: update.Message.Chat.Username})

	kb := &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			{Text: "Создать папку", CallbackData: "1"},
			{Text: "Уже тут был", CallbackData: "3"},
		},
	},
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        startMessage,
		ReplyMarkup: kb,
	})
}

//
//func (a *api) CreateUser(ctx context.Context, b *bot.Bot, update *models.Update) {
//	a.app.CreateUser(ctx, app.User{TgUserID: update.Message.Chat.ID, Name: update.Message.Chat.Username})
//	b.SendMessage(ctx, &bot.SendMessageParams{
//		ChatID: update.Message.Chat.ID,
//		Text:   "Привет!",
//	})
//}

func (a *api) CreateFolder(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	a.f.Transition(chatID, stateCreateFolder, chatID, ctx)
}

func (a *api) CreateFolderCallbackHandler(f *fsm.FSM, args ...any) {
	chatID := args[0].(int64)
	a.b.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Введи название:",
	})
}

func (a *api) DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	userStatus := a.f.Current(chatID)
	switch userStatus {
	case stateCreateFolder:
		fold, err := a.app.CreateFolder(ctx, app.Folder{
			Name:     update.Message.Text,
			TgUserID: chatID,
		})
		if err != nil {
			a.b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   fmt.Sprintf("Папка с именем %s уже существует.\nВведите другое Название:", update.Message.Text),
			})
			return
		}

		a.f.Set(chatID, "folderID", fold.ID)
		kb := &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Создать карточку", CallbackData: "2"},
			},
		},
		}
		//Написать текст что ты уже в этой папке и можешь писать карточки
		a.b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      chatID,
			Text:        "Папка " + fold.Name + " создана",
			ReplyMarkup: kb,
		})
	case stateCreateCardName:
		folderID, ok := a.f.Get(chatID, "folderID")
		if !ok {
			log.Println(ok)
			return
		}
		card, err := a.app.CreateNameCard(ctx, app.Card{Name: update.Message.Text, FolderID: folderID.(uuid.UUID)})
		if err != nil {
			a.b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   fmt.Sprintf("Карточка с именем %s уже существует.\nВведите другое Название:", update.Message.Text),
			})
			return
		}
		a.f.Set(chatID, "cardID", card.ID)

		a.b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Введи текст карточки: ",
		})
		a.f.Transition(chatID, stateCreateTextCard, chatID, ctx)
	case stateCreateTextCard:
		folderID, ok := a.f.Get(chatID, "folderID")
		if !ok {
			log.Println(ok)
			return
		}
		folderID = folderID.(uuid.UUID)
		cardID, ok := a.f.Get(chatID, "cardID")
		if !ok {
			log.Println(ok)
			return
		}
		err := a.app.CreateTextCard(ctx, app.Card{Text: update.Message.Text, ID: cardID.(uuid.UUID)})
		if err != nil {
			log.Println(err)
			return
		}
		kb := &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Создать карточку", CallbackData: "2"},
				{Text: "В начало", CallbackData: "3"},
				{Text: "Удалить карточку", CallbackData: "6"},
				{Text: "Назад", CallbackData: fmt.Sprintf("i %s", folderID)},
			},
		},
		}
		a.b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "Card " + " создана",
			ReplyMarkup: kb,
		})
	case stateDeleteCard:
		err := a.app.DeleteCard(ctx, update.Message.Text)
		if err != nil {
			log.Println(err)
			return
		}
		kb := &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Получить созданные папки", CallbackData: "4"},
				{Text: "Создать карточку", CallbackData: "2"},
				{Text: "Удалить карточку", CallbackData: "7"},
			},
		}}
		a.b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "Карточка удалена.Куда дальше?",
			ReplyMarkup: kb,
		})
	default:
		a.b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Не пиши просто так!",
		})
	}
}

func (a *api) CreateCardName(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	a.f.Transition(chatID, stateCreateCardName, chatID, ctx)
}

func (a *api) CreateCardNameCallbackHandler(f *fsm.FSM, args ...any) {
	chatID := args[0].(int64)
	a.b.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Введите название карточки: ",
	})
}

func (a *api) BackStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			{Text: "Создать папку", CallbackData: "1"},
			{Text: "Получить созданные папки", CallbackData: "4"},
		},
	},
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        "Выбери действие: ",
		ReplyMarkup: kb,
	})
}

func (a *api) GetAllFolders(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	all, err := a.app.GetAllFolders(ctx, chatID)
	if err != nil {
		log.Println(err)
		return
	}
	var inlineKeyboard [][]models.InlineKeyboardButton
	for _, folder := range all {
		button := models.InlineKeyboardButton{
			Text:         folder.Name,
			CallbackData: fmt.Sprintf("i %s", folder.ID),
		}
		inlineKeyboard = append(inlineKeyboard, []models.InlineKeyboardButton{button})
	}
	kb := &models.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "Выберите папку",
		ReplyMarkup: kb,
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func (a *api) GetFolder(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	id := strings.Trim(update.CallbackQuery.Data, "i ")
	folderID, err := uuid.Parse(id)
	if err != nil {
		log.Println(err)
		return
	}
	a.f.Set(chatID, "folderID", folderID)
	kb := models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			{Text: "Рандомизировать", CallbackData: "8"},
			{Text: "Создать карточку", CallbackData: "2"},
			{Text: "Удалить папку", CallbackData: "5"},
			{Text: "Посмотреть карточки/Удалить", CallbackData: "6"},
		},
	}}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        "Ты в папке, выбери действие: ",
		ReplyMarkup: kb,
	})
}

func (a *api) DeleteFolder(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	folderID, ok := a.f.Get(chatID, "folderID")
	if !ok {
		log.Println(ok)
	}
	err := a.app.DeleteFolder(ctx, folderID.(uuid.UUID))
	if err != nil {
		log.Println(err)
		return
	}
	kb := models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			{Text: "В начало", CallbackData: "3"},
		},
	}}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        "Папка удалена.",
		ReplyMarkup: kb,
	})
}

func (a *api) GetCards(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	folderID, ok := a.f.Get(chatID, "folderID")
	if !ok {
		log.Println(ok)
	}
	cards, err := a.app.OpenFolder(ctx, folderID.(uuid.UUID))
	if err != nil {
		log.Println(err)
		return
	}

	for _, e := range cards {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   fmt.Sprintf("%s - %s \n", e.Name, e.Text),
		})
	}
	kb := models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			{Text: "Удалить карточку", CallbackData: "7"},
			{Text: "Back", CallbackData: "3"},
		},
	}}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "Действия с карточками: ",
		ReplyMarkup: kb,
	})
}

func (a *api) DeleteCard(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	a.f.Transition(chatID, stateDeleteCard, chatID, ctx)
}

func (a *api) DeleteCardCallbackHandler(f *fsm.FSM, args ...any) {
	chatID := args[0].(int64)
	a.b.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Введи название удаляемой карточки:",
	})
}

func (a *api) Random(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	folderID, ok := a.f.Get(chatID, "folderID")
	if !ok {
		log.Println(ok)
	}
	cards, err := a.app.OpenFolder(ctx, folderID.(uuid.UUID))
	if err != nil {
		log.Println(err)
		return
	}
	var idCard []uuid.UUID
	for _, e := range cards {
		idCard = append(idCard, e.ID)
	}
	rand.Shuffle(len(idCard), func(i, j int) {
		idCard[i], idCard[j] = idCard[j], idCard[i]
	})
	log.Println(idCard)
	RandomID[update.CallbackQuery.Message.Message.Chat.ID] = Str{
		Counter:  0,
		CardList: idCard,
	}
	a.sendCard(ctx, b, update)
}

func (a *api) sendCard(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	state := RandomID[chatID]
	cardID := state.CardList[state.Counter]
	card, err := a.app.GetCard(ctx, cardID)
	if err != nil {
		log.Println(err)
		return
	}
	kb := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Показать описание", CallbackData: "9"},
			},
		},
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        fmt.Sprintf("Карточка: %s", card.Name),
		ReplyMarkup: kb,
	})
}

func (a *api) getText(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	state := RandomID[chatID]
	cardID := state.CardList[state.Counter]
	card, err := a.app.GetCard(ctx, cardID)
	if err != nil {
		log.Println(err)
		return
	}
	kb := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Следующая карточка", CallbackData: "10"},
				{Text: "Закончить", CallbackData: "3"},
			},
		},
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        fmt.Sprintf("Описание: %s", card.Text),
		ReplyMarkup: kb,
	})
}

func (a *api) nextCard(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	state := RandomID[chatID]
	state.Counter++
	if state.Counter >= len(state.CardList) {
		state.Counter = 0
	}
	RandomID[chatID] = state
	a.sendCard(ctx, b, update)
}
