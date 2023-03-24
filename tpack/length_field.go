package tpack

import (
	"encoding/binary"
	"errors"
	"github.com/sisobobo/tinx/tiface"
	"github.com/sisobobo/tinx/tlog"
	"github.com/sisobobo/tinx/tpkg/bufio"
)

type LengthFieldDecoder struct {
	BaseConnect
	order               binary.ByteOrder
	lengthFieldOffset   int
	lengthFieldLength   int
	lengthAdjustment    int
	initialBytesToStrip bool
}

func (c *LengthFieldDecoder) Read(channel tiface.IChannel, reader *bufio.Reader) ([]byte, error) {
	var (
		offset     []byte
		adjustment []byte
		lenData    []byte
		arr        []byte
		result     []byte
		l          int
		err        error
	)
	if c.lengthFieldOffset > 0 {
		offset, err = reader.Pop(c.lengthFieldOffset)
		if err != nil {
			tlog.Warnf("read error : %v", err)
			return nil, err
		}
		if !c.initialBytesToStrip {
			result = append(result, offset...)
		}
	}
	lenData, err = reader.Pop(c.lengthFieldLength)
	if err != nil {
		tlog.Warnf("read error : %v", err)
		return nil, err
	}
	if !c.initialBytesToStrip {
		result = append(result, lenData...)
	}
	if c.lengthFieldLength == 2 {
		l = int(c.order.Uint16(lenData))
	} else if c.lengthFieldLength == 4 {
		l = int(c.order.Uint32(lenData))
	} else if c.lengthFieldLength == 8 {
		l = int(c.order.Uint64(lenData))
	} else {
		return nil, errors.New("lengthFieldLength error")
	}
	if c.lengthAdjustment > 0 {
		adjustment, err = reader.Pop(c.lengthAdjustment)
		if err != nil {
			tlog.Warnf("read error : %v", err)
			return nil, err
		}
		if !c.initialBytesToStrip {
			result = append(result, adjustment...)
		}
	}
	arr, err = reader.Pop(l)
	if err != nil {
		tlog.Warnf("read error : %v", err)
		return nil, err
	}
	result = append(result, arr...)
	return result, nil
}

func NewLengthFieldDecoder(offset, len, adjustmentLen int, strip bool, order binary.ByteOrder) LengthFieldDecoder {
	return LengthFieldDecoder{
		order:               order,
		lengthFieldOffset:   offset,
		lengthFieldLength:   len,
		lengthAdjustment:    adjustmentLen,
		initialBytesToStrip: strip,
	}
}
