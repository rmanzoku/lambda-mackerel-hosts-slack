package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	slack "github.com/ashwanthkumar/slack-go-webhook"
	"github.com/aws/aws-lambda-go/lambda"
	mkr "github.com/mackerelio/mackerel-client-go"
)

// Event from event.json
type Event map[string]interface{}

func postSlack(webhookurl, channel, text string) error {
	payload := slack.Payload{
		Text:    "```" + text + "```",
		Channel: channel,
	}

	err := slack.Send(webhookurl, "", payload)
	if len(err) > 0 {
		return err[0]
	}

	return nil
}

// HandleRequest is main handler
func HandleRequest(ctx context.Context, event Event) (string, error) {
	text := ""

	client := mkr.NewClient(os.Getenv("MACKEREL_APIKEY"))

	services, err := client.FindServices()
	if err != nil {
		return "", err
	}

	for _, s := range services {
		params := &mkr.FindHostsParam{Service: s.Name}
		hosts, err := client.FindHosts(params)
		if err != nil {
			return "", err
		}

		if len(hosts) == 0 {
			continue
		}

		text += fmt.Sprintf("Service: %s / Count: %d\n", s.Name, len(hosts))

		var hostnames []string
		for _, h := range hosts {
			hostnames = append(hostnames, fmt.Sprintf("%s %s\n", h.Name, h.DisplayName))
		}
		sort.Strings(hostnames)
		for _, l := range hostnames {
			text += l
		}
		text += "\n"
	}

	log.Print(text)

	err = postSlack(os.Getenv("SLACK_WEBHOOK_URL"), os.Getenv("SLACK_CHANNEL"), text)
	if err != nil {
		return "", err
	}

	return "ok", nil
}

func main() {
	lambda.Start(HandleRequest)
}
