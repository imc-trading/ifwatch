package backend

type Publisher interface {
	Backend() Backend
	Send(string, interface{}) error
	Close() error
}
