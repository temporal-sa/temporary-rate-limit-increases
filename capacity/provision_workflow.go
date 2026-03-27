package capacity

import (
	"fmt"
	"go.temporal.io/sdk/temporal"
	"net/http"
	"strconv"
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

var (
	unauthorized = strconv.Itoa(http.StatusUnauthorized)
	forbidden    = strconv.Itoa(http.StatusForbidden)
)

// ProvisionTRUWorkflow raises the namespace capacity limit, then starts the
// DeprovisionTRUWorkflow as an asynchronous Child Workflow with an abandon
// policy so the parent can return immediately while cleanup runs independently.
func ProvisionTRUWorkflow(ctx workflow.Context, input ProvisionTRUInput) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{unauthorized, forbidden},
		},
	}
	actCtx := workflow.WithActivityOptions(ctx, ao)

	var a *Activities
	err := workflow.ExecuteActivity(actCtx, a.AddTRUs, AddTRUInput{
		Namespace: input.Namespace,
		APSLimit:  input.APSLimit,
	}).Get(ctx, nil)
	if err != nil {
		return err
	}

	cwo := workflow.ChildWorkflowOptions{
		WorkflowID:        fmt.Sprintf("deprovision-%s", input.Namespace),
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}
	childCtx := workflow.WithChildOptions(ctx, cwo)

	childFuture := workflow.ExecuteChildWorkflow(childCtx, DeprovisionTRUWorkflow, DeprovisionTRUInput{
		Namespace:          input.Namespace,
		MinutesToProvision: input.MinutesToProvision,
	})

	var childExec workflow.Execution
	err = childFuture.GetChildWorkflowExecution().Get(ctx, &childExec)
	if err != nil {
		return err
	}

	return nil
}
