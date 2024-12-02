package requestsender

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender/dto"
)

func TestWithGzip(t *testing.T) {
	tests := []struct {
		name        string
		comp        dataCompressor
		params      dto.ParamsDto
		wantHeaders map[string]string
		wantData    string
		wantErr     bool
	}{
		{
			name: "Test with error on read data",
			params: func() dto.ParamsDto {
				r := readerMock{}
				r.On("Read", mock.Anything).Return(0, errors.New("error on read"))
				return dto.ParamsDto{Buf: &r, Headers: make(map[string]string)}
			}(),
			wantHeaders: map[string]string{},
			wantErr:     true,
		},
		{
			name: "Test with error on compress data",
			comp: func() *dataCompressorMock {
				c := dataCompressorMock{}
				c.On("Compress", mock.Anything).Return(nil, errors.New("error on compress"))
				return &c
			}(),
			params:      dto.ParamsDto{Buf: bytes.NewReader([]byte("test")), Headers: make(map[string]string)},
			wantHeaders: map[string]string{},
			wantErr:     true,
		},
		{
			name: "Test with succeed compress",
			comp: func() *dataCompressorMock {
				c := dataCompressorMock{}
				c.On("Compress", mock.Anything).Return(bytes.NewReader([]byte("compressedData")), nil)
				return &c
			}(),
			params: dto.ParamsDto{Buf: bytes.NewReader([]byte("test")), Headers: make(map[string]string)},
			wantHeaders: map[string]string{
				"Content-Encoding": "gzip",
				"Accept-Encoding":  "gzip",
			},
			wantData: "compressedData",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := WithGzip(tt.comp)
			require.Equal(t, tt.wantErr, f(&tt.params) != nil)
			assert.Equal(t, tt.wantHeaders, tt.params.Headers)
			if tt.wantData != "" {
				data, err := io.ReadAll(tt.params.Buf)
				require.NoError(t, err)
				assert.Equal(t, tt.wantData, string(data))
			}
		})
	}
}
