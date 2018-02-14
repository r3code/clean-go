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
			e := detectedErrors[0].Err
			aerr, ok := GetHTTPAPIError(e)
			if !ok {
				// send error as JSON
				internalError := NewHTTPAPIError(http.StatusInternalServerError, e.Error(), "InternalServerError", "")
				c.IndentedJSON(internalError.Code, internalError)
				return
			}

			c.IndentedJSON(aerr.Code, aerr)
			return
		}

	}
}
