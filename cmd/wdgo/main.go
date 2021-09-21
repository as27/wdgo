package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	flagLogFile   = flag.String("logFile", "", "write log into this file")
	flagAppFile   = flag.String("file", "wdgo.txt", "define the file where the boardlist is stored")
	flagEventPath = flag.String("events", "events", "path where the events of each board are stored")
)

func main() {
	flag.Parse()
	buf := bufio.NewWriter(os.Stdout)
	defer buf.Flush()
	log.SetOutput(buf)
	if *flagLogFile != "" {
		logfile, _ := os.OpenFile(*flagLogFile, os.O_CREATE, 0777)
		defer logfile.Close()
		log.SetOutput(logfile)
	}
	lockFile := fmt.Sprintf("%s.lock", *flagAppFile)
	// check lock
	_, err := os.Stat(lockFile)
	if !os.IsNotExist(err) {
		fmt.Println("there is another instance running")
		os.Exit(2)
	}
	// set lock
	lock, err := os.Create(lockFile)
	if err != nil {
		log.Println(err)
		log.Println("can not creat lock-file")
		os.Exit(1)
	}
	lock.Close()
	// remove lock after main is finished
	defer os.Remove(lockFile)

	a := newApp(
		appPaths{
			app:   *flagAppFile,
			event: *flagEventPath,
		},
	)
	if err := a.run(); err != nil {
		panic(err)
	}
}
