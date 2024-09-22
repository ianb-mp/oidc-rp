package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"golang.org/x/net/context"
)

const issuerURL = "https://token.actions.githubusercontent.com"
const clientID = "foobar"

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
		// Parse the JWT from the request body
		rawIDToken := r.FormValue("jwt")

		// Verify the JWT
		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			http.Error(w, "Failed to verify JWT", http.StatusBadRequest)
			return
		}

		// Print the subject of the JWT
		fmt.Fprintf(w, "Subject: %s", idToken.Subject)
	})

	// Start the HTTP server
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
