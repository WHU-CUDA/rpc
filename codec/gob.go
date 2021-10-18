package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

// the instance of GobCodec
type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

var _ Codec = (*GobCodec)(nil)

func (g *GobCodec) Close() error {
	return g.conn.Close()
}

// read and decode the header of the message
func (g *GobCodec) ReadHeader(header *Header) error {
	return g.dec.Decode(header)
}

// read and decode the body of message
func (g *GobCodec) ReadBody(body interface{}) error {
	return g.dec.Decode(body)
}

// write the information into the message
func (g *GobCodec) Write(header *Header, body interface{}) (err error) {
	defer func() {
		_ = g.buf.Flush()
		if err != nil {
			_ = g.Close()
		}
	}()

	if err := g.enc.Encode(header); err != nil {
		log.Println("rpc codec: gob error encoding header; ", err)
		return err
	}
	if err := g.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body; ", err)
		return err
	}
	return nil
}

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}
