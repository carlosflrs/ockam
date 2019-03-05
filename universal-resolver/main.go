package main

import (
	"fmt"
	"github.com/ockam-network/did"
	"github.com/ockam-network/ockam/node"
	ockamHttp "github.com/ockam-network/ockam/node/remote/http"
	"log"
	"net/http"
	"os"
	"regexp"
)

func run() error {
	mux := http.NewServeMux()
	mux.Handle("/hello/", handleHello())
	mux.Handle(fmt.Sprintf("/%s/identifiers/ockam/", Version()), handleGetEntity())
	mux.Handle(fmt.Sprintf("/%s/identifiers/ockam/{ockam}/claim/{claim}", Version()), handleGetClaim())

	log.Println("Listening on", "8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		return err
	}

	return nil
}

func handleHello() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("handleHello")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	})
}

func handleGetEntity() http.Handler {
	// test.ockam.network needs to be replaced
	ockamNode, err := node.New(node.PeerDiscoverer(ockamHttp.Discoverer("test.ockam.network", 26657)))
	exitOnError(err)

	err = ockamNode.Sync()
	exitOnError(err)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ockamChain := ockamNode.Chain()

		// print some information about the chain
		fmt.Printf("Chain ID: %s\n", ockamChain.ID())

		regex, _ := regexp.Compile("did:ockam:[a-zA-Z0-9]*")

		// Fetch Entity
		id, err := did.Parse(regex.FindString(r.URL.Path))
		if err != nil {
			exitOnError(err)
		}
		fmt.Printf("did: %s\n", id.String())
		bytes, _, err := ockamNode.FetchEntity(id.String())
		//bytes, _, err := ockamNode.FetchClaim(id.String())

		if err != nil {
			exitOnError(err)
		}

		respondWithJson(w, r, http.StatusOK, bytes)
	})
}

func handleGetClaim() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
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
