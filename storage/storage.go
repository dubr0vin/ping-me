package storage

import (
    "encoding/json"
    "os"
)

const (
    FileName = "storage.json"
)


type ImportData struct {
    Id string
    Params []string
}
type UserData struct {
    Imports []ImportData
    Minutes []uint32
}
type Data struct {
    Users map[int64]UserData
}

func (d *Data) Save() error {
    file, err := os.Create(FileName)
    if err != nil {
        return err
    }
    defer file.Close()
    decoder := json.NewEncoder(file)
    err = decoder.Encode(d)
    if err != nil {
        return err
    }
    return nil
}

func Load() (*Data, error) {
    file, err := os.Open(FileName)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    var v Data
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&v)
    if err != nil {
        return nil, err
    }
    return &v, nil
}