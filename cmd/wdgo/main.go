package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

var (
	flagLogFile = flag.Bool("logFile", false, "write log into log.txt")
)

func main() {
	buf := bufio.NewWriter(os.Stdout)
	defer buf.Flush()
	log.SetOutput(buf)
	if *flagLogFile {
		logfile, _ := os.OpenFile("log.txt", os.O_CREATE, 0777)
		defer logfile.Close()
		log.SetOutput(logfile)
	}
	a := newApp()
	if err := a.run(); err != nil {
		panic(err)
	}
}
