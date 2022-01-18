package structure

import (
	"cardreader/common/tlv"
)

type Application struct {
	raw         []byte
	FCITemplate *struct {
		DFName      []byte `tlv:"84" json:"DFName,omitempty"` // Dedicated File Name
		Proprietary *struct {
			AppLabel         string `tlv:"50" json:"appLabel,omitempty"`           // Application Label
			AppPriority      byte   `tlv:"87" json:"appPriority,omitempty"`        // Application Priority Indicator
			SFI              []byte `tlv:"88" json:"SFI,omitempty"`                // Short File Identifier
			LangPref         string `tlv:"5F2D" json:"langPref,omitempty"`         // Language Preference
			DOPL             []byte `tlv:"9F38" json:"DOPL,omitempty"`             // Processing Options Data Object List (PDOL)
			AppPrefName      string `tlv:"9F12" json:"AppPrefName,omitempty"`      //  Application Preferred Name
			IssuerCodeTabIdx byte   `tlv:"9F11" json:"issuerCodeTabIdx,omitempty"` //  Issuer Code Table Index
			IssuerData       *struct {
				LogEntry []byte `tlv:"9F4D" json:"LogEntry,omitempty"` //  Log Entry
			} `tlv:"BF0C" json:"issuerData,omitempty"` // Issuer Discretionary Data
		} `tlv:"A5" json:"proprietary,omitempty"`
	} `tlv:"6F" json:"FCITemplate,omitempty"` // FCI File Control Information Template
}

func ParseApplication(raw []byte) (*Application, error) {
	var app = Application{raw: raw}
	_, err := tlv.Unmarshal(raw, &app)
	return &app, err
}
