package internal

import (
	"log"
)

type AppSettings struct {
	Logger                *log.Logger
	textListenerUrl       string
	compressedListenerUrl string
	HttpListenerUrl       string
}

