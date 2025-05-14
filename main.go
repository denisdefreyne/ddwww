package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const usage = `Usage:
    ddwww [-r PATH] [-p PORT]

Options:
    -r, --root  Specify directory to serve.
    -p, --port  Specify port to serve on.

PORT defaults to 3000, and ROOT defaults to the current working directory.`

func main() {
	// Get working directory (default)
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// Define flags
	var (
		rootFlag string
		portFlag int
	)

	// Parse flags
	flag.StringVar(&rootFlag, "root", wd, "Specify directory to serve")
	flag.StringVar(&rootFlag, "r", wd, "Specify directory to serve")
	flag.IntVar(&portFlag, "port", 3000, "Specify port to serve on")
	flag.IntVar(&portFlag, "p", 3000, "Specify port to serve on")
	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }
	flag.Parse()

	// Set up HTTP handler
	fs := http.FileServer(http.Dir(rootFlag))
	http.Handle("/", loggingMiddleware(noCachingMiddleware(fs)))

	// Serve
	addr := fmt.Sprintf("localhost:%v", portFlag)
	log.Printf("Serving at <http://%v>", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
