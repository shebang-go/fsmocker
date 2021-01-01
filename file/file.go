package file

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/shebang-go/fsmocker/testdouble"
)

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

type FS struct {
	*testdouble.TestDouble
	PathStubs     map[string]*FileInfo
	AbsPathPrefix string
	AbsPathError  error
	t             *testing.T
}

func CreateFS(td *testdouble.TestDouble, opts ...testdouble.Option) *FS {
	t := &FS{
		TestDouble:    td,
		PathStubs:     make(map[string]*FileInfo),
		AbsPathPrefix: "",
	}

	for _, opt := range opts {
		opt(&t.OptionData)
	}
	t.PathStubs["/"] = &FileInfo{FName: "/", Path: "/", FIsDir: true}
	return t
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
			fs.TestDouble.Log("*FS.getFile(): op=%s path=%s -> return pre-configured error %v", op, path, v.Error)
			return nil, v.Error

		}
		fs.TestDouble.Log("*FS.getFile(): op=%s path=%s -> return os.FileInfo", op, path)
		return v, nil
	}
	fs.TestDouble.Log("*FS.getFile(): op=%s path=%s -> return os.ErrNotExist", op, path)
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
	if v, ok := fs.PathStubs[dirname]; ok {
		if v.Error != nil {
			fs.TestDouble.Log("op=ReadDir path=%s -> return pre-configured error %v", dirname, v.Error)
			return nil, v.Error
		}
		if !v.IsDir() {
			fs.TestDouble.Log("op=ReadDir path=%s -> return os.ErrInvalid NOT A DIRECTORY", dirname)
			return nil, os.ErrInvalid
		}
	} else {
		fs.TestDouble.Log("op=ReadDir path=%s -> return os.ErrNotExist", dirname)
		return nil, os.ErrNotExist
	}
	tmpFiles := fs.getDirEntries(dirname)
	retval := make([]os.FileInfo, 0)
	for _, v := range tmpFiles {
		if v.Error != nil {
			fs.TestDouble.Log("op=ReadDir path=%s -> return pre-configured error %v", filepath.Join(dirname, v.Name()), v.Error)
			return nil, v.Error
		}
		retval = append(retval, v)

	}
	fs.TestDouble.Log("op=ReadDir path=%s -> return []os.FileInfo", dirname)
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
func (fs *FS) Walk(root string, walkFn filepath.WalkFunc) error {

	if v, ok := fs.PathStubs[root]; ok {
		if v.Error != nil {
			fs.TestDouble.Log("op=Walk path=%s -> return pre-configured error %v", root, v.Error)
			return v.Error
		}
		if !v.IsDir() {
			fs.TestDouble.Log("op=Walk path=%s -> return os.ErrInvalid NOT A DIRECTORY", root)
			return os.ErrInvalid
		}
	} else {
		fs.TestDouble.Log("op=Walk path=%s -> return os.ErrNotExist", root)
		return os.ErrNotExist
	}

	keys := []string{}
	for k, _ := range fs.PathStubs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if strings.HasPrefix(k, root) {
			fi, err := fs.getFile(k, "walk")
			fs.TestDouble.Log("op=Walk path=%s -> calling walkFn", k)
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
