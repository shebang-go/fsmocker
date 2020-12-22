# fsmocker

`fsmocker` is a little tool library which provides test doubles for file system
related methods with side effects. Currently, only "stubbing" is supported (pre
configured behaviour).

## Example

```go
stub := fsmocker.NewStub([]string{"/home/john[file1(data=test)]"})

data, _ := stub.ReadFile("/home/john/file1")
fmt.Printf("data: %s\n", data)

fi, _ := stub.Stat("/home/john/file1")
fmt.Printf("name:%s isdir:%t\n", fi.Name(), fi.IsDir())

fi, _ = stub.Stat("/home/john")
fmt.Printf("name:%s isdir:%t\n", fi.Name(), fi.IsDir())

stubError := fsmocker.NewStub([]string{"/home/baddir(err=baderror)"})
_, err := stubError.Stat("/home/baddir")
fmt.Printf("err:%s\n", err)

// Output:
// data: testa
// name:file1 isdir:false
// name:john isdir:true
// err:baderror
```

## Supported Methods

-   Stat
-   ReadDir
-   ReadFile

## Syntax

Stubs are created using path expressions:

This is a directory

```
/somedir
/home/user
```

    By default all elements of a path are considered directories.

This is a file

```
/somedir/file(isdir=false)
```

    Tags inside `()` control the behaviour of a file. `isdir=false` creates
    a file.

This is a directory with files

```
/somedir[filemock1.txt,filemock2.txt]
```

    Files inside `[]` are direct child nodes of its parent directory. Files can
    have tags, too.

This is file error

```
/somedir/filemock.txt(isdir=false, err=baderror)
```

    Tags `isdir` and `err` are used to create a file with an error condition
    when accessed.

This is a directory with a file error (different approach)

```
/somedir[filemock.txt(err=baderror)]
```
