package tiface

import "github.com/sisobobo/tinx/tpkg/bufio"

type IPack interface {
	Connect(channel IChannel)                                    //连接回调
	Read(channel IChannel, reader *bufio.Reader) ([]byte, error) //读取包
	UnPack(channel IChannel, arr []byte) IMessage                //解析包
	Pack(channel IChannel, msg IMessage) []byte                  //封装包
	DisConnect(channel IChannel)                                 //断开连接回调
}
