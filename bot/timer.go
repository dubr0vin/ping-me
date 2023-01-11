package bot

import (
	"errors"
	"fmt"
	"pingMe/imports"
	csv "pingMe/imports/csv"
	"pingMe/storage"
    "strings"
    "time"
)

func setupTicker() {
	ticker := time.NewTicker(1 * time.Minute)
	quit := make(chan struct{})
	go func() {
		everyMinute()
		for {
			select {
			case <-ticker.C:
				everyMinute()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

var (
	userTasks = make(map[int64][]*time.Timer)
	data      *storage.Data
)

func setupData() {
	var err error
	data, err = storage.Load()
	if err != nil {
		panic(err)
	}
}

func getEvents(importData storage.ImportData) ([]imports.Event, error) {
	var getter imports.Import
	var err error
	switch importData.Id {
	case "csv":
		getter, err = csv.NewImport(importData.Params)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unknown type " + importData.Id)
	}
	return getter.GetEvents()
}
func everyMinute() {
	for user, userDatas := range data.Users {

		events := make([]imports.Event, 0)
		for _, importData := range userDatas.Imports {
			addEvents, err := getEvents(importData)
			if err != nil {
				if err := send(err.Error(), user); err != nil {
					fmt.Println(err.Error())
				}
			}
			events = append(events, addEvents...)
		}
		for _, task := range userTasks[user] {
			task.Stop()
		}
		userTasks[user] = userTasks[user][0:0]
		for _, event := range events {
			for _, minutes := range userDatas.Minutes {
				tm := time.UnixMilli(event.When - int64(minutes*60*1000))
				if tm.After(time.Now()) {
					d := tm.Sub(time.Now())
					userTasks[user] = append(userTasks[user],
                        time.AfterFunc(d, notifyUzerFunc(user, event)),
					)
				}
			}

		}
	}
}

func notifyUzerFunc(user int64, event imports.Event) func () {
    return func() {
        if err := send(buildMessage(event.Data), user); err != nil {
            fmt.Println(err)
        }
    }
}

func buildMessage(data []imports.DataElement) string {
	result := ""
    for _, v := range data {
        if v.Value != "" {
            result += strings.ReplaceAll(v.Key, "?", "") + ": " + v.Value + "\n"
        }
    }
	return result
}
