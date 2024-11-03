package metricsidentifier

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

func TestNewJSONIdentifier(t *testing.T) {
	tests := []struct {
		name string
		want *jsonMetricsIdentifier
	}{
		{
			name: "Test OK NewJSONIdentifier",
			want: &jsonMetricsIdentifier{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewJSONIdentifier())
		})
	}
}

func TestJSONMetricsIdentifierProcess(t *testing.T) {
	type modelValue struct {
		Value float64
		Delta int64
	}
	tests := []struct {
		name       string
		body       string
		modelValue *modelValue
		want       *models.Metrics
		wantErr    bool
	}{
		{
			name:       "OK Counter test",
			body:       "{\"id\":\"test\",\"type\":\"counter\",\"delta\":10,\"value\":null}",
			modelValue: &modelValue{Delta: 10},
			want:       &models.Metrics{ID: "test", MType: metrics.Counter},
			wantErr:    false,
		},
		{
			name:       "OK Gauge test",
			body:       "{\"id\":\"test\",\"type\":\"gauge\",\"delta\":null,\"value\":10.1}",
			modelValue: &modelValue{Value: 10.1},
			want:       &models.Metrics{ID: "test", MType: metrics.Gauge},
			wantErr:    false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i := &jsonMetricsIdentifier{}
			request := httptest.NewRequest("", "/test", strings.NewReader(test.body))
			m, err := i.Process(request)
			require.Equal(t, test.wantErr, err != nil)
			if !test.wantErr {
				assert.Equal(t, test.want.ID, m.ID)
				assert.Equal(t, test.want.MType, m.MType)
				if test.modelValue != nil && test.modelValue.Delta != 0 {
					assert.Equal(t, test.modelValue.Delta, *m.Delta)
				}
				if test.modelValue != nil && test.modelValue.Value != 0 {
					assert.InDelta(t, test.modelValue.Value, *m.Value, 0.0001)
				}
			}
		})
	}
}

func Example_jsonMetricsIdentifier_Process() {
	body := "{\"id\":\"test\",\"type\":\"counter\",\"delta\":10,\"value\":null}"
	request := httptest.NewRequest("", "/test", strings.NewReader(body))
	i := &jsonMetricsIdentifier{}

	m, err := i.Process(request)
	if err != nil {
		return
	}

	fmt.Println(m.ID)
	fmt.Println(m.MType)
	fmt.Println(*m.Delta)
	fmt.Println(m.Value)

	// Output:
	// test
	// counter
	// 10
	// <nil>
}
