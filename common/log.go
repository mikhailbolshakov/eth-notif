package common

import "log"

func LogErr(err error) {
	log.Printf("[error] %s\n", err.Error())
}

func LogDbg(format string, v ...any) {
	log.Printf(format+"\n", v...)
}
