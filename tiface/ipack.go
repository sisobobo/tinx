package tiface

import (
	"bufio"
)

type IPack interface {
	GetMaxFrameLength() uint32                   //获取包的最大长度
	UnPack(reader *bufio.Reader) ([]byte, error) //粘包/分包处理
	Decode(dat []byte) (IMessage, error)         //解析成IMessage
}
