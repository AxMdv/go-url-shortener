// Package compressor provides gzip compressing and decompressing of data.
package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
)

// compressWriter implements the http.ResponseWriter and allows
// compressing the transmitted data for the server and set the correct HTTP headers.
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// NewCompressWriter returns a new compressWriter.
func NewCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header returns http.Header.
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write writes compressed data.
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)

}

// WriteHeader sets response status code.
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 || statusCode == 409 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close closes compress writer.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressWriter implements the io.ReadCloser and allows
// decompressing the data received from the client transparently for the server.
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// NewCompressReader returns new compressReader.
func NewCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read reads gzip-compressed data.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close closes gzip reader.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
