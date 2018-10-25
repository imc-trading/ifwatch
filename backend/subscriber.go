package backend

type MessageHandler func(string, []byte)

type Subscriber interface {
	Backend() Backend
	Start() error
	Stop() error
	AddHandler(MessageHandler)
	Versions(string, MessageHandler) error
}
