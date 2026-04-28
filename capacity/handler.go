package capacity

import (
	"context"
	"errors"
	"fmt"
	"go.temporal.io/sdk/client"
)

const TaskQueue = "capacity-management"

// HandleProvisionRequest starts a ProvisionTRUWorkflow execution and returns
// after the Temporal Cluster accepts the request. It does not wait
// for the Workflow to complete.
func HandleProvisionRequest(c client.Client, namespace string, apsLimit int32, minutesToProvision int32) error {
	id := generateProvisioningId(namespace)
	preExistingRun := c.GetWorkflow(context.Background(), id, "").GetRunID()
	if preExistingRun != "" {
		return errors.New("provisioning request already in-progress")
	}
	deprovisionId := fmt.Sprintf("deprovision-%s", namespace)
	preExistingDeprovisionRun := c.GetWorkflow(context.Background(), deprovisionId, "").GetRunID()
	if preExistingDeprovisionRun != "" {
		err := c.CancelWorkflow(context.Background(), deprovisionId, preExistingDeprovisionRun)
		if err != nil {
			return errors.Join(errors.New("unable to cancel deprovisioning workflow"), err)
		}
	}
	options := client.StartWorkflowOptions{
		ID:        id,
		TaskQueue: TaskQueue,
	}

	input := ProvisionTRUInput{
		Namespace:          namespace,
		APSLimit:           apsLimit,
		MinutesToProvision: minutesToProvision,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, ProvisionTRUWorkflow, input)
	if err != nil {
		return fmt.Errorf("failed to start workflow: %w", err)
	}

	fmt.Printf("Successfully started provision workflow. ID: %s, RunID: %s\n",
		we.GetID(), we.GetRunID())
	return nil
}
