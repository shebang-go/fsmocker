// Package fsmocker provides test doubles for file systems methods which have
// side effects (ex: Stat(p string) (os.FileInfo, error))
package fsmocker

import (
	"os"

	"github.com/shebang-go/fsmocker/file"
	"github.com/shebang-go/fsmocker/parser"
	"github.com/shebang-go/fsmocker/testdouble"
)

// Stub represents a file system stub.
type Stub struct {
	testDouble testdouble.TestDouble
	// FS is the test file system.
	FS *file.FS
}

var (
	//WithLogging is an option func to enable logging output in tests.
	WithLogging = testdouble.WithLogging
)

// NewStub creates a new stub.
func NewStub(paths []string, opts ...testdouble.Option) *Stub {

	stub := &Stub{
		testDouble: *testdouble.NewTestDouble(opts...),
		FS:         file.CreateFS(opts...),
	}

	for _, v := range paths {
		stub.FS.AddFiles(parser.Parse(v))
	}
	return stub
}

// Stat is a stub for os.Stat
func (st *Stub) Stat(path string) (os.FileInfo, error) {
	return st.FS.Stat(path)
}

// ReadFile is a stub for ioutil.ReadFile
func (st *Stub) ReadFile(path string) ([]byte, error) {
	return st.FS.ReadFile(path)
}

// ReadDir is a stub for ioutil.ReadDir
func (st *Stub) ReadDir(path string) ([]os.FileInfo, error) {
	return st.FS.ReadDir(path)
}

// Abs is a stub for filepath.Abs
func (st *Stub) Abs(p string) (string, error) {
	return st.FS.Abs(p)
}
