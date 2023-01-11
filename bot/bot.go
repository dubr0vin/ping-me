package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"pingMe/storage"
	"strconv"
	"strings"
)

var (
	bot *tgbotapi.BotAPI
)

func send(text string, user int64) error {
	_, err := bot.Send(tgbotapi.NewMessage(user, text))
	return err
}

func getFromStorageOrFake(userId int64) storage.UserData {
	user, ok := data.Users[userId]
	if !ok {
		user = storage.UserData{
			Imports: make([]storage.ImportData, 0),
			Minutes: make([]uint32, 0),
		}
	}
	return user
}
func botImport(userId int64, args []string) {
	if len(args) == 0 {
		if err := send("Not enough arguments", userId); err != nil {
			fmt.Println(err)
		}
	} else {
		user := getFromStorageOrFake(userId)
		user.Imports = append(user.Imports,
			storage.ImportData{
				Id:     args[0],
				Params: args[1:],
			},
		)
		data.Users[userId] = user
		data.Save()
		if err := send("Done", userId); err != nil {
			fmt.Println(err)
		}
	}
}

func botImports(userId int64, args []string) {
	user := getFromStorageOrFake(userId)
	message := ""
	for i, userImport := range user.Imports {
		message += strconv.Itoa(i) + ") " + userImport.Id + " " + strings.Join(userImport.Params, " ")
	}
	if err := send(message, userId); err != nil {
		fmt.Println(err)
	}
}

func botRemoveImport(userId int64, args []string) {
	if len(args) == 0 {
		if err := send("Not enough arguments", userId); err != nil {
			fmt.Println(err)
		}
        return
	}
	result, err := strconv.Atoi(args[0])
	if err != nil {
		if err := send("First argument must be integer", userId); err != nil {
			fmt.Println(err)
        }
        return
	}
	user := getFromStorageOrFake(userId)
	if !(0 <= result && result < len(user.Imports)) {
		if err := send("First argument must in [0; "+strconv.Itoa(len(user.Imports))+")", userId); err != nil {
			fmt.Println(err)
        }
        return
	}
	user.Imports = append(user.Imports[0:result], user.Imports[result+1:]...)
	data.Users[userId] = user
	data.Save()
    if err := send("Done", userId); err != nil {
        fmt.Println(err)
    }
}
func botOffset(userId int64, args []string) {
    if len(args) == 0 {
        if err := send("Not enough arguments", userId); err != nil {
            fmt.Println(err)
        }
        return
    }
    result, err := strconv.Atoi(args[0])
    if err != nil {
        if err := send("First argument must be integer", userId); err != nil {
            fmt.Println(err)
        }
        return
    }
    user := getFromStorageOrFake(userId)
    if !(0 <= result) {
        if err := send("First argument must in [0; +inf)", userId); err != nil {
            fmt.Println(err)
        }
        return
    }
    user.Minutes = append(user.Minutes, uint32(result))
    data.Users[userId] = user
    data.Save()
    if err := send("Done", userId); err != nil {
        fmt.Println(err)
    }
}

func botOffsets(userId int64, args []string) {
    user := getFromStorageOrFake(userId)
    message := ""
    for i, offset := range user.Minutes {
        message += strconv.Itoa(i) + ") " + strconv.Itoa(int(offset)) + "\n"
    }
    if err := send(message, userId); err != nil {
        fmt.Println(err)
    }
}

func botRemoveOffset(userId int64, args []string) {
    if len(args) == 0 {
        if err := send("Not enough arguments", userId); err != nil {
            fmt.Println(err)
        }
        return
    }
    result, err := strconv.Atoi(args[0])
    if err != nil {
        if err := send("First argument must be integer", userId); err != nil {
            fmt.Println(err)
        }
        return
    }
    user := getFromStorageOrFake(userId)
    if !(0 <= result && result < len(user.Minutes)) {
        if err := send("First argument must in [0; "+strconv.Itoa(len(user.Minutes))+")", userId); err != nil {
            fmt.Println(err)
        }
        return
    }
    user.Minutes = append(user.Minutes[0:result], user.Minutes[result+1:]...)
    data.Users[userId] = user
    data.Save()
    if err := send("Done", userId); err != nil {
        fmt.Println(err)
    }
}

func RunBot() {
	var err error
	bot, err = tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	setupData()
	setupTicker()
	for update := range updates {
		if update.Message == nil || update.Message.Text == "" {
			continue
		}
		items := strings.Split(update.Message.Text, " ")
		args := items[1:]
		switch items[0] {
		case "/import":
			botImport(update.Message.From.ID, args)
		case "/imports":
			botImports(update.Message.From.ID, args)
		case "/remove_import":
			botRemoveImport(update.Message.From.ID, args)
        case "/add_offset":
            botOffset(update.Message.From.ID, args)
        case "/offsets":
            botOffsets(update.Message.From.ID, args)
        case "/remove_offset":
            botRemoveOffset(update.Message.From.ID, args)
		}
        if strings.Contains(update.Message.Text, "#bug") {
            if _, err := bot.CopyMessage(tgbotapi.NewCopyMessage(-1001807236405, update.Message.Chat.ID, update.Message.MessageID)); err != nil {
                fmt.Println(err)
            }
        }
	}
}
