package codec

import "io"

type Header struct {
	ServiceMethod string // 服务名和方法名，通常与 Go 语言中的结构体和方法相映射。
	Seq uint64 // 请求的序号，也可以认为是某个请求的 ID，用来区分不同的请求
	Error string // 错误信息，客户端置为空，服务端如果如果发生错误，将错误信息置于 Error 中
}

// Codec 消息体编解码接口
// for the different codec implements
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}


type NewCodecFunc func(io.ReadWriteCloser) Codec

// Type 通过类型可以获得不同得消息编解码得构造方法
type Type string

const (
	GobType Type = "application/gob"
	JsonType Type = "application/json"
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
