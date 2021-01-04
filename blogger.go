package blogger

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const (
	//holds our static assets e.g css,javascript
	staticDir = "/static/"
	//holds our html markup files.
	pagesDir = "pages"
	//relative path of the static directory
	staticRelPath = "./static"
)

//Start the http server
func Start() {

	portAddr := flag.String("port", "3000", "HTTP listen port address")

	flag.Parse()

	fs := http.FileServer(http.Dir(staticRelPath))

	mux := http.NewServeMux()

	mux.Handle(staticDir, http.StripPrefix(staticDir, fs))
	mux.HandleFunc("/", serveHTML)

	srv := http.Server{
		Addr:         fmt.Sprintf(":%v", *portAddr),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	sigc := make(chan os.Signal, 1) //buffered channel

	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	errc := make(chan error, 1) //buffered channel

	go func() {
		fmt.Printf("Server listening on address %v\n", srv.Addr)
		errc <- srv.ListenAndServe()
	}()

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

func serveHTML(w http.ResponseWriter, r *http.Request) {
	fPath := filepath.Join(pagesDir, filepath.Clean(r.URL.Path))

	info, err := os.Stat(fPath)
	if err != nil {
		//If it is a os.PathError (file does not exist)
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	//Do not serve a directory here
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tpl, err := template.ParseFiles(fPath)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
