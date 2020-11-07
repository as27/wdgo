package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

var (
	flagLogFile   = flag.String("logFile", "", "write log into log.txt")
	flagAppFile   = flag.String("file", "wdgo.txt", "define the file where the boardlist is stored")
	flagEventPath = flag.String("events", "events", "path where the events of each board are stored")
)

func main() {
	flag.Parse()
	buf := bufio.NewWriter(os.Stdout)
	defer buf.Flush()
	log.SetOutput(buf)
	if *flagLogFile != "" {
		logfile, _ := os.OpenFile("log.txt", os.O_CREATE, 0777)
		defer logfile.Close()
		log.SetOutput(logfile)
	}
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
