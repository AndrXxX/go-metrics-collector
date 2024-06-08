package utils

import (
	"errors"
	"github.com/AndrXxX/go-metrics-collector/internal/mocks"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestRequestSender_Post(t *testing.T) {
	type fields struct {
		ub URLBuilder
		c  Client
	}
	type args struct {
		params      URLParams
		contentType string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Positive test #1",
			fields: fields{
				ub: NewMetricURLBuilder("host"),
				c: &mocks.MockClient{
					PostFunc: func(url, contentType string, body io.Reader) (*http.Response, error) {
						return nil, nil
					},
				},
			},
			args:    args{params: map[string]any{}, contentType: ""},
			wantErr: false,
		},
		{
			name: "Positive test #2 with body",
			fields: fields{
				ub: NewMetricURLBuilder("host"),
				c: &mocks.MockClient{
					PostFunc: func(url, contentType string, body io.Reader) (*http.Response, error) {
						return &http.Response{Header: http.Header{}, Body: io.NopCloser(strings.NewReader("test"))}, nil
					},
				},
			},
			args:    args{params: map[string]any{}, contentType: ""},
			wantErr: false,
		},
		{
			name: "Error test #1",
			fields: fields{
				ub: NewMetricURLBuilder("host"),
				c: &mocks.MockClient{
					PostFunc: func(url, contentType string, body io.Reader) (*http.Response, error) {
						return nil, errors.New("error from web server")
					},
				},
			},
			args:    args{params: map[string]any{}, contentType: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RequestSender{
				ub: tt.fields.ub,
				c:  tt.fields.c,
			}
			if err := s.Post(tt.args.params, tt.args.contentType); (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewRequestSender(t *testing.T) {
	type args struct {
		ub URLBuilder
		c  Client
	}
	tests := []struct {
		name string
		args args
		want *RequestSender
	}{
		{
			name: "Test New RequestSender #1 (Alloc)",
			args: args{ub: NewMetricURLBuilder(""), c: http.DefaultClient},
			want: &RequestSender{ub: NewMetricURLBuilder(""), c: http.DefaultClient},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRequestSender(tt.args.ub, tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRequestSender() = %v, want %v", got, tt.want)
			}
		})
	}
}
