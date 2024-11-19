package metricsvaluesetter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
)

func TestFactory(t *testing.T) {
	t.Run("Test for get Factory", func(t *testing.T) {
		f := &factory{
			setters: map[string]setter{
				metrics.Counter: &counterValueSetter{},
				metrics.Gauge:   &gaugeValueSetter{},
			},
		}
		assert.Equal(t, f, Factory())
	})
}

func Test_factory_CounterValueSetter(t *testing.T) {
	t.Run("Test for get counterValueSetter", func(t *testing.T) {
		f := &factory{}
		assert.Equal(t, &counterValueSetter{}, f.CounterValueSetter())
	})
}

func Test_factory_GaugeValueSetter(t *testing.T) {
	t.Run("Test for get gaugeValueSetter", func(t *testing.T) {
		f := &factory{}
		assert.Equal(t, &gaugeValueSetter{}, f.GaugeValueSetter())
	})
}

func Test_factory_SetterByType(t *testing.T) {
	tests := []struct {
		name      string
		setters   map[string]setter
		mType     string
		want      setter
		wantExist bool
	}{
		{
			name: "Test for gaugeValueSetter with key gauge",
			setters: map[string]setter{
				"gauge": &gaugeValueSetter{},
			},
			mType:     "gauge",
			want:      &gaugeValueSetter{},
			wantExist: true,
		},
		{
			name:      "Test for not exist setter",
			setters:   map[string]setter{},
			mType:     "unknown",
			want:      nil,
			wantExist: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &factory{setters: tt.setters}
			s, ok := f.SetterByType(tt.mType)
			assert.Equal(t, tt.want, s)
			assert.Equal(t, tt.wantExist, ok)
		})
	}
}
