package cache

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

type responseCaptureWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func newResponseCaptureWriter(w gin.ResponseWriter) *responseCaptureWriter {
	return &responseCaptureWriter{
		ResponseWriter: w,
		body:           bytes.NewBuffer(nil),
	}
}

func (w *responseCaptureWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}
