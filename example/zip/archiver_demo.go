package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mholt/archiver/v4"
)

// 压缩
func main() {
	filenames := map[string]string{
		//"F:\\go-toolbox\\README.md": "README.md",
		"F:\\go-toolbox\\": "",
	}
	files, err := archiver.FilesFromDisk(nil, filenames)
	if err != nil {
		panic(fmt.Sprintf("[1]%s", err))
	}

	out, err := os.Create("F:\\zip-test\\example.zip")
	if err != nil {
		panic(fmt.Sprintf("[2]%s", err))
	}
	defer out.Close()

	format := archiver.CompressedArchive{
		Archival: archiver.Zip{},
	}
	err = format.Archive(context.Background(), out, files)
	if err != nil {
		panic(fmt.Sprintf("[3]%s", err))
	}

	fmt.Println("Done!")
}
