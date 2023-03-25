package main

import (
	"github.com/nnnewb/nf/cmd"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	lumberjackLogger := lumberjack.Logger{
		Filename:   "nf.log",
		MaxSize:    100,
		MaxAge:     30,
		MaxBackups: 30,
		LocalTime:  true,
		Compress:   true,
	}
	log.SetOutput(io.MultiWriter(os.Stdout, &lumberjackLogger))
}

func main() {
	cmd.Execute()
}
