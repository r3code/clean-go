package web

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/r3code/clean-go/engine"
)

type (
	greetingManager struct {
		engine.GreetingManager
	}
)

// wire up the greetings routes
func initGreetings(e *gin.Engine, f engine.ServiceCreator, endpoint string) {
	greeter := &greetingManager{f.NewGreetingManager()}
	g := e.Group(endpoint)
	{
		g.GET("", greeter.listGreeting)
		g.POST("", greeter.addGreeting)
	}
}

// list converts the parameters into an engine
// request and then marshalls the results based
// on the format requested, returning either an
// html rendered page or JSON (to simulate basic
// content negotiation). It's simpler if the UI
// is a SPA and the web interface is just an API.
func (gm greetingManager) listGreeting(c *gin.Context) {
	ctx := getContext(c)
	count, err := strconv.Atoi(c.Query("count"))
	if err != nil || count == 0 {
		count = 5
	}
	req := &engine.ListGreetingsRequest{
		Count: count,
	}
	res, err := gm.ListGreetings(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}
	if c.Query("format") == "json" {
		c.JSON(http.StatusOK, res.Greetings)
	} else {
		c.HTML(http.StatusOK, "guestbook.html", res)
	}
	log.Printf("List Result: %v\n", res)
}

// add accepts a form post and creates a new
// greoting in the system. It could be made a
// lot smarter and automatically check for the
// content type to handle forms, JSON etc...
func (gm greetingManager) addGreeting(c *gin.Context) {
	ctx := getContext(c)
	req := &engine.AddGreetingRequest{
		Author:  c.PostForm("Author"),
		Content: c.PostForm("Content"),
	}
	res, err := gm.CreateGreeting(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}
	// TODO: set a flash cookie for "added"
	// if this was a web request, otherwise
	// send a nice JSON response ...
	log.Printf("Add Result: %v\n", res)
	c.Redirect(http.StatusFound, "/")
}
