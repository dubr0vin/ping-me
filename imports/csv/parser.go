package data_getter

import (
	"encoding/csv"
	"errors"
    "io"
    "net/http"
	"pingMe/imports"
	"strconv"
	"time"
)

func NewImport(params []string) (*Import, error) {
    if len(params) == 0 {
        return nil, errors.New("no enough params")
    }
	return &Import{
        url: params[0],
	}, nil
}

type Import struct {
	url string
}

func (s *Import) GetEvents() ([]imports.Event, error) {
	resp, err := http.Get(s.url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Wrong status code" + strconv.Itoa(resp.StatusCode))
	}
	csvReader := csv.NewReader(resp.Body);
	defer func(Body io.ReadCloser) {
        err := Body.Close();
        if err != nil {
            panic(err)
        }
    }(resp.Body)
	records, err := csvReader.ReadAll();
	if err != nil {
		return nil, err
	}
	names := records[0]
	result := make([]imports.Event, len(records)-1)
	for i, record := range records {
		if i < 1 {
			continue
		}
		data := make([]imports.DataElement, 0)
		timeStr := ""
		for j := 1; j < len(names) && names[j] != ""; j++ {
			if names[j] == "Дата" || names[j] == "Day" {
				timeStr += record[j] + "~"
			} else if names[j] == "Время" || names[j] == "Time" {
				timeStr += record[j] + "~"
			} else {
                data = append(data, imports.DataElement{
                    Key: names[j],
                    Value: record[j],
                })
			}
		}
		parsed, err := time.Parse("1/2/2006~03:04:05 PM~MST", timeStr+"MSK")
		if err != nil {
			return nil, err
		}
		result = append(result, imports.Event{
			When: parsed.UnixMilli(),
			Data: data,
		})
	}
	return result, nil
}
