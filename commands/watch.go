package commands

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/concourse/atc"
	"github.com/concourse/fly/eventstream"
	"github.com/tedsuo/rata"
	"github.com/vito/go-sse/sse"
)

func Watch(c *cli.Context) {
	atcURL := c.GlobalString("atcURL")
	insecure := c.GlobalBool("insecure")

	atcRequester := newAtcRequester(atcURL, insecure)

	build := getBuild(c, atcRequester.httpClient, atcRequester.RequestGenerator)

	eventSource := &sse.EventSource{
		Client: atcRequester.httpClient,
		CreateRequest: func() *http.Request {
			logOutput, err := atcRequester.CreateRequest(
				atc.BuildEvents,
				rata.Params{"build_id": strconv.Itoa(build.ID)},
				nil,
			)
			if err != nil {
				log.Fatalln(err)
			}

			return logOutput
		},
	}

	exitCode, err := eventstream.RenderStream(eventSource)
	if err != nil {
		log.Println("failed to render stream:", err)
		os.Exit(1)
	}

	eventSource.Close()

	os.Exit(exitCode)
}
