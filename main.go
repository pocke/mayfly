package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	var n int
	flag.IntVar(&n, "t", 600, "time as second (default 600)")
	flag.Parse()

	err := Do(flag.Args(), n)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Do(cmd []string, n int) error {
	if len(cmd) == 0 {
		return fmt.Errorf("command expected.")
	}

	for {
		c := exec.Command(cmd[0], cmd[1:]...)
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		err := c.Start()
		if err != nil {
			return err
		}
		fmt.Printf("[mayfly] Start command > %s\n", strings.Join(cmd, " "))

		DieOrTimeout(c, n)
	}
}

func DieOrTimeout(c *exec.Cmd, n int) {
	die := make(chan struct{})
	go func() {
		c.Wait()
		die <- struct{}{}
	}()
	timeout := time.After(time.Duration(n) * time.Second)

	select {
	case <-die:
		return
	case <-timeout:
		c.Process.Kill()
		return
	}
}
