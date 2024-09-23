package main

import (
	"bytes"
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
	// clientID corresponds with `aud` field, which will be set to the Github user/org
	clientID = "https://github.com/ianb-mp"
	// limit is the maximum size of the request body
	limit = 1024 * 1024
)

type JWTRequest struct {
	JWT string `json:"jwt"`
}

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

		var req JWTRequest

		// Decode the request body into the JWTRequest struct
		err := json.NewDecoder(io.LimitReader(r.Body, limit)).Decode(&req)
		if err != nil {
			log.Printf("Failed to decode request body: %v", err)
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		// Verify the JWT
		idToken, err := verifier.Verify(ctx, req.JWT)
		if err != nil {
			log.Printf("Failed to verify JWT: %v", err)
			http.Error(w, "Failed to verify JWT", http.StatusBadRequest)
			return
		}

		buf := bytes.NewBuffer(nil)
		PrintJSON(buf, idToken)
		log.Printf("raw token: %s", req.JWT)
		log.Printf("JWT verified: %s", buf.String())

		// check other claims like subject

		// return access token
	})

	// Start the HTTP server
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// https://www.reddit.com/r/golang/comments/gritgv/comment/fs3gdtu/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button
func PrintJSON(w io.Writer, obj interface{}) {
	bytes, _ := json.MarshalIndent(obj, "", "  ")
	fmt.Fprintf(w, "%s\n", string(bytes))
}
