package util

import (
	"log"

	"github.com/natefinch/lumberjack"
)

func LoadLogFile(filepath string, filename string, amountMB int, numBackups int) {

	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	log.SetOutput(&lumberjack.Logger{
		Filename:   "logs/" + filename,
		MaxSize:    amountMB,   // megabytes after which new file is created
		MaxBackups: numBackups, // number of backups
		MaxAge:     28,         //days
	})

}
