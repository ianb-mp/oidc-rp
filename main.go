package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"golang.org/x/net/context"
)

const (
	// issuerURL corresponds with `iss` field
	issuerURL = "https://token.actions.githubusercontent.com"
	// clientID corresponds with `aud` field, which will be the Github repo
	clientID = "https://github.com/ianb-mp"
	// limit is the maximum size of the request body
	limit = 1024 * 1024
)

func main() {
	// Define a flag for the listen port
	port := flag.Int("port", 8080, "HTTP listen port")
	flag.Parse()

	// Create a new context
	ctx := context.Background()

	// Set up the OpenID Connect configuration
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		log.Fatal(err)
	}

	// Set up the verifier
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})

	// Define the HTTP handler for the JWT endpoint
	http.HandleFunc("/jwt", func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(io.LimitReader(r.Body, limit))
		if err != nil {
			fmt.Printf("Failed to read JWT: %v", err)
			http.Error(w, "Failed to read JWT", http.StatusBadRequest)
			return
		}

		// Parse the JWT from the request body
		rawIDToken := string(body)

		log.Printf("rawIDToken: %q", rawIDToken)

		// Verify the JWT
		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			fmt.Printf("Failed to verify JWT: %v", err)
			http.Error(w, "Failed to verify JWT", http.StatusBadRequest)
			return
		}

		// dump token
		PrintJSON(idToken)
		//fmt.Printf("%+v", idToken)
	})

	// Start the HTTP server
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// https://www.reddit.com/r/golang/comments/gritgv/comment/fs3gdtu/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button
func PrintJSON(obj interface{}) {
	bytes, _ := json.MarshalIndent(obj, "", "  ")
	fmt.Println(string(bytes))
}
