package capacity

import (
	"go.temporal.io/sdk/temporal"
	"time"

	"go.temporal.io/sdk/workflow"
)

// DeprovisionTRUWorkflow sleeps for the TTL duration then reverts the namespace
// to on-demand capacity. It is started as an asynchronous Child Workflow by
// ProvisionTRUWorkflow and runs independently after the parent completes.
func DeprovisionTRUWorkflow(ctx workflow.Context, input DeprovisionTRUInput) error {
	err := workflow.Sleep(ctx, time.Duration(input.MinutesToProvision)*time.Minute)
	if err != nil {
		return err
	}

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{unauthorized, forbidden, badRequest},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var a *Activities
	return workflow.ExecuteActivity(ctx, a.RemoveTRUs, input).Get(ctx, nil)
}
