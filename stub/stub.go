// Package fsmocker provides test doubles for file systems methods which have
// side effects (ex: Stat(p string) (os.FileInfo, error))
package stub

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shebang-go/fsmocker/file"
	"github.com/shebang-go/fsmocker/parser"
	"github.com/shebang-go/fsmocker/testdouble"
)

type Option func(*Stub)

type Stuber interface {
	TestDouble(v ...testdouble.TestDoubler) testdouble.TestDoubler
	Options(opts ...Option)
	Config(p string) file.Configer
	Stat(path string) (os.FileInfo, error)
	ReadDir(path string) ([]os.FileInfo, error)
	Walk(root string, walkFn filepath.WalkFunc) error
	WriteFile(filename string, data []byte, perm os.FileMode) error
	Abs(p string) (string, error)
}

// Stub represents a file system stub.
type Stub struct {
	testDouble testdouble.TestDouble
	// fs is the test file system.
	fs *file.FS
}

var (
//WithLogging is an option func to enable logging output in tests.
// WithLogging = testdouble.WithLogging
)

func WithGlobalOptions(opts ...testdouble.Option) Option {
	return func(stub *Stub) {
		for _, opt := range opts {
			opt(&stub.testDouble)
		}
	}
}

// NewStub creates a new stub.
func NewStub(paths []string, opts ...Option) Stuber {

	stub := &Stub{
		testDouble: testdouble.TestDouble{},
	}
	stub.fs = file.CreateFS(&stub.testDouble)
	for _, v := range paths {
		stub.fs.AddFiles(parser.Parse(v))
	}

	for _, opt := range opts {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>", opt)
		opt(stub)
	}
	return stub
}

// ConfigRaw provides access to stubs
// func (st *Stub) ConfigRaw(p string) *file.FileInfo {
// 	if v, ok := st.fs.PathStubs[p]; ok {
// 		return v
// 	}
// 	return nil
// }
// Config provides access to stubs
func (st *Stub) TestDouble(v ...testdouble.TestDoubler) testdouble.TestDoubler {
	return &st.testDouble
}

// Config provides access to stubs
func (st *Stub) Config(p string) file.Configer {
	return st.fs.Config(p)
}

func (st *Stub) Options(opts ...Option) {
	for _, opt := range opts {
		opt(st)
	}
}

// Stat is a stub for os.Stat
func (st *Stub) Stat(path string) (os.FileInfo, error) {
	return st.fs.Stat(path)
}

// ReadFile is a stub for ioutil.ReadFile
func (st *Stub) ReadFile(path string) ([]byte, error) {
	return st.fs.ReadFile(path)
}

// ReadDir is a stub for ioutil.ReadDir
func (st *Stub) ReadDir(path string) ([]os.FileInfo, error) {
	return st.fs.ReadDir(path)
}

// Walk is a stub for filepath.Walk
func (st *Stub) Walk(root string, walkFn filepath.WalkFunc) error {
	return st.fs.Walk(root, walkFn)
}

// WriteFile is a stub for ioutil.WriteFile
func (st *Stub) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return st.fs.WriteFile(filename, data, perm)
}

// Abs is a stub for filepath.Abs
func (st *Stub) Abs(p string) (string, error) {
	return st.fs.Abs(p)
}
