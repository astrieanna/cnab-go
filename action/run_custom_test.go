package action

import (
	"io/ioutil"
	"testing"

	"github.com/deislabs/cnab-go/claim"
	"github.com/deislabs/cnab-go/driver"

	"github.com/deislabs/cnab-go/bundle"
	"github.com/stretchr/testify/assert"
)

// makes sure RunCustom implements Action interface
var _ Action = &RunCustom{}

func TestRunCustom(t *testing.T) {
	out := ioutil.Discard

	rc := &RunCustom{
		Driver: &mockDriver{
			shouldHandle: true,
			Result: driver.OperationResult{
				Outputs: map[string]string{
					"/tmp/some/path": "SOME CONTENT",
				},
			},
			Error: nil,
		},
		Action: "test",
	}
	c := newClaim()
	err := rc.Run(c, mockSet, out)
	assert.NoError(t, err)
	assert.Equal(t, claim.StatusSuccess, c.Result.Status)
	assert.Equal(t, "test", c.Result.Action)
	assert.Equal(t, map[string]string{"some-output": "SOME CONTENT"}, c.Outputs)

	// Make sure we don't allow forbidden custom actions
	c = newClaim()
	rc.Action = "install"
	err = rc.Run(c, mockSet, out)
	assert.Error(t, err)
	assert.Empty(t, c.Outputs)

	// Get rid of custom actions, and this should fail
	c = newClaim()
	rc.Action = "test"
	c.Bundle.Actions = map[string]bundle.Action{}
	err = rc.Run(c, mockSet, out)
	assert.Error(t, err, "Unknown action should fail")
	assert.Empty(t, c.Outputs)
}
