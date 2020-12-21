package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var (
	portAddr   = flag.String("http", ":8080", "HTTP listen address")
	contentDir = flag.String("content", "content/", "Content directory")
	index      = flag.Bool("index", false, "Show index")
)

func main() {
	flag.Parse()

	if strings.Index(*portAddr, ":") != 0 {
		fmt.Printf("invalid port in address")
		flag.Usage()
		os.Exit(1)
	}

	info, err := os.Stat(*contentDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//holds our html files and their paths
	cmap := map[string]string{}

	//check if the content passed in is a directory. if it is a directory we traverse its contents.pulling out only html files
	if isDir := info.IsDir(); isDir {
		files, err := ioutil.ReadDir(*contentDir)
		if err != nil {
			fmt.Printf("ioutil.ReadDir(%s): error %s", *contentDir, err)
		}

		for _, f := range files {
			//We are only interested in html files for now
			if strings.HasSuffix(f.Name(), ".html") {
				// fmt.Println(f.Name())
				fp := filepath.Join(*contentDir, f.Name())
				key := strings.Split(f.Name(), ".")[0]
				cmap[key] = fp
			}
		}

		// fmt.Printf("%+v\n", cmap)
	}

	mux := http.NewServeMux()

	for k, v := range cmap {
		//Let's show index here instead of using another loop
		if *index {
			fmt.Printf("/%s\n", k)
		}
		//construct our routes here
		mux.HandleFunc(fmt.Sprintf("/%s", k), serveFile(v))
	}

	srv := http.Server{
		Addr:         *portAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	sigc := make(chan os.Signal, 1) //buffered channel

	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	errc := make(chan error, 1) //buffered channel
	go func(portAddr *string) {
		fmt.Printf("Server listening on address %v\n", *portAddr)
		errc <- srv.ListenAndServe()
	}(portAddr)

	for {
		select {
		case sig := <-sigc:

			fmt.Printf("Received shutdown signal: %v\n", sig)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				srv.Close()
				fmt.Printf("could not shut down server gracefully: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Server shutting down gracefully")
			os.Exit(0)

		case err := <-errc:
			if err != http.ErrServerClosed {
				fmt.Printf("server error: %v\n", err)
				os.Exit(1)
			}
		}
	}

}

func serveFile(path string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}
