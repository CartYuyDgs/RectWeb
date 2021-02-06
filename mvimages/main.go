package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

const Dirname = "D:\\bonc\\sftp\\aaa\\"

func main() {
	files, _ := ioutil.ReadDir(Dirname)
	for _, f := range files {
		fmt.Println(f.Name())
		os.Rename(Dirname+f.Name(), Dirname+"aa"+f.Name())
		time.Sleep(400 * time.Millisecond)
		fmt.Println(f.Name(), path.Base(Dirname))
	}

}
