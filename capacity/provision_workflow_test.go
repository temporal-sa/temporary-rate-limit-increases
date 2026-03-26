package capacity

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_ProvisionTRUWorkflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	env.RegisterWorkflow(DeprovisionTRUWorkflow)

	var a *Activities
	env.OnActivity(a.AddTRUs, mock.Anything, mock.Anything).Return(nil)

	input := ProvisionTRUInput{
		Namespace:          "test-namespace",
		APSLimit:           2000,
		MinutesToProvision: 5,
	}

	env.ExecuteWorkflow(ProvisionTRUWorkflow, input)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	env.AssertExpectations(t)
}
