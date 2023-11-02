package middleware

import (
	"net/http"
	"net/url"
	"bytes"
	"io"
	"context"
	"github.com/julienschmidt/httprouter"
)

// URL decode the body of the request so there are no issues 
func ParseFormData(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			// Check the content type
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
					http.Error(w, "Unsupported Content-Type", http.StatusBadRequest)
					return
			}

			// Create a buffer to store the request body
			var buffer bytes.Buffer
			bodyReader := io.TeeReader(r.Body, &buffer)

			// Read the request body into a byte slice
			b, err := io.ReadAll(bodyReader)
			if err != nil {
					http.Error(w, "Error reading request body", http.StatusBadRequest)
					return
			}

			// Reset the request body with the buffer
			r.Body = io.NopCloser(&buffer)

			// Decode the URL-encoded data
			decodedData, err := url.QueryUnescape(string(b))
			if err != nil {
					http.Error(w, "Error decoding URL-encoded data", http.StatusBadRequest)
					return
			}

			// Set the decoded data as a value in the request context
			ctx := context.WithValue(r.Context(), "body", decodedData)

			// Create a new request with the updated context
			r = r.WithContext(ctx)

			// Call the next handler with the updated request
			n(w, r, ps)
	}
}