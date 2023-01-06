package tpack

/**
分割符解析器
*/
import (
	"bufio"
	"github.com/sisobobo/tinx/tiface"
	"github.com/sisobobo/tinx/tlog"
)

type Delimiter struct {
	delim          byte   //分隔符
	maxFrameLen    uint32 //包最大长度
	stripDelimiter bool   //是否除去分隔符
}

func (d *Delimiter) GetMaxFrameLength() uint32 {
	return d.maxFrameLen
}

func (d *Delimiter) UnPack(reader *bufio.Reader) ([]byte, error) {
	line, err := reader.ReadSlice(d.delim)
	if err == bufio.ErrBufferFull {
		tlog.Error("frame length exceeds %d - discarding", d.maxFrameLen)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if d.stripDelimiter {
		line = line[:len(line)-1]
	} else {

	}
	return line, nil
}

func (d *Delimiter) Decode(data []byte) (tiface.IMessage, error) {
	panic("please override func UnPack(data []byte) (tiface.IMessage, error)")
}

func NewDelimiter(delim byte, maxFrameLen uint32, stripDelimiter bool) tiface.IPack {
	return &Delimiter{
		delim:          delim,
		maxFrameLen:    maxFrameLen,
		stripDelimiter: stripDelimiter,
	}
}
