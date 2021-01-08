// build !integration
// build example

package fsmocker_test

// import (
// 	"fmt"
//
// 	"github.com/shebang-go/fsmocker"
// )
//
// func ExampleFsmocker() {
//
// 	stub := fsmocker.NewStub([]string{"/home/john[file1(data=test)]"})
//
// 	data, _ := stub.ReadFile("/home/john/file1")
// 	fmt.Printf("data:%s\n", data)
//
// 	fi, _ := stub.Stat("/home/john/file1")
// 	fmt.Printf("name:%s isdir:%t\n", fi.Name(), fi.IsDir())
//
// 	fi, _ = stub.Stat("/home/john")
// 	fmt.Printf("name:%s isdir:%t\n", fi.Name(), fi.IsDir())
//
// 	stubError := fsmocker.NewStub([]string{"/home/baddir(err=baderror)"})
// 	_, err := stubError.Stat("/home/baddir")
// 	fmt.Printf("err:%s\n", err)
//
// 	// Output:
// 	// data:test
// 	// name:file1 isdir:false
// 	// name:john isdir:true
// 	// err:baderror
// }
