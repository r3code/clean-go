package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSONAPIErrorReporter catches errors, transaltes them to HTTPAPIError if possible or wraps as InternalServerError if err not found in registered HTTPAPIErrors, then sends as it as JSON
func JSONAPIErrorReporter() gin.HandlerFunc {
	return jsonAPIErrorReporterT(gin.ErrorTypeAny)
}

func jsonAPIErrorReporterT(errType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedErrors := c.Errors.ByType(errType)
		if len(detectedErrors) > 0 {
			var aerr *HTTPAPIError
			err := detectedErrors[0].Err
			switch err.(type) {
			case *HTTPAPIError:
				aerr = err.(*HTTPAPIError)
			default:
				ok := false
				aerr, ok = GetHTTPAPIError(err)
				if !ok {
					// send error as JSON
					aerr = NewHTTPAPIError(http.StatusInternalServerError, err.Error(), "InternalServerError", "###")
				}
			}

			c.IndentedJSON(aerr.Code, gin.H{
				"error": aerr})
			return
		}

	}
}
