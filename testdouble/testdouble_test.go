package testdouble

// func TestWithLogging(t *testing.T) {
// 	type args struct {
// 		t *testing.T
// 	}
// 	tests := []struct {
// 		name       string
// 		args       args
// 		want       *TestDouble
// 		optionData *TestDouble
// 	}{
// 		{
// 			name: "noError",
// 			args: args{t: &testing.T{}},
// 			want: &TestDouble{OptionData: OptionData{t: &testing.T{}, enableLogging: true}},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := WithLogging(tt.args.t)
// 			td := &TestDouble{}
// 			got(td)
// 			assert.EqualValues(t, td, tt.want)
// 		})
// 	}
// }
