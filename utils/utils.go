package utils

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func CheckErr(err error, s ...string) {
	if err != nil {
		if len(s) == 0 {
			log.Fatal(err)
		} else {
			for _, message := range s {
				log.Fatal(message)
			}
			log.Fatal(err)
		}
	}
}

func StopOnErr(err error, s ...string) {
	if err != nil {
		if len(s) == 0 {
			log.Fatal(err)
		} else {
			for _, message := range s {
				log.Fatal(message)
			}
			log.Fatal(err)
		}
		os.Exit(-1)
	}
}
