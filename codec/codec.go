package codec

import "io"

type Header struct {
	ServiceMethod string // ex: ServiceName.MethodName
	Seq           uint64 // the sequence number of the request, to classfiy different requests.
	Error         string // Error message, client is empty, if server has error,put Error.ErrorMessage into this.
}

// Codec encode and decode interface
// for the different codec implements
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type NewCodecFunc func(io.ReadWriteCloser) Codec

// Type across different Type get different constructor
type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
