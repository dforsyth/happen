package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

type ctx struct {
	lastCheck time.Time
	files     []string
	cmd       string
}

func runCommand(ctx *ctx) {
	log.Printf("running command")
	b, e := exec.Command("/bin/sh", "-c", ctx.cmd).CombinedOutput()
	if e != nil {
		log.Println(e)
	}

	log.Printf("command output:\n%s", b)
}

func loop(ctx *ctx) {
	for ctx.lastCheck = time.Now(); ; {
		time.Sleep(time.Second)
		for _, f := range ctx.files {
			i, e := os.Stat(f)
			if e != nil {
				log.Println(e)
				continue
			}

			if i.ModTime().After(ctx.lastCheck) {
				ctx.lastCheck = time.Now()
				runCommand(ctx)
				break
			}
		}
	}
}

func usage() {
	log.Printf("usage: %s \"<command>\" <files>", os.Args[0])
	os.Exit(0)
}

func main() {
	if len(os.Args) < 3 {
		usage()
	}

	ctx := &ctx{
		files: os.Args[2:],
		cmd:   os.Args[1],
	}

	log.Printf("command to run is \"%s\"", ctx.cmd)
	log.Printf("watched files are %s", ctx.files)

	loop(ctx)
}
