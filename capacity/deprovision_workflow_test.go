package capacity

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_DeprovisionTRUWorkflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	var a *Activities
	env.OnActivity(a.RemoveTRUs, mock.Anything, mock.Anything).Return(nil)

	input := ProvisionInput{
		Namespace: "test-namespace",
		APSLimit:  2000,
	}

	env.ExecuteWorkflow(DeprovisionTRUWorkflow, input)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	env.AssertExpectations(t)
}
