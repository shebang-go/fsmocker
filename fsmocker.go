// Package fsmocker provides test doubles for file systems methods which have
// side effects (ex: Stat(p string) (os.FileInfo, error))
package fsmocker

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/shebang-go/fsmocker/file"
	"github.com/shebang-go/fsmocker/stub"
	"github.com/shebang-go/fsmocker/testdouble"
)

type Stub interface {
	Config(p string) file.Configer
	Options(opts ...stub.Option)
	Stat(path string) (os.FileInfo, error)
	ReadFile(path string) ([]byte, error)
	ReadDir(path string) ([]os.FileInfo, error)
	Walk(root string, walkFn filepath.WalkFunc) error
	WriteFile(filename string, data []byte, perm os.FileMode) error
	Abs(p string) (string, error)
}
type TestDoubleOption func(td *testdouble.TestDouble)
type StubOption stub.Option

// type StubOption func(st *stub.Stub)

type stubOption struct{}

func WithLogging(t *testing.T) TestDoubleOption {
	return func(td *testdouble.TestDouble) {
		td.EnableLogging(t)
	}
}
func WithGlobalOptions(opts ...TestDoubleOption) StubOption {
	return func(st *stub.Stub) {
		for _, opt := range opts {
			// td := st.(*stub.Stub)
			td := st.TestDouble()
			log.Println(">>>", td, opt)
			opt(td.(*testdouble.TestDouble))
		}
	}
}

func NewStub(paths []string, opts ...StubOption) *stub.Stub {
	o := []stub.Option{}
	for _, v := range opts {
		o = append(o, stub.Option(v))
		// o = append(o, stub.Stub(stub.Option(v)))
	}
	st := stub.NewStub(paths, o...)
	return st.(*stub.Stub)
}
