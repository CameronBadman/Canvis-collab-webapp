package internal

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// NewProxy returns a new HTTP handler that proxies requests to the specified service
func NewProxy(serviceURL string, pathPrefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Proxy: Original request URL: %s", r.URL.String())
		log.Printf("Proxy: Service URL: %s", serviceURL)
		log.Printf("Proxy: Path prefix: %s", pathPrefix)

		// Create a new URL for the service
		serviceURLParsed, err := url.Parse(serviceURL)
		if err != nil {
			http.Error(w, "Invalid service URL", http.StatusInternalServerError)
			log.Printf("Error parsing service URL: %v", err)
			return
		}

		// Update the request URL path to remove the prefix
		r.URL.Path = strings.TrimPrefix(r.URL.Path, pathPrefix)
		log.Printf("Proxy: Updated request path: %s", r.URL.Path)

		// Create a new request to forward
		proxyReq, err := http.NewRequest(r.Method, serviceURLParsed.String()+r.URL.Path, r.Body)
		if err != nil {
			http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
			log.Printf("Error creating proxy request: %v", err)
			return
		}

		log.Printf("Proxy: Final proxy request URL: %s", proxyReq.URL.String())

		// Copy headers from original request to proxy request
		copyHeader(proxyReq.Header, r.Header)

		// Forward the request to the service
		resp, err := http.DefaultClient.Do(proxyReq)
		if err != nil {
			http.Error(w, "Error communicating with service", http.StatusBadGateway)
			log.Printf("Error communicating with service: %v", err)
			return
		}
		defer resp.Body.Close()

		log.Printf("Received response from service: %d\n", resp.StatusCode)

		// Copy response headers
		copyHeader(w.Header(), resp.Header)

		// Write the status code
		w.WriteHeader(resp.StatusCode)

		// Copy the response body
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Printf("Error copying response body: %v", err)
		}
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
