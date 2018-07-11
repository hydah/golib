package render

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type (
	Render interface {
		Render(interface{}, http.ResponseWriter) error
	}
	JSON  struct{}
	JSONP struct {
		Callback string
	}
	XML   struct{}
	TEXT  struct{}
	BYTES struct{}
	HTML  struct{}
	VIDEO struct{}
	IMAGE struct{}
	RAW   struct{}
)

const (
	// ContentType header constant.
	ContentType = "Content-Type"
	// ContentJSON header value for JSON data.
	ContentJSON = "application/json; charset=utf-8"
	// ContentJSONP header value for JSONP data.
	ContentJSONP = "application/javascript"
	// ContentXML header value for XML data.
	ContentXML = "text/xml"
	// ContentPlain header value for Text data.
	ContentPlain = "text/plain; charset=utf-8"
	// ContentPlain header value for Text data.
	ContentHTML = "text/html"
	// ContentPlain header value for Text data.
	ContentVIDEO = "video/mp4"
	// ContentPlain header value for Text data.
	ContentIMAGE = "image/jpeg"
)

// Render an JSON response.
func (c JSON) Render(data interface{}, w http.ResponseWriter) error {
	result, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set(ContentType, ContentJSON)
	w.Write(result)
	return nil
}

// Render an JSONP response.
func (c JSONP) Render(data interface{}, w http.ResponseWriter) error {
	result, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set(ContentType, ContentJSONP)
	w.Write([]byte(c.Callback + "("))
	w.Write(result)
	w.Write([]byte(");"))
	return nil
}

// Render an XML response.
func (c XML) Render(data interface{}, w http.ResponseWriter) error {
	result, err := xml.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set(ContentType, ContentXML)
	w.Write(result)
	return nil
}

// Render an Text response.
func (c TEXT) Render(data interface{}, w http.ResponseWriter) error {
	w.Header().Set(ContentType, ContentPlain)
	_, err := w.Write([]byte(data.(string)))
	return err
}

// Render an Bytes response.
func (c BYTES) Render(data interface{}, w http.ResponseWriter) error {
	w.Header().Set(ContentType, ContentPlain)
	_, err := w.Write(data.([]byte))
	return err
}

// Render an HTML response.
func (c HTML) Render(data interface{}, w http.ResponseWriter) error {
	w.Header().Set(ContentType, ContentHTML)
	_, err := w.Write([]byte(data.(string)))
	return err
}

// Render an HTML response.
func (c VIDEO) Render(data interface{}, w http.ResponseWriter) error {
	w.Header().Set(ContentType, ContentVIDEO)
	_, err := w.Write([]byte(data.([]byte)))
	return err
}

// Render an HTML response.
func (c IMAGE) Render(data interface{}, w http.ResponseWriter) error {
	w.Header().Set(ContentType, ContentIMAGE)
	_, err := w.Write([]byte(data.([]byte)))
	return err
}

// Render an Raw response.
func (c RAW) Render(data interface{}, w http.ResponseWriter) error {
	_, err := w.Write([]byte(data.([]byte)))
	return err
}
