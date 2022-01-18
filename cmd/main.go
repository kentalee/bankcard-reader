package main

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/dumacp/smartcard/pcsc"

	"cardreader/common/command"
	"cardreader/common/structure"
	"cardreader/internal/controller"
)

const contactless = false

func main() {
	var err error
	var card pcsc.Card
	{ /* connect to card */
		var ctx *pcsc.Context
		log.Printf("============= Connect Card ==============")
		if ctx, err = pcsc.NewContext(); err != nil {
			panic(err)
		}
		if card, err = controller.SelectCardAuto(ctx, contactless); err != nil {
			panic(err)
		}
		defer func() {
			_ = card.DisconnectCard()
		}()
	}
	{ /* reset to card */
		log.Printf("============== Reset Card ===============")
		var respOfAtr []byte // Answer To Reset
		if respOfAtr, err = card.ATR(); err != nil {
			panic(err)
		}
		log.Printf("ATR: %02x\n", respOfAtr)
	}
	var application *structure.Application
	{ /* select card application */
		log.Printf("============== Select App ===============")
		if application, err = controller.ParseCardApplication(card); err != nil {
			panic(err)
		}
		PrintObj(application)
	}
	{ /* select card application */
		var response *command.Response
		log.Printf("============== Select App ===============")
		if response, err = command.SelectApplication([]byte{
			0xA0, 0x00, 0x00, 0x03, 0x33, 0x01, 0x01, 0x02,
		}).Exec(card.Apdu); err != nil {
			panic(err)
		} else if application, err = structure.ParseApplication(response.Content()); err != nil {
			panic(err)
		}
		PrintObj(application)
	}
	// command.SelectApplication().Exec(card.Apdu)
}

func PrintObj(obj interface{}) {
	var j bytes.Buffer
	je := json.NewEncoder(&j)
	je.SetIndent("", "  ")
	_ = je.Encode(obj)
	log.Println(j.String())
}
