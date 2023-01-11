package imports


type DataElement struct {
    Key, Value string
}
type Event struct {
    When int64
    Data []DataElement
}

type Import interface {
    GetEvents() ([]Event, error)
}