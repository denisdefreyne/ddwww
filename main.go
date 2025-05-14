package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
)

const usage = `Usage:
    ddwww [-r PATH] [-p PORT]

Options:
    -r, --root     Specify directory to serve (default: .).
    -p, --port     Specify port to serve on (default: 3000).
    -v, --version  Show application version (default: false).`

func main() {
	// Get working directory (default)
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// Define flags
	var (
		rootFlag    string
		portFlag    int
		versionFlag bool
	)

	// Parse flags
	flag.StringVar(&rootFlag, "root", wd, "Specify directory to serve")
	flag.StringVar(&rootFlag, "r", wd, "Specify directory to serve")
	flag.IntVar(&portFlag, "port", 3000, "Specify port to serve on")
	flag.IntVar(&portFlag, "p", 3000, "Specify port to serve on")
	flag.BoolVar(&versionFlag, "version", false, "Show version")
	flag.BoolVar(&versionFlag, "v", false, "Show version")
	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }
	flag.Parse()

	// Handle version flag
	if versionFlag {
		fmt.Printf(
			"ddwww %s\ncommit %s\nbuilt at %s\n",
			buildVersion, buildCommit, buildDate)
		os.Exit(0)
	}

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
