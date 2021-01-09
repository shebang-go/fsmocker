package file

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/shebang-go/fsmocker/testdouble"
)

var Once sync.Once

// FileInfo represents a file (see os.FileInfo)
type FileInfo struct {
	// FName is the name of the file
	FName string
	// FSize is the size of the file in bytes
	FSize int64
	// FMode is the file mode
	FMode os.FileMode
	// FModTime is the file modification time
	FModTime time.Time
	// FIsDir is true for a directory
	FIsDir bool

	// Error holds a pre-configured error for a file stub.
	Error error
	// Data holds data for a file stub
	Data []byte
	// Path is the full path of the file
	Path string
}

type Configer interface {
	Data(...[]byte) []byte
	Error(...error) error
	Mode(...os.FileMode) os.FileMode
}

type setter struct {
	fi *FileInfo
}

func (s *setter) Data(v ...[]byte) []byte {
	if len(v) == 1 {
		s.fi.Data = v[0]
	}
	return s.fi.Data
}

func (s *setter) Mode(v ...os.FileMode) os.FileMode {
	if len(v) == 1 {
		s.fi.FMode = v[0]
	}
	return s.fi.Mode()
}

func (s *setter) Error(v ...error) error {
	if len(v) == 1 {
		s.fi.Error = v[0]
	}
	return s.fi.Error
}

type Option func(*FS)

func WithFiles(files []*FileInfo) Option {
	return func(fs *FS) {
		fs.AddFiles(files)
	}
}

type FS struct {
	*testdouble.TestDouble
	PathStubs     map[string]*FileInfo
	AbsPathPrefix string
	AbsPathError  error
	t             *testing.T
}

func CreateFS(td *testdouble.TestDouble, opts ...Option) *FS {
	fs := &FS{
		TestDouble:    td,
		PathStubs:     make(map[string]*FileInfo),
		AbsPathPrefix: "",
	}

	for _, opt := range opts {
		opt(fs)
	}
	fs.PathStubs["/"] = &FileInfo{FName: "/", Path: "/", FIsDir: true}
	return fs
}

func (fs *FS) FileInfo(p string) os.FileInfo {
	if v, ok := fs.PathStubs[p]; ok {
		return v
	}
	return nil
}

// Config provides access to stubs
func (fs *FS) Config(p string) Configer {
	var fi *FileInfo
	if v, ok := fs.PathStubs[p]; ok {
		fi = v
	}
	if fi == nil {
		return nil
	}
	s := &setter{fi: fi}
	return s
}

func (fs *FS) AddFiles(in []*FileInfo) {
	for _, v := range in {
		if v.Path != "" && v.Path != "/" {
			fs.PathStubs[v.Path] = v
		}
	}
}

func (fs *FS) getFile(path string, op string) (*FileInfo, error) {
	if v, ok := fs.PathStubs[path]; ok {
		if v.Error != nil {
			fs.TestDouble.Log("return pre-configured error").Path(path).Operation("getFile").Error(v.Error).Done()
			return nil, v.Error

		}
		// fs.TestDouble.Log("return os.FileInfo").Path(path).Operation("getFile").Done()
		return v, nil
	}
	fs.TestDouble.Log("return os.ErrNotExist").Path(path).Operation("getFile").Error(os.ErrNotExist).Done()
	return nil, os.ErrNotExist
}

func (fs *FS) getDirEntries(dirname string) map[string]*FileInfo {

	tmpFiles := make(map[string]*FileInfo)
	for k, v := range fs.PathStubs {
		if strings.HasPrefix(k, dirname) || dirname == "/" {
			f := strings.TrimPrefix(k, dirname)
			g := strings.Split(f, string(os.PathSeparator))
			fname := ""
			if dirname == "/" && len(g) > 0 {

				fname = g[0]
			} else {
				if len(g) == 2 {
					fname = g[1]
				}
			}
			if fname != "" {
				if _, ok := tmpFiles[fname]; !ok {
					tmpFiles[fname] = v
				}
			}
		}
	}
	return tmpFiles
}

