package main

import (
	"fmt"

	"github.com/caseyhird/nginkc-go/nginkc"
)

type MyApp struct {
	Name string
}

func (app MyApp) Call(req nginkc.Request) nginkc.Response {
	return nginkc.Response{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "text/plain"},
		Body:       fmt.Sprintf("App %s called with %s", app.Name, req.String()),
	}
}

func main() {
	app := MyApp{
		Name: "test app",
	}
	nginkc.Serve(app)
}
