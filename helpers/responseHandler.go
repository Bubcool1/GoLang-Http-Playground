package helpers

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"

	"gopkg.in/yaml.v3"
)

func FormatResponse(w http.ResponseWriter, content any, req *http.Request) error {
	contentType := req.Header.Get("Content-Type")

	// Default to JSON if not set
	if contentType == "" {
		contentType = "application/json"
	}

	// Validate it's an allowed type
	allowedTypes := map[string]convertFunc{
		"application/json": convertToJSON,
		"application/xml":  convertToXML,
		"application/yaml": convertToYAML,
		// "text/csv":         true,
		// "text/html":        true,
		// "text/plain":       true,
	}

	if allowedTypes[contentType] == nil {
		return errors.New("invalid content type")
	}

	w.Header().Set("Content-Type", contentType)
	if err := allowedTypes[contentType](w, content); err != nil {
		return err
	}
	println("200 status code, response sent: " + req.RequestURI)

	return nil
}

type convertFunc func(w http.ResponseWriter, content any) error

func convertToJSON(w http.ResponseWriter, content any) error {
	json.NewEncoder(w).Encode(StructToMap(content))

	return nil
}

func convertToXML(w http.ResponseWriter, content any) error {
	xml.NewEncoder(w).Encode(content)

	return nil
}

func convertToYAML(w http.ResponseWriter, content any) error {
	yaml.NewEncoder(w).Encode(StructToMap(content))

	return nil
}
