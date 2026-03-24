package main

import (
	"github.com/temporal-sa/temporary-rate-limit-increases/capacity"
	"go.temporal.io/sdk/client"
	"log"
	"os"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	targetNamespace := os.Getenv("TEMPORAL_CLOUD_NAMESPACE")
	var newLimit int32 = 1000

	err = capacity.HandleProvisionRequest(c, targetNamespace, newLimit)
	if err != nil {
		log.Fatalln("Failed to execute provision request", err)
	}
}
