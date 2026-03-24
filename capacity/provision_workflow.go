package capacity

import (
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

// ProvisionTRUWorkflow raises the namespace capacity limit, then starts the
// DeprovisionTRUWorkflow as an asynchronous Child Workflow with an abandon
// policy so the parent can return immediately while cleanup runs independently.
func ProvisionTRUWorkflow(ctx workflow.Context, input ProvisionInput) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
	}
	actCtx := workflow.WithActivityOptions(ctx, ao)

	var a *Activities
	err := workflow.ExecuteActivity(actCtx, a.AddTRUs, input).Get(ctx, nil)
	if err != nil {
		return err
	}

	cwo := workflow.ChildWorkflowOptions{
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}
	childCtx := workflow.WithChildOptions(ctx, cwo)

	childFuture := workflow.ExecuteChildWorkflow(childCtx, DeprovisionTRUWorkflow, input)

	var childExec workflow.Execution
	err = childFuture.GetChildWorkflowExecution().Get(ctx, &childExec)
	if err != nil {
		return err
	}

	return nil
}
