package storagesaver

import (
	"bufio"
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/memory"
	"github.com/AndrXxX/go-metrics-collector/internal/services/utils"
)

var errorOnFailedAccessRights = false

var checkRights = sync.OnceFunc(func() {
	name := "./check.tmp"
	_, _ = os.OpenFile(name, os.O_CREATE, 0111)
	_, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, permission)
	if err != nil {
		errorOnFailedAccessRights = true
	}
	_ = os.Remove(name)
})

func Test_storageSaver_Restore(t *testing.T) {
	checkRights()
	tests := []struct {
		name          string
		path          string
		beforeRestore func()
		afterRestore  func()
		wantErr       bool
		want          map[string]*models.Metrics
	}{
		{
			name: "Test with error on read file",
			path: "./test_r1.json",
			beforeRestore: func() {
				_, _ = os.OpenFile("./test_r1.json", os.O_CREATE, 0111)
			},
			afterRestore: func() {
				_ = os.Remove("./test_r1.json")
			},
			wantErr: errorOnFailedAccessRights,
			want:    map[string]*models.Metrics{},
		},
		{
			name: "Test with incorrect data in file",
			path: "./test_r2.json",
			beforeRestore: func() {
				f, _ := os.OpenFile("./test_r2.json", os.O_RDWR|os.O_CREATE, 0666)
				bufWriter := bufio.NewWriter(f)
				_, _ = bufWriter.Write([]byte("{{{sdafjs;lfkasjd /n\n"))
				_ = bufWriter.Flush()
			},
			afterRestore: func() {
				_ = os.Remove("./test_r2.json")
			},
			wantErr: false,
			want:    map[string]*models.Metrics{},
		},
		{
			name: "Test with restore data",
			path: "./test_r3.json",
			beforeRestore: func() {
				f, _ := os.OpenFile("./test_r3.json", os.O_RDWR|os.O_CREATE, 0666)
				bufWriter := bufio.NewWriter(f)
				_, _ = bufWriter.Write([]byte("{\"id\":\"test1\",\"type\":\"type1\",\"value\":10.1}"))
				_ = bufWriter.Flush()
			},
			afterRestore: func() {
				_ = os.Remove("./test_r3.json")
			},
			wantErr: false,
			want: map[string]*models.Metrics{
				"test1": {ID: "test1", MType: "type1", Value: utils.Pointer(10.1)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := memory.New[*models.Metrics]()
			ss := New(tt.path, &s, []time.Duration{})
			if tt.beforeRestore != nil {
				tt.beforeRestore()
			}
			assert.Equal(t, tt.wantErr, ss.Restore(context.Background()) != nil)
			if tt.afterRestore != nil {
				tt.afterRestore()
			}
			assert.EqualValues(t, tt.want, s.All(context.Background()))
		})
	}
}

func Test_storageSaver_Save(t *testing.T) {
	checkRights()
	tests := []struct {
		name       string
		path       string
		vals       []*models.Metrics
		beforeSave func()
		afterSave  func()
		wantErr    bool
	}{
		{
			name: "Test with error on save file",
			path: "./test_s1.json",
			beforeSave: func() {
				_, _ = os.OpenFile("./test_s1.json", os.O_CREATE, 0111)
			},
			afterSave: func() {
				_ = os.Remove("./test_s1.json")
			},
			wantErr: errorOnFailedAccessRights,
		},
		{
			name: "Test with succeed save data",
			path: "./test_s2.json",
			vals: []*models.Metrics{
				{ID: "test1", MType: "type1", Value: utils.Pointer(10.1)},
			},
			afterSave: func() {
				f, _ := os.OpenFile("./test_s2.json", os.O_RDONLY, 0666)
				scanner := bufio.NewScanner(f)
				vals := make([]string, 0)
				for scanner.Scan() {
					vals = append(vals, scanner.Text())
				}

				assert.Equal(t, []string{"{\"id\":\"test1\",\"type\":\"type1\",\"value\":10.1}"}, vals)
				_ = os.Remove("./test_s2.json")
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := memory.New[*models.Metrics]()
			ss := New(tt.path, &s, []time.Duration{})
			if tt.vals != nil {
				for _, val := range tt.vals {
					s.Insert(context.Background(), val.ID, val)
				}
			}
			if tt.beforeSave != nil {
				tt.beforeSave()
			}
			assert.Equal(t, tt.wantErr, ss.Save(context.Background()) != nil)
			if tt.afterSave != nil {
				tt.afterSave()
			}
		})
	}
}

func Test_storageSaver_openFile(t *testing.T) {
	checkRights()
	tests := []struct {
		name       string
		path       string
		ri         []time.Duration
		beforeOpen func()
		afterOpen  func()
		wantErr    bool
		wantFile   bool
	}{
		{
			name: "Test with error on open file with repeat",
			path: "./test_o1.json",
			ri:   []time.Duration{50 * time.Millisecond},
			beforeOpen: func() {
				_, _ = os.OpenFile("./test_o1.json", os.O_CREATE, 0111)
			},
			afterOpen: func() {
				_ = os.Remove("./test_o1.json")
			},
			wantErr:  errorOnFailedAccessRights,
			wantFile: !errorOnFailedAccessRights,
		},
		{
			name: "Test with succeed open file",
			path: "./test_o1.json",
			afterOpen: func() {
				_ = os.Remove("./test_o1.json")
			},
			wantErr:  false,
			wantFile: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := New(tt.path, nil, tt.ri)
			if tt.beforeOpen != nil {
				tt.beforeOpen()
			}
			got, err := ss.openFile(tt.path, os.O_RDONLY|os.O_CREATE)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantFile, got != nil)
			if got != nil {
				_ = got.Close()
			}
			if tt.afterOpen != nil {
				tt.afterOpen()
			}
		})
	}
}
