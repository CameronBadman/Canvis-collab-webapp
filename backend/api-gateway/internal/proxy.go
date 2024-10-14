package internal

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// NewProxy returns a new HTTP handler that proxies requests to the specified service
func NewProxy(serviceURL string, pathPrefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a new URL for the service
		serviceURLParsed, _ := url.Parse(serviceURL)

		// Update the request URL path to remove the prefix
		r.URL.Path = strings.TrimPrefix(r.URL.Path, pathPrefix)

		// Update the request to point to the service URL
		r.URL.Scheme = serviceURLParsed.Scheme
		r.URL.Host = serviceURLParsed.Host
		r.Host = serviceURLParsed.Host

		// Forward the request to the service
		resp, err := http.DefaultTransport.RoundTrip(r)
		if err != nil {
			http.Error(w, "Error communicating with service", http.StatusBadGateway)
			log.Printf("Error communicating with service: %v", err)
			return
		}
		defer resp.Body.Close()

		log.Printf("Received response from service: %d\n", resp.StatusCode)

		// Read response from service
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading response from service", http.StatusInternalServerError)
			log.Printf("Error reading response body: %v", err)
			return
		}

		// Copy response headers and status code
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	}
}
