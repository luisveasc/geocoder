package util

import (
	"log"
	"strconv"

	"github.com/natefinch/lumberjack"
)

func LoadLogFile(filename string, amountMB string, numBackups string) {

	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	_amountMB, _ := strconv.Atoi(amountMB)
	_numBackups, _ := strconv.Atoi(numBackups)

	log.SetOutput(&lumberjack.Logger{

		Filename:   "logs/" + filename,
		MaxSize:    _amountMB,   // megabytes after which new file is created
		MaxBackups: _numBackups, // number of backups
		MaxAge:     28,          //days
	})

}
