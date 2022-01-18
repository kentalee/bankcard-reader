package command

const ( // 卡片头部程序, 响应中包含后续程序Id
	// PSEOfChipCard 芯片卡 0x315041592E5359532E4444463031
	PSEOfChipCard = "1PAY.SYS.DDF01"
	// PSEOfRFIDCard 非接卡 0x325041592E5359532E4444463031
	PSEOfRFIDCard = "2PAY.SYS.DDF01"
)

func SelectApplication(name []byte) *Command {
	return &Command{
		cla: 0x00, ins: InsSelectFile, p1: 0x04, p2: 0x00,
		cmd: name,
	}
}

func SelectPSEOfChipCard() *Command {
	return SelectApplication([]byte(PSEOfChipCard))
}
func SelectPSEOfRFIDCard() *Command {
	return SelectApplication([]byte(PSEOfRFIDCard))
}

func ReadMore(length int) *Command {
	return &Command{
		cla: 0x00, ins: 0xC0, p1: 0x00, p2: 0x00, ne: length,
	}
}
