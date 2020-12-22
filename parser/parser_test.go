package parser

import (
	"errors"
	"testing"

	"github.com/shebang-go/fsmocker/file"
	"github.com/stretchr/testify/assert"
)

// func TestCreateParser(t *testing.T) {
// 	type args struct {
// 		v string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want *Parser
// 	}{
// 		{
// 			name: "singleFileWithoutSep",
// 			args: args{v: "file"},
// 			want: &Parser{elements: []string{"file"}},
// 		},
// 		{
// 			name: "singleFileWithSep",
// 			args: args{v: "/file"},
// 			want: &Parser{elements: []string{"file"}},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := CreateParser(tt.args.v); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("CreateParser() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestParser_Parse(t *testing.T) {
	type fields struct {
		elements []string
	}
	type args struct {
		v string
	}
	tests := []struct {
		name  string
		input string
		args  args
		want  []*file.FileInfo
	}{
		{
			name:  "simpleDir",
			input: "dir",
			want: []*file.FileInfo{
				{FName: "dir", Path: "/dir", FIsDir: true},
			},
		},
		{
			name:  "dirWithTags",
			input: "dir(err=test)",
			want: []*file.FileInfo{
				{FName: "dir", FIsDir: true, Path: "/dir", Error: errors.New("test")},
			},
		},
		{
			name:  "dirWithTwoTags",
			input: "dir(err=test, data=test)",
			want: []*file.FileInfo{
				{FName: "dir", FIsDir: true, Path: "/dir", Error: errors.New("test"), Data: []byte("test")},
			},
		},
		{
			name:  "dirWithInvalidTag",
			input: "dir(err=test, invalid=test)",
			want: []*file.FileInfo{
				{FName: "dir", FIsDir: true, Path: "/dir", Error: errors.New("test")},
			},
		},
		{
			name:  "dirWithFiles",
			input: "dir[file1, file2]",
			want: []*file.FileInfo{
				{FName: "dir", FIsDir: true, Path: "/dir"},
				{FName: "file1", Path: "/dir/file1"},
				{FName: "file2", Path: "/dir/file2"},
			},
		},
		{
			name:  "dirWithFilesAndTags",
			input: "dir[file1(err=someerr), file2(data=testdata)]",
			want: []*file.FileInfo{
				{FName: "dir", FIsDir: true, Path: "/dir"},
				{FName: "file1", Path: "/dir/file1", Error: errors.New("someerr")},
				{FName: "file2", Path: "/dir/file2", Data: []byte("testdata")},
			},
		},
		{
			name:  "complexPath1",
			input: "/home/barbara/dir[file1(err=someerr), file2(data=testdata)]",
			want: []*file.FileInfo{
				{FName: "home", FIsDir: true, Path: "/home"},
				{FName: "barbara", FIsDir: true, Path: "/home/barbara"},
				{FName: "dir", FIsDir: true, Path: "/home/barbara/dir"},
				{FName: "file1", Path: "/home/barbara/dir/file1", Error: errors.New("someerr")},
				{FName: "file2", Path: "/home/barbara/dir/file2", Data: []byte("testdata")},
			},
		},
		{
			name:  "complexPath2",
			input: "/home/barbara/dir[file1(err=someerr), file2(data=testdata)]/subdir",
			want: []*file.FileInfo{
				{FName: "home", FIsDir: true, Path: "/home"},
				{FName: "barbara", FIsDir: true, Path: "/home/barbara"},
				{FName: "dir", FIsDir: true, Path: "/home/barbara/dir"},
				{FName: "subdir", FIsDir: true, Path: "/home/barbara/dir/subdir"},
				{FName: "file1", Path: "/home/barbara/dir/file1", Error: errors.New("someerr")},
				{FName: "file2", Path: "/home/barbara/dir/file2", Data: []byte("testdata")},
			},
		},
		{
			name:  "complexPath3",
			input: "/home/barbara[notes.txt(data=somenote)]/dir[file1(err=someerr), file2(data=testdata)]/subdir",
			want: []*file.FileInfo{
				{FName: "home", FIsDir: true, Path: "/home"},
				{FName: "barbara", FIsDir: true, Path: "/home/barbara"},
				{FName: "notes.txt", Path: "/home/barbara/notes.txt", Data: []byte("somenote")},
				{FName: "dir", FIsDir: true, Path: "/home/barbara/dir"},
				{FName: "subdir", FIsDir: true, Path: "/home/barbara/dir/subdir"},
				{FName: "file1", Path: "/home/barbara/dir/file1", Error: errors.New("someerr")},
				{FName: "file2", Path: "/home/barbara/dir/file2", Data: []byte("testdata")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.input)
			assert.ElementsMatch(t, tt.want, got)
			// assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestParser_parseFilename(t *testing.T) {
	type fields struct {
		elements []string
	}
	type args struct {
		v string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "1",
			args: args{v: "file"},
			want: "file",
		},
		{
			name: "2",
			args: args{v: "dir[file1(err=someerr) file2()]"},
			want: "dir",
		},
		{
			name: "3",
			args: args{v: "file1(data=testdata)"},
			want: "file1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseFilename(tt.args.v); got != tt.want {
				t.Errorf("Parser.parseFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}
