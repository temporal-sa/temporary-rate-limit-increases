package main

import (
	"github.com/temporal-sa/temporary-rate-limit-increases/capacity"
	"go.temporal.io/sdk/client"
	"log"
	"os"
	"strconv"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	targetNamespace := os.Getenv("TEMPORAL_CLOUD_NAMESPACE")
	if targetNamespace == "" {
		log.Fatalln("TEMPORAL_CLOUD_NAMESPACE missing and required")
	}
	minutesToProvisionRaw := os.Getenv("MINUTES_TO_PROVISION")
	if minutesToProvisionRaw == "" {
		log.Fatalln("MINUTES_TO_PROVISION missing and required")
	}
	minutesToProvisionParsed, err := strconv.ParseInt(minutesToProvisionRaw, 10, 32)
	if err != nil {
		log.Fatalln("Unable to parse MINUTES_TO_PROVISION: " + err.Error())
	}
	minutesToProvision := int32(minutesToProvisionParsed)
	var newLimit int32 = 1000

	err = capacity.HandleProvisionRequest(c, targetNamespace, newLimit, minutesToProvision)
	if err != nil {
		log.Fatalln("Unable to execute provision request", err)
	}
}
