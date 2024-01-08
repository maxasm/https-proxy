package certs

import (
	"log"
	"os"
	"fmt"
)

type WarningLogger struct {}

func (wl *WarningLogger) Write(buffer []byte) (int, error) {
	return fmt.Fprintf(os.Stdout, "\033[91m%s\033[0m", string(buffer))
}

type DebugLogger struct {}

func (wl *DebugLogger) Write(buffer []byte) (int, error) {
	return fmt.Fprintf(os.Stdout, "\033[92m%s\033[0m", string(buffer))
}

var dl *log.Logger = log.New(&DebugLogger{}, "[DEBUG] :", log.Lshortfile)
var wl *log.Logger = log.New(&WarningLogger{}, "[WARNING] :", log.Lshortfile)
