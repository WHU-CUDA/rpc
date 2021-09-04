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
