package stub

import (
	"errors"
	"testing"

	"github.com/shebang-go/fsmocker/file"
	"github.com/shebang-go/fsmocker/testdouble"
	"github.com/stretchr/testify/assert"
)

func TestNewStub(t *testing.T) {
	type args struct {
		paths []string
		opts  []Option
	}
	tests := []struct {
		name string
		args args
		want *Stub
	}{
		{
			name: "noError",
			args: args{
				paths: []string{
					"/folder1[file1]",
					"/folder2",
					"/folder2[file2]",
				},
			},
			want: &Stub{
				testDouble: *testdouble.NewTestDouble().(*testdouble.TestDouble),
				fs: &file.FS{
					TestDouble: testdouble.NewTestDouble().(*testdouble.TestDouble),
					PathStubs: map[string]*file.FileInfo{
						"/":              {FName: "/", Path: "/", FIsDir: true},
						"/folder1":       {FName: "folder1", Path: "/folder1", FIsDir: true},
						"/folder1/file1": {FName: "file1", Path: "/folder1/file1"},
						"/folder2":       {FName: "folder2", Path: "/folder2", FIsDir: true},
						"/folder2/file2": {FName: "file2", Path: "/folder2/file2"},
					},
					AbsPathPrefix: "",
				},
			},
		},
		{
			name: "withGlobalOption",
			args: args{
				paths: []string{},
				opts: []Option{
					WithGlobalOptions(testdouble.WithError(errors.New("test"))),
				},
			},
			want: &Stub{
				testDouble: *testdouble.NewTestDouble(testdouble.WithError(errors.New("test"))).(*testdouble.TestDouble),
				fs: &file.FS{
					TestDouble: testdouble.NewTestDouble(testdouble.WithError(errors.New("test"))).(*testdouble.TestDouble),
					PathStubs: map[string]*file.FileInfo{
						"/": {FName: "/", Path: "/", FIsDir: true},
					},
					AbsPathPrefix: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewStub(tt.args.paths, tt.args.opts...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, got.(*Stub).testDouble, got.(*Stub).testDouble)
		})
	}
}

// func TestStub_ConfigRaw(t *testing.T) {
// 	type args struct {
// 		p string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want *file.FileInfo
// 	}{
// 		{
// 			name: "path=/folder1",
// 			args: args{
// 				p: "/folder1",
// 			},
// 			want: &file.FileInfo{
// 				FName:  "folder1",
// 				FIsDir: true,
// 				Path:   "/folder1",
// 			},
// 		},
// 		{
// 			name: "path=/folder1/file1",
// 			args: args{
// 				p: "/folder1/file1",
// 			},
// 			want: &file.FileInfo{
// 				FName: "file1",
// 				Path:  "/folder1/file1",
// 			},
// 		},
// 		{
// 			name: "path=nil",
// 			args: args{
// 				p: "/invalid",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			st := NewStub([]string{"/folder1[file1]"})
// 			got := st.ConfigRaw(tt.args.p)
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

func TestStub_Options(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "withGlobalOptions",
			args: args{
				opts: []Option{
					WithGlobalOptions(testdouble.WithError(errors.New("test"))),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st1 := NewStub([]string{"/folder1[file1]"})
			st2 := NewStub([]string{"/folder1[file1]"})
			st2.Options(tt.args.opts...)
			assert.NotEqual(t, st1, st2)
		})
	}
}
