package filters

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/cosiner/zerver"
)

type gzipWriter struct {
	gw *gzip.Writer
	http.ResponseWriter
	needClose bool
}

func (w *gzipWriter) Write(data []byte) (int, error) {
	return w.gw.Write(data)
}

func (w *gzipWriter) Close() error {
	err := w.gw.Close()
	if w.needClose {
		err = w.ResponseWriter.(io.Closer).Close()
	}
	return err
}

type flateWriter struct {
	fw *flate.Writer
	http.ResponseWriter
	needClose bool
}

func (w *flateWriter) Write(data []byte) (int, error) {
	return w.fw.Write(data)
}

func (w *flateWriter) Close() error {
	err := w.fw.Close()
	if w.needClose {
		err = w.ResponseWriter.(io.Closer).Close()
	}
	return err
}

func gzipWrapper(w http.ResponseWriter, needClose bool) (http.ResponseWriter, bool) {
	return &gzipWriter{
		gw:             gzip.NewWriter(w),
		ResponseWriter: w,
		needClose:      needClose,
	}, true
}

func flateWrapper(w http.ResponseWriter, needClose bool) (http.ResponseWriter, bool) {
	fw, _ := flate.NewWriter(w, flate.DefaultCompression)
	return &flateWriter{
		fw:             fw,
		ResponseWriter: w,
		needClose:      needClose,
	}, true
}

func CompressFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	encoding := req.AcceptEncodings()
	if strings.Contains(encoding, zerver.ENCODING_GZIP) {
		resp.SetContentEncoding(zerver.ENCODING_GZIP)
		resp.Wrap(gzipWrapper)
	} else if strings.Contains(encoding, zerver.ENCODING_DEFLATE) {
		resp.SetContentEncoding(zerver.ENCODING_DEFLATE)
		resp.Wrap(flateWrapper)
	}
	chain(req, resp)
}