func (fs *FS) ReadDir(dirname string) ([]os.FileInfo, error) {

	if err := fs.requireDir(dirname, "ReadDir"); err != nil {
		return nil, err
	}

	tmpFiles := fs.getDirEntries(dirname)
	retval := make([]os.FileInfo, 0)
	for _, v := range tmpFiles {
		if v.Error != nil {
			fs.TestDouble.Log("return pre-configured error").Path(dirname).Operation("ReadDir").Error(v.Error).Done()
			return nil, v.Error
		}
		retval = append(retval, v)

	}
	fs.TestDouble.Log("return return []os.FileInfo").Path(dirname).Operation("ReadDir").Done()
	return retval, nil
}

func (fs *FS) Stat(path string) (os.FileInfo, error) {

	fi, err := fs.getFile(path, "Stat")
	if err != nil {
		return nil, err
	}
	return fi, nil
}

func (fs *FS) ReadFile(path string) ([]byte, error) {

	fi, err := fs.getFile(path, "ReadFile")
	if err != nil {
		return nil, err
	}
	return fi.Data, nil
}

// Walk is a stub for filepath.Walk
func (fs *FS) requireDir(path string, op string) error {
	if v, ok := fs.PathStubs[path]; ok {
		if v.Error != nil {
			fs.TestDouble.Log("return pre-configured error").Path(path).Operation(op).Error(v.Error).Done()
			return v.Error
		}
		if !v.IsDir() {
			fs.TestDouble.Log("return pre-configured error").Path(path).Operation(op).Error(os.ErrInvalid).Done()
			return os.ErrInvalid
		}
	} else {
		fs.TestDouble.Log("return pre-configured error").Path(path).Operation(op).Error(os.ErrNotExist).Done()
		return os.ErrNotExist
	}
	return nil
}

// Walk is a stub for filepath.Walk
func (fs *FS) Walk(root string, walkFn filepath.WalkFunc) error {

	Once.Do(func() {

	})

	if err := fs.requireDir(root, "Walk"); err != nil {
		return err
	}

	keys := []string{}
	for k, _ := range fs.PathStubs {
		keys = append(keys, k)
	}
	// log.Println(">>>>>> Walk keys before sort", keys)
	sort.Strings(keys)
	// log.Println(">>>>>> Walk keys after sort", keys)
	for _, k := range keys {
		if strings.HasPrefix(k, root) {
			fi, err := fs.getFile(k, "walk")
			fs.TestDouble.Log("calling walkFn").Path(k).Operation("Walk").Done()
			walkFn(k, fi, err)
		}
	}
	return nil
}

func (fs *FS) WriteFile(filename string, data []byte, perm os.FileMode) error {

	return nil
}

func (fs *FS) Abs(p string) (string, error) {
	if fs.AbsPathError != nil {
		return "", fs.AbsPathError
	}
	return filepath.Join(fs.AbsPathPrefix, p), nil
}

// NewFile creates a new file. It is used to simplify the interface when only
// names are used.
func NewFile(name string, args ...interface{}) os.FileInfo {
	var size int64
	var mode os.FileMode
	var modTime time.Time

	for _, v := range args {
		switch s := v.(type) {
		case int64:
			size = s
			break
		case os.FileMode:
			mode = s
			break
		case time.Time:
			modTime = s
			break
		}
	}

	return &FileInfo{
		FName:    name,
		FSize:    size,
		FMode:    mode,
		FModTime: modTime,
	}

}

func NewDir(name string, args ...interface{}) os.FileInfo {
	fi := NewFile(name, append(args, true)...)
	fi.(*FileInfo).FIsDir = true
	return fi
}

// Name returns the name of the file
func (fi *FileInfo) Name() string { return fi.FName }

// Size returns the size of the file
func (fi *FileInfo) Size() int64 { return fi.FSize }

// Mode returns the FileMode of the file
func (fi *FileInfo) Mode() os.FileMode { return fi.FMode }

// ModTime returns the modification time of the file
func (fi *FileInfo) ModTime() time.Time { return fi.FModTime }

// Sys is not used but here to satisfy the interface
func (fi *FileInfo) Sys() interface{} { return nil }

// IsDir returns true if the file is a directory.
func (fi *FileInfo) IsDir() bool {
	return fi.FIsDir
}
