package stack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFromSlice(t *testing.T) {
	type testCase[T any] struct {
		name     string
		elements []T
		want     *Stack[T]
	}
	tests := []testCase[int64]{
		{
			name:     "test1",
			elements: []int64{1, 2, 3},
			want:     &Stack[int64]{elements: []int64{1, 2, 3}},
		},
		{
			name:     "test1",
			elements: []int64{5, 2, 1},
			want:     &Stack[int64]{elements: []int64{5, 2, 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewFromSlice(tt.elements))
		})
	}
}

func TestStackAll(t *testing.T) {
	type testCase[T any] struct {
		name string
		s    *Stack[T]
		want []T
	}
	tests := []testCase[int64]{
		{
			name: "test method All with 1, 2, 3",
			s:    &Stack[int64]{elements: []int64{1, 2, 3}},
			want: []int64{1, 2, 3},
		},
		{
			name: "test method All with 5, 2, 1",
			s:    &Stack[int64]{elements: []int64{5, 2, 1}},
			want: []int64{5, 2, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.s.All())
		})
	}
}

func TestStackCopy(t *testing.T) {
	type testCase[T any] struct {
		name string
		s    *Stack[T]
		want *Stack[T]
	}
	tests := []testCase[int64]{
		{
			name: "test Copy with 1, 2, 3",
			s:    &Stack[int64]{elements: []int64{1, 2, 3}},
			want: &Stack[int64]{elements: []int64{1, 2, 3}},
		},
		{
			name: "test Copy with 5, 11, 44",
			s:    &Stack[int64]{elements: []int64{5, 11, 44}},
			want: &Stack[int64]{elements: []int64{5, 11, 44}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newStack := tt.s.Copy()
			assert.Equal(t, tt.s, newStack)
			assert.Equal(t, tt.s, newStack)
			tt.s.Pop()
			assert.NotEqual(t, tt.s.All(), newStack.All())
		})
	}
}

func TestStackPop(t *testing.T) {
	type want[T any] struct {
		val T
		ok  bool
		s   *Stack[T]
	}
	type testCase[T any] struct {
		name string
		s    *Stack[T]
		want want[T]
	}
	tests := []testCase[int64]{
		{
			name: "test Pop with 5, 11, 44",
			s:    &Stack[int64]{elements: []int64{5, 11, 44}},
			want: want[int64]{val: 44, ok: true, s: &Stack[int64]{elements: []int64{5, 11}}},
		},
		{
			name: "test Pop with 44, -1, 12, -1",
			s:    &Stack[int64]{elements: []int64{44, -1, 12, -1}},
			want: want[int64]{val: -1, ok: true, s: &Stack[int64]{elements: []int64{44, -1, 12}}},
		},
		{
			name: "test Pop with empty slice",
			s:    &Stack[int64]{elements: []int64{}},
			want: want[int64]{val: 0, ok: false, s: &Stack[int64]{elements: []int64{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := tt.s.Pop()
			assert.Equal(t, val, tt.want.val)
			assert.Equal(t, ok, tt.want.ok)
			assert.Equal(t, tt.s, tt.want.s)
		})
	}
}

func TestStackPush(t *testing.T) {
	type testCase[T any] struct {
		name  string
		s     *Stack[T]
		value T
		want  *Stack[T]
	}
	tests := []testCase[int64]{
		{
			name:  "test Push 5 with 5, 11, 44",
			s:     &Stack[int64]{elements: []int64{5, 11, 44}},
			value: 5,
			want:  &Stack[int64]{elements: []int64{5, 11, 44, 5}},
		},
		{
			name:  "test Push 5 with empty slice",
			s:     &Stack[int64]{elements: []int64{}},
			value: 5,
			want:  &Stack[int64]{elements: []int64{5}},
		},
	}
	for _, tt := range tests {
		tt.s.Push(tt.value)
		assert.Equal(t, tt.s, tt.want)
	}
}

func TestStackShift(t *testing.T) {
	type want[T any] struct {
		val T
		ok  bool
		s   *Stack[T]
	}
	type testCase[T any] struct {
		name string
		s    *Stack[T]
		want want[T]
	}
	tests := []testCase[int64]{
		{
			name: "test Shift with 5, 11, 44",
			s:    &Stack[int64]{elements: []int64{5, 11, 44}},
			want: want[int64]{val: 5, ok: true, s: &Stack[int64]{elements: []int64{11, 44}}},
		},
		{
			name: "test Shift with 44, -1, 12, -1",
			s:    &Stack[int64]{elements: []int64{44, -1, 12, -1}},
			want: want[int64]{val: 44, ok: true, s: &Stack[int64]{elements: []int64{-1, 12, -1}}},
		},
		{
			name: "test Shift with empty slice",
			s:    &Stack[int64]{elements: []int64{}},
			want: want[int64]{val: 0, ok: false, s: &Stack[int64]{elements: []int64{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := tt.s.Shift()
			assert.Equal(t, val, tt.want.val)
			assert.Equal(t, ok, tt.want.ok)
			assert.Equal(t, tt.s, tt.want.s)
		})
	}
}
