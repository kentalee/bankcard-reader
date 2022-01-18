package controller

import (
	"errors"
	"strings"

	"github.com/dumacp/smartcard/pcsc"
)

// todo dialog select

func SelectCardAuto(pcscCtx *pcsc.Context, contactless bool) (card pcsc.Card, err error) {
	var readers []string
	if readers, err = pcsc.ListReaders(pcscCtx); err != nil {
		return nil, err
	}
	var keyWord string
	if contactless {
		keyWord = "Contactless Reader"
	} else {
		keyWord = "Contact Reader"
	}
	var reader string
	for i, r := range readers {
		if strings.Contains(r, keyWord) {
			reader = readers[i]
		}
	}
	if reader == "" {
		return nil, errors.New("reader not found")
	}
	return SelectProtocolAuto(pcscCtx, reader)
}

func SelectProtocolAuto(pcscCtx *pcsc.Context, reader string) (card pcsc.Card, err error) {
	if card, err = pcsc.NewReader(pcscCtx, reader).ConnectCardPCSC_T0(); err != nil {
		card, err = pcsc.NewReader(pcscCtx, reader).ConnectCardPCSC()
		if err != nil {
			return nil, err
		}
	}
	return card, nil
}
