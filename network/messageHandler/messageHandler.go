package messageHandler

import (
	"log"
	"sync"

	"capnproto.org/go/capnp/v3"
)

type MessageHandler struct {
	FuncChan   chan func()
	handlers   []func(*capnp.Message, string)
	handlerMux sync.Mutex
}

func NewMessageHandler() *MessageHandler {
	handler := &MessageHandler{
		FuncChan: make(chan func(), 100),
	}
	go handler.run()
	return handler
}

func (handler *MessageHandler) AddHandler(fun func(*capnp.Message, string)) {
	handler.handlerMux.Lock()
	defer handler.handlerMux.Unlock()
	handler.handlers = append(handler.handlers, fun)
}

func (handler *MessageHandler) run() {
	for f := range handler.FuncChan {
		f()
	}
}

func (handler *MessageHandler) HandleMessage(msg *capnp.Message, source string) {
	handler.handlerMux.Lock()
	defer handler.handlerMux.Unlock()
	if len(handler.handlers) == 1 {
		handler.handlers[0](msg, source)
	} else if len(handler.handlers) > 0 {
		nextHandler := handler.handlers[0]
		handler.handlers = handler.handlers[1:]
		nextHandler(msg, source)
	} else {
		log.Println("No more handlers available")
	}
}
