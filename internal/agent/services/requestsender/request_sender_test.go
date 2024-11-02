package requestsender

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/AndrXxX/go-metrics-collector/internal/services/hashgenerator"
)

type closableReadableBodyMock struct {
	mock.Mock
	io.Reader
}

func (m *closableReadableBodyMock) Close() error {
	return nil
}

func (m *closableReadableBodyMock) Read(_ []byte) (n int, err error) {
	return 0, nil
}

type mockClient struct {
	mock.Mock
}

func (m *mockClient) Do(r *http.Request) (*http.Response, error) {
	args := m.Called(r)
	resp, _ := args.Get(0).(*http.Response)
	return resp, args.Error(1)
}

func TestRequestSender_Post(t *testing.T) {
	type fields struct {
		c client
	}
	type args struct {
		url         string
		contentType string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		data    []byte
		wantErr bool
	}{
		{
			name: "Positive test #1",
			fields: fields{
				c: func() *mockClient {
					c := mockClient{}
					c.On("Do", mock.Anything).Return(nil, nil)
					return &c
				}(),
			},
			args:    args{url: "", contentType: ""},
			wantErr: false,
		},
		{
			name: "Positive test #2 with body",
			fields: fields{
				c: func() *mockClient {
					c := mockClient{}
					c.On("Do", mock.Anything).
						Return(&http.Response{Header: http.Header{}, Body: &closableReadableBodyMock{}}, nil)
					return &c
				}(),
			},
			args:    args{url: "", contentType: ""},
			wantErr: false,
		},
		{
			name: "Positive test #3 with data",
			fields: fields{
				c: func() *mockClient {
					c := mockClient{}
					c.On("Do", mock.Anything).Return(nil, nil)
					return &c
				}(),
			},
			data:    []byte("test"),
			args:    args{url: "", contentType: ""},
			wantErr: false,
		},
		{
			name: "Error test #1",
			fields: fields{
				c: func() *mockClient {
					c := mockClient{}
					c.On("Do", mock.Anything).Return(nil, errors.New("error from web server"))
					return &c
				}(),
			},
			args:    args{url: "", contentType: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.fields.c, hashgenerator.Factory().SHA256(), "test")
			err := s.Post(tt.args.url, tt.args.contentType, tt.data)
			assert.Equal(t, tt.wantErr, err != nil, fmt.Errorf("post() error = %v, wantErr %v", err, tt.wantErr))
		})
	}
}

func TestNewRequestSender(t *testing.T) {
	type args struct {
		c client
	}
	tests := []struct {
		name string
		args args
		want *RequestSender
	}{
		{
			name: "Test New RequestSender #1 (Alloc)",
			args: args{c: http.DefaultClient},
			want: &RequestSender{c: http.DefaultClient},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := New(tt.args.c, nil, "")
			assert.Equal(t, tt.want, rs)
		})
	}
}
