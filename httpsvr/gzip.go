package httpsvr

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"net"
	"net/http"
	"strings"
)

const (
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderContentEncoding = "Content-Encoding"
	HeaderContentLength   = "Content-Length"
	HeaderContentType     = "Content-Type"
	HeaderVary            = "Vary"

	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

// Gzip returns a Handler that adds gzip compression to all requests
func Gzip(compressionLevel int) HandlerFunc {
	return func(ctx *Context) {
		if !strings.Contains(ctx.Req.Header.Get(HeaderAcceptEncoding), "gzip") {
			return
		}

		ctx.Writer.Header().Set(HeaderContentEncoding, "gzip")
		ctx.Writer.Header().Set(HeaderVary, HeaderAcceptEncoding)
		gz := gzip.NewWriter(ctx.Writer)

		gz, err := gzip.NewWriterLevel(ctx.Writer, compressionLevel)
		if err != nil {
			ctx.Next()
			return
		}
		defer gz.Close()

		gzw := gzipResponseWriter{gz, ctx.Writer}
		ctx.Writer = gzw
		ctx.Next()
	}
}

type gzipResponseWriter struct {
	w *gzip.Writer
	ResponseWriter
}

func (c gzipResponseWriter) Write(p []byte) (int, error) {
	if len(c.Header().Get(HeaderContentType)) == 0 {
		c.Header().Set(HeaderContentType, http.DetectContentType(p))
	}
	return c.w.Write(p)
}

func (grw gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := grw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}
