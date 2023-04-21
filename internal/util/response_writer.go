package util

import (
	"log"
	"net/http"
)

func WriteBytesToRW(w http.ResponseWriter, b []byte) {
	if _, err := w.Write(b); err != nil {
		log.Println(err.Error())
	}
}

func WriteStringToRW(w http.ResponseWriter, s string) {
	WriteBytesToRW(w, []byte(s))
}
