// Package logger provee un logger global sencillo con niveles.
package logger

import (
	"log"
	"os"
)

var (
	Info  = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile)
	Error = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)
	Debug = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile)
)
