package dto

import (
	"io"
)

type ParamsDto struct {
	Buf     io.Reader
	Data    []byte
	Headers map[string]string
}
