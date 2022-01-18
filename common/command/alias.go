package command

const (
	InsEraseBinary   = 0x0E
	InsVerify        = 0x20
	InsManageChannel = 0x70
	InsExternalAuth  = 0x82
	InsGetChallenge  = 0x84
	InsInternalAuth  = 0x88
	InsSelectFile    = 0xA4
	InsReadBinary    = 0xB0
	InsReadRecord    = 0xB2
	InsGetResponse   = 0xC0
	InsEnvelope      = 0xC2
	InsGetData       = 0xCA
	InsWriteBinary   = 0xD0
	InsWriteRecord   = 0xD2
	InsUpdateBinary  = 0xD6
	InsPutData       = 0xDA
	InsUpdateData    = 0xDC
	InsAppendRecord  = 0xE2
)
