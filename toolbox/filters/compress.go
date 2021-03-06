package filters

import (
	"compress/flate"
	"compress/gzip"
	"strings"

	"github.com/cosiner/zerver"
)

type gzipResponse struct {
	gzipWriter *gzip.Writer
	zerver.Response
}

type flateResponse struct {
	flateWriter *flate.Writer
	zerver.Response
}

func (gr *gzipResponse) Write(data []byte) (int, error) {
	return gr.gzipWriter.Write(data)
}

func (fr *flateResponse) Write(data []byte) (int, error) {
	return fr.flateWriter.Write(data)
}

func CompressFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	encoding := req.AcceptEncodings()
	if strings.Contains(encoding, zerver.ENCODING_GZIP) {
		gzw := gzip.NewWriter(resp)
		gresp := gzipResponse{gzw, resp}
		resp.SetContentEncoding(zerver.ENCODING_GZIP)
		defer gzw.Close()
		chain(req, &gresp)
		gzw.Close()
	} else if strings.Contains(encoding, zerver.ENCODING_DEFLATE) {
		flw, _ := flate.NewWriter(resp, flate.DefaultCompression)
		fresp := flateResponse{flw, resp}
		resp.SetContentEncoding(zerver.ENCODING_DEFLATE)
		defer flw.Close()
		chain(req, &fresp)
	} else {
		chain(req, resp)
	}
}
