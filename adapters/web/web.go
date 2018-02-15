package web

import (
	"errors"
	"fmt"
	"net/http"

	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"

	"github.com/r3code/clean-go/engine"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// NewWebAdapter creates a new web adaptor which will
// handle the web interface and pass calls on to the
// engine to do the real work (that's why the engine
// factory is passed in - so anything that *it* needs
// is unknown to this).
// Because the web adapter ends up quite lightweight
// it easier to replace. We could use any one of the
// Go web routers / frameworks (Gin, Echo, Goji etc...)
// or just stick with the standard framework. Changing
// should be far less costly.
func NewWebAdapter(e engine.ServiceCreator, log bool) http.Handler {
	var router *gin.Engine
	if log {
		router = gin.Default()
	} else {
		router = gin.New()
	}
	// router.Use(gin.ErrorLogger())
	router.Use(JSONAPIErrorReporter())
	// router.Use(gin.Recovery()) default recover for unexpected errors
	router.Use(nice.Recovery(recoveryHandler))
	router.LoadHTMLGlob("templates/*")

	router.GET("/test/error", func(c *gin.Context) {
		c.Error(errors.New("Some unregisterd error in app"))
		return
	})

	initGreetingManager(router, e, "/greetings")

	return router
}

func recoveryHandler(c *gin.Context, err interface{}) {
	apiError := NewHTTPAPIError(http.StatusForbidden, "Internal Server Error", "UnexpectedError", "")
	apiError.DebugInfo = fmt.Sprintf("%#v", err)
	c.IndentedJSON(apiError.Code, gin.H{
		"error": apiError})
}
