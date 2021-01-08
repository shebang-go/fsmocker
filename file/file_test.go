package file

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/shebang-go/fsmocker/testdouble"
	"github.com/stretchr/testify/assert"
)

func TestFS_ReadDir(t *testing.T) {
	type fields struct {
		PathStubs map[string]*FileInfo
	}
	type args struct {
		dirname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []os.FileInfo
		wantErr bool
	}{
		{
			name: "noError",
			fields: fields{PathStubs: map[string]*FileInfo{
				"/home/maggy":              {FName: "maggy", FIsDir: true},
				"/home/maggy/file1":        {FName: "file1"},
				"/home/maggy/file2":        {FName: "file2"},
				"/home/maggy/subdir":       {FName: "subdir", FIsDir: true},
				"/home/maggy/subdir/file3": {FName: "file3"},
			}},
			args: args{dirname: "/home/maggy"},
			want: []os.FileInfo{
				&FileInfo{FName: "file1"},
				&FileInfo{FName: "file2"},
				&FileInfo{FName: "subdir", FIsDir: true},
			},
		},
		{
			name: "errorReadDir",
			fields: fields{PathStubs: map[string]*FileInfo{
				"/home/maggy":              {FName: "maggy", FIsDir: true, Error: errors.New("errorReadDir")},
				"/home/maggy/file1":        {FName: "file1"},
				"/home/maggy/file2":        {FName: "file2"},
				"/home/maggy/subdir":       {FName: "subdir", FIsDir: true},
				"/home/maggy/subdir/file3": {FName: "file3"},
			}},
			args:    args{dirname: "/home/maggy"},
			wantErr: true,
		},
		{
			name: "errorReadFile",
			fields: fields{PathStubs: map[string]*FileInfo{
				"/home/maggy":       {FName: "maggy", FIsDir: true},
				"/home/maggy/file1": {FName: "file1", Error: errors.New("errorFile")},
			}},
			args:    args{dirname: "/home/maggy"},
			wantErr: true,
		},
		{
			name:    "errorNotExist",
			fields:  fields{PathStubs: map[string]*FileInfo{}},
			args:    args{dirname: "/home/maggy"},
			wantErr: true,
		},
		{
			name: "errorInvalidDir",
			fields: fields{PathStubs: map[string]*FileInfo{

				"/home/maggy": {FName: "maggy"},
			}},
			args:    args{dirname: "/home/maggy"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FS{
				PathStubs:  tt.fields.PathStubs,
				TestDouble: testdouble.NewTestDouble(testdouble.WithLogging(t)).(*testdouble.TestDouble),
			}
			got, err := fs.ReadDir(tt.args.dirname)
			assert.ElementsMatch(t, got, tt.want)
			if (err != nil) != tt.wantErr {
				t.Errorf("FS.ReadDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFS_getDirEntries(t *testing.T) {
	type fields struct {
		PathStubs map[string]*FileInfo
	}
	type args struct {
		dirname string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]*FileInfo
	}{
		{
			name: "noError",
			fields: fields{PathStubs: map[string]*FileInfo{
				"/home/maggy/file1":        {FName: "file1"},
				"/home/maggy/file2":        {FName: "file2"},
				"/home/maggy/subdir":       {FName: "subdir", FIsDir: true},
				"/home/maggy/subdir/file3": {FName: "file3"},
			}},
			args: args{dirname: "/home/maggy"},
			want: map[string]*FileInfo{
				"file1":  {FName: "file1"},
				"file2":  {FName: "file2"},
				"subdir": {FName: "subdir", FIsDir: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FS{
				PathStubs:  tt.fields.PathStubs,
				TestDouble: testdouble.NewTestDouble(testdouble.WithLogging(t)).(*testdouble.TestDouble),
			}
			if got := fs.getDirEntries(tt.args.dirname); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FS.getDirEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFS_Stat(t *testing.T) {
	type fields struct {
		PathStubs map[string]*FileInfo
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    os.FileInfo
		wantErr bool
	}{
		{
			name: "noErrorDir",
			fields: fields{PathStubs: map[string]*FileInfo{
				"/home/maggy":       {FName: "maggy", FIsDir: true},
				"/home/maggy/file1": {FName: "file1"},
			}},
			args: args{path: "/home/maggy"},
			want: &FileInfo{FName: "maggy", FIsDir: true},
		},
		{
			name: "noErrorFile",
			fields: fields{PathStubs: map[string]*FileInfo{
				"/home/maggy/file1": {FName: "file1"},
			}},
			args: args{path: "/home/maggy/file1"},
			want: &FileInfo{FName: "file1"},
		},
		{
			name: "error",
			fields: fields{PathStubs: map[string]*FileInfo{
				"/home/maggy/file1": {FName: "file1", Error: errors.New("errorStat")},
			}},
			args:    args{path: "/home/maggy/file1"},
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FS{
				PathStubs:  tt.fields.PathStubs,
				TestDouble: testdouble.NewTestDouble(testdouble.WithLogging(t)).(*testdouble.TestDouble),
			}
			got, err := fs.Stat(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("FS.Stat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FS.Stat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFS_ReadFile(t *testing.T) {
	type fields struct {
		PathStubs map[string]*FileInfo
		t         *testing.T
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "noError",
			fields: fields{PathStubs: map[string]*FileInfo{
				"/home/maggy/file1": {FName: "file1", Data: []byte("test")},
			}},
			args: args{path: "/home/maggy/file1"},
			want: []byte("test"),
		},
		{
			name: "error",
			fields: fields{PathStubs: map[string]*FileInfo{
				"/home/maggy/file1": {FName: "file1", Error: errors.New("errorReadFile")},
			}},
			args:    args{path: "/home/maggy/file1"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FS{
				PathStubs:  tt.fields.PathStubs,
				TestDouble: testdouble.NewTestDouble(testdouble.WithLogging(t)).(*testdouble.TestDouble),
				t:          tt.fields.t,
			}
			got, err := fs.ReadFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("FS.ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FS.ReadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFS_Walk(t *testing.T) {
	type fields struct {
		TestDouble    testdouble.TestDouble
		PathStubs     map[string]*FileInfo
		AbsPathPrefix string
		AbsPathError  error
		t             *testing.T
	}
	type args struct {
		root   string
		walkFn filepath.WalkFunc
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr error
	}{
		{
			name: "noError(/home/maggy)",
			fields: fields{PathStubs: map[string]*FileInfo{
				"/home/maggy":              {FName: "maggy", FIsDir: true},
				"/home/maggy/file1":        {FName: "file1"},
				"/home/maggy/file2":        {FName: "file2"},
				"/home/maggy/subdir":       {FName: "subdir", FIsDir: true},
				"/home/maggy/subdir/file3": {FName: "file3"},
			}},
			args: args{root: "/home/maggy"},
			want: []string{"/home/maggy", "/home/maggy/file1", "/home/maggy/file2", "/home/maggy/subdir", "/home/maggy/subdir/file3"},
		},
		{
			name: "noError(/home/maggy/subdir)",
			fields: fields{PathStubs: map[string]*FileInfo{
				"/home/maggy":              {FName: "maggy", FIsDir: true},
				"/home/maggy/file1":        {FName: "file1"},
				"/home/maggy/file2":        {FName: "file2"},
				"/home/maggy/subdir":       {FName: "subdir", FIsDir: true},
				"/home/maggy/subdir/file3": {FName: "file3"},
			}},
			args: args{root: "/home/maggy/subdir"},
			want: []string{"/home/maggy/subdir", "/home/maggy/subdir/file3"},
		},
		{
			name:    "errorInvalidPath",
			fields:  fields{PathStubs: map[string]*FileInfo{}},
			args:    args{root: "/invalid"},
			wantErr: errors.New("file does not exist"),
		},
		{
			name: "errorNotADirectory",

			fields: fields{PathStubs: map[string]*FileInfo{
				"/home/maggy":              {FName: "maggy", FIsDir: true},
				"/home/maggy/file1":        {FName: "file1"},
				"/home/maggy/file2":        {FName: "file2"},
				"/home/maggy/subdir":       {FName: "subdir", FIsDir: true},
				"/home/maggy/subdir/file3": {FName: "file3"},
			}},
			args:    args{root: "/home/maggy/file1"},
			wantErr: errors.New("invalid argument"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FS{
				PathStubs:  tt.fields.PathStubs,
				TestDouble: testdouble.NewTestDouble(testdouble.WithLogging(t)).(*testdouble.TestDouble),
				t:          tt.fields.t,
			}
			got := []string{}
			err := fs.Walk(tt.args.root, func(path string, f os.FileInfo, err error) error {
				got = append(got, path)
				return err
			})
			if tt.wantErr == nil {
				assert.Exactly(t, tt.want, got)
			} else {
				assert.EqualError(t, tt.wantErr, err.Error())
			}
		})
	}
}

func TestFS_Config(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name  string
		files []*FileInfo
		args  args
		want  *Configer
	}{
		{
			name: "path=/folder1",
			args: args{
				p: "/folder1",
			},
			files: []*FileInfo{
				{FName: "/", FIsDir: true, Path: "/"},
				{FName: "folder1", FIsDir: true, Path: "/folder1"},
				{FName: "file1", Path: "/folder1/file1"},
			},
		},
		{
			name: "path=/folder1/file1",
			args: args{
				p: "/folder1/file1",
			},
			files: []*FileInfo{
				{FName: "/", FIsDir: true, Path: "/"},
				{FName: "folder1", FIsDir: true, Path: "/folder1"},
				{FName: "file1", Path: "/folder1/file1"},
			},
		},
		{
			name: "path=nil",
			args: args{
				p: "/invalid",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := CreateFS(testdouble.NewTestDouble(testdouble.WithLogging(t)).(*testdouble.TestDouble), WithFiles(tt.files))

			configer := st.Config(tt.args.p)
			if configer != nil {
				currentData := st.Config(tt.args.p).Data()
				assert.Equal(t, []byte(nil), currentData)

				newData := st.Config(tt.args.p).Data([]byte("test"))
				assert.NotEqual(t, currentData, newData)
				assert.Equal(t, []byte("test"), st.Config(tt.args.p).Data())

				currentMode := st.Config(tt.args.p).Mode()
				assert.Equal(t, os.FileMode(0x0), currentMode)

				newMode := st.Config(tt.args.p).Mode(0x750)
				assert.Equal(t, os.FileMode(0x750), newMode)

				currentErr := st.Config(tt.args.p).Error()
				assert.Equal(t, nil, currentErr)

				newError := st.Config(tt.args.p).Error(errors.New("test"))
				assert.Equal(t, errors.New("test"), newError)
			} else {
				assert.Equal(t, nil, configer)
			}
		})
	}
}
