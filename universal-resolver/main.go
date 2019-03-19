package main

import (
	"fmt"
	"github.com/ockam-network/did"
	"github.com/ockam-network/ockam/node"
	ockamHttp "github.com/ockam-network/ockam/node/remote/http"
	"log"
	"net/http"
	"os"
	"strings"
)

func run() error {
	mux := http.NewServeMux()
	mux.Handle(fmt.Sprintf("/%s/identifiers/ockam/", Version()), handleGetEntity())

	log.Println("Listening on", "8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		return err
	}

	return nil
}

func handleGetEntity() http.Handler {
	// test.ockam.network needs to be replaced
	ockamNode, err := node.New(node.PeerDiscoverer(ockamHttp.Discoverer("test.ockam.network", 26657)))
	exitOnError(err)

	err = ockamNode.Sync()
	exitOnError(err)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		did_path := strings.Split(r.URL.Path, fmt.Sprintf("/%s/identifiers/ockam/", Version()))

		id, err := did.Parse(did_path[1])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
			return
		}

		//Fetch Entity
		fmt.Printf("Entity did: %s\n", id.String())
		bytes, _, err := ockamNode.FetchEntity(id.String())

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
			return
		}

		respondWithJson(w, r, http.StatusOK, bytes)
	})
}

func respondWithJson(w http.ResponseWriter, r *http.Request, code int, payload []byte) {
	w.Header().Set("Content-Type", "application/json-ld")
	w.WriteHeader(code)
	w.Write(payload)
}

func main() {
	log.Fatal(run())
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

// Version returns the current version of Resolver
func Version() string {
	version := "1.0"
	return version
}
