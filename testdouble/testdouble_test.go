package testdouble

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithLogging(t *testing.T) {
	type args struct {
		t *testing.T
	}
	tests := []struct {
		name       string
		args       args
		want       *OptionData
		optionData *OptionData
	}{
		{
			name: "noError",
			args: args{t: &testing.T{}},
			want: &OptionData{t: &testing.T{}, enableLogging: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithLogging(tt.args.t)
			optionData := &OptionData{}
			got(optionData)
			assert.EqualValues(t, optionData, tt.want)
		})
	}
}
