// json编码，待实现
package codec

import (
	"bufio"
	"encoding/json"
	"io"
)

type JsonCodec struct {
	conn io.ReadWriteCloser
	buf bufio.Writer
	dec *json.Decoder
	enc *json.Encoder
}

func (j *JsonCodec) Close() error {
	return j.conn.Close()
}

func (j *JsonCodec) ReadHeader(header *Header) error  {
	return j.dec.Decode(header)
}


