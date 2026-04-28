package main

import (
	"github.com/temporal-sa/temporary-rate-limit-increases/capacity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	"os"
)

func main() {
	apiKey := os.Getenv("TEMPORAL_CLOUD_API_KEY")
	if apiKey == "" {
		log.Fatalln("TEMPORAL_CLOUD_API_KEY missing and required")
	}

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, capacity.TaskQueue, worker.Options{})

	w.RegisterWorkflow(capacity.ProvisionTRUWorkflow)
	w.RegisterWorkflow(capacity.DeprovisionTRUWorkflow)

	activities := &capacity.Activities{
		APIKey: apiKey,
	}
	w.RegisterActivity(activities)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
