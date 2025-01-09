package main

import (
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

func main() {
	c := exec.Command("bash")
	f, err := pty.Start(c)
	if err != nil {
		panic(err)
	}

	go func() {
		f.Write([]byte("ls\n"))
		time.Sleep(1000 * time.Millisecond)
		f.Write([]byte("whoami\n"))
		time.Sleep(1000 * time.Millisecond)
		f.Write([]byte("sleep 3\n"))
		f.Write([]byte{4})
		time.Sleep(1000 * time.Millisecond)
	}()
	io.Copy(os.Stdout, f)
}
