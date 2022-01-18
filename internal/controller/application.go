package controller

import (
	"errors"

	"github.com/dumacp/smartcard/pcsc"

	"cardreader/common/command"
	"cardreader/common/status"
	"cardreader/common/structure"
)

func ParseCardApplication(card pcsc.Card) (app *structure.Application, err error) {
	var response *command.Response
	if response, err = command.SelectPSEOfChipCard().Exec(card.Apdu); err == nil {
		return structure.ParseApplication(response.Content())
	} else if errors.Is(err, status.Ex6A82FileNotFound) {
		if response, err = command.SelectPSEOfRFIDCard().Exec(card.Apdu); err == nil {
			return structure.ParseApplication(response.Content())
		}
	}
	return nil, err
}
