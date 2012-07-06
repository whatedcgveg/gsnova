package event

import (
	"bytes"
	"errors"
	"snappy"
	"strconv"
)

type CompressEvent struct {
	CompressType uint32
	Ev           Event
}

func (ev *CompressEvent) Encode(buffer *bytes.Buffer) {
	if ev.CompressType != COMPRESSOR_NONE && ev.CompressType != COMPRESSOR_SNAPPY {
		ev.CompressType = COMPRESSOR_NONE
	}
	EncodeUInt64Value(buffer, uint64(ev.CompressType))
	//ev.ev.Encode(buffer);
	var buf bytes.Buffer
	EncodeRawEvent(&buf, ev.Ev)
	switch ev.CompressType {
	case COMPRESSOR_NONE:
		buffer.Write(buf.Bytes())
	case COMPRESSOR_SNAPPY:
		evbuf := make([]byte, 0)
		newbuf, err := snappy.Encode(evbuf, buf.Bytes())
		buffer.Write(newbuf)
	}
}
func (ev *CompressEvent) Decode(buffer *bytes.Buffer) (err error) {
	ev.CompressType, err = DecodeUInt32Value(buffer)
	if err != nil {
		return
	}
	switch ev.CompressType {
	case COMPRESSOR_NONE:
		//success, ev.Ev, _ = ParseEvent(buffer)
		//return success
	case COMPRESSOR_SNAPPY:
		b := make([]byte, 0, 0)
		b, err = snappy.Decode(b, buffer.Bytes())
		if err != nil {
			return
		}
		//success, ev.Ev, _ = ParseEvent(bytes.NewBuffer(newbuf))
		//return success
	default:
		return errors.New("Not supported compress type:" + strconv.Itoa(int(ev.CompressType)))
	}
	return nil
}
