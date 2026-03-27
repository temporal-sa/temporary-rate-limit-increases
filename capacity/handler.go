package capacity

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
)

const TaskQueue = "capacity-management"

// HandleProvisionRequest starts a ProvisionTRUWorkflow execution and returns
// after the Temporal Cluster accepts the request. It does not wait
// for the Workflow to complete.
func HandleProvisionRequest(c client.Client, namespace string, apsLimit int32) error {
	options := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("provision-%s", namespace),
		TaskQueue: TaskQueue,
	}

	input := ProvisionTRUInput{
		Namespace:          namespace,
		APSLimit:           apsLimit,
		MinutesToProvision: 5,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, ProvisionTRUWorkflow, input)
	if err != nil {
		return fmt.Errorf("failed to start workflow: %w", err)
	}

	fmt.Printf("Successfully started provision workflow. ID: %s, RunID: %s\n",
		we.GetID(), we.GetRunID())
	return nil
}
