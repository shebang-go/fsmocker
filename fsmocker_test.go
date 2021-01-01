// Package fsmocker provides test doubles for file systems methods which have
// side effects (ex: Stat(p string) (os.file.FileInfo, error))

package fsmocker

import (
	"os"
	"testing"

	"github.com/shebang-go/fsmocker/file"
	"github.com/shebang-go/fsmocker/testdouble"
	"github.com/stretchr/testify/assert"
)

// func TestStub_Stat(t *testing.T) {
// 	type fields struct {
// 		options *testdouble.OptionData
// 	}
// 	type args struct {
// 		path string
// 	}
// 	tests := []struct {
// 		name      string
// 		fields    fields
// 		testPaths []string
// 		args      args
// 		want      os.FileInfo
// 		wantErr   bool
// 	}{
// 		{
// 			name: "noError",
// 			args: args{path: "/home/barbara/notes.txt"},
// 			testPaths: []string{
// 				"/home/barbara[notes.txt(data=somenote)]/dir[file1(err=someerr), file2(data=testdata)]/subdir",
// 			},
// 			want: &file.FileInfo{FName: "notes.txt", Path: "/home/barbara/notes.txt", Data: []byte("somenote")},
// 		},
// 		{
// 			name: "errorStub",
// 			args: args{path: "/home/barbara/dir/file1"},
// 			testPaths: []string{
// 				"/home/barbara[notes.txt(data=somenote)]/dir[file1(err=someerr), file2(data=testdata)]/subdir",
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			stub := NewStub(tt.testPaths)
// 			got, err := stub.Stat(tt.args.path)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Stub.Stat() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Stub.Stat() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func createTestFiles(t *testing.T, paths []string) {
//
// }
//
func TestStub_ReadDir(t *testing.T) {
	type fields struct {
		options *testdouble.OptionData
		fs      *file.FS
	}
	type args struct {
		path string
	}
	tests := []struct {
		name      string
		fields    fields
		testPaths []string
		args      args
		want      []os.FileInfo
		wantErr   bool
	}{
		{
			name: "noError",
			args: args{path: "/home/barbara"},
			testPaths: []string{
				"/home/barbara[notes.txt(data=somenote)]/dir[file1, file2(data=testdata)]/subdir",
			},
			want: []os.FileInfo{
				&file.FileInfo{FName: "notes.txt", Path: "/home/barbara/notes.txt", Data: []byte("somenote")},
				&file.FileInfo{FName: "dir", Path: "/home/barbara/dir", FIsDir: true},
			},
		},
		{
			name: "rootDir",
			args: args{path: "/"},
			testPaths: []string{
				"/",
				"/dir",
			},
			want: []os.FileInfo{
				&file.FileInfo{FName: "dir", Path: "/dir", FIsDir: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := NewStub(tt.testPaths)
			got, err := stub.ReadDir(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stub.ReadDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.ElementsMatch(t, got, tt.want)
		})
	}
}

func TestNewStub(t *testing.T) {
	type args struct {
		paths []string
		opts  []testdouble.Option
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
				testDouble: *testdouble.NewTestDouble(),
				FS: &file.FS{
					TestDouble: testdouble.NewTestDouble(),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewStub(tt.args.paths, tt.args.opts...)
			assert.Equal(t, tt.want, got)
		})
	}
}
