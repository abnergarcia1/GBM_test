package data

type IStorageData interface {
	Query(model interface{}, query string, args ...interface{}) (err error)
	Connect() error
	Disconnect()
}
