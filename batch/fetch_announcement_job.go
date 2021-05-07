package batch

import (
	"log"

	"github.com/goodgoodjm/k-pioneer/kstartup"
)

var api = kstartup.New()

func run() {
	text, err := api.GetAnnouncements()
	if err != nil {
		panic(err)
	}
	log.Println(text)
}
