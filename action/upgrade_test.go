package action

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/deislabs/cnab-go/claim"
	"github.com/deislabs/cnab-go/driver"

	"github.com/stretchr/testify/assert"
)

// makes sure Upgrade implements Action interface
var _ Action = &Upgrade{}

func TestUpgrade_Run(t *testing.T) {
	out := ioutil.Discard

	c := newClaim()
	upgr := &Upgrade{Driver: &mockDriver{
		shouldHandle: true,
		Result: driver.OperationResult{
			Outputs: map[string]string{
				"/tmp/some/path": "SOME CONTENT",
			},
		},
		Error: nil,
	}}
	err := upgr.Run(c, mockSet, out)
	assert.NoError(t, err)
	assert.NotEqual(t, c.Created, c.Modified, "Claim was not updated with modified time stamp during upgrade action")
	assert.Equal(t, claim.ActionUpgrade, c.Result.Action)
	assert.Equal(t, claim.StatusSuccess, c.Result.Status)
	assert.Equal(t, map[string]string{"some-output": "SOME CONTENT"}, c.Outputs)

	c = newClaim()
	upgr = &Upgrade{Driver: &mockDriver{
		Error:        errors.New("I always fail"),
		shouldHandle: false,
	}}
	err = upgr.Run(c, mockSet, out)
	assert.Error(t, err)
	assert.Empty(t, c.Outputs)

	c = newClaim()
	upgr = &Upgrade{Driver: &mockDriver{
		Error:        errors.New("I always fail"),
		shouldHandle: true,
	}}
	err = upgr.Run(c, mockSet, out)
	assert.Error(t, err)
	assert.NotEmpty(t, c.Result.Message, "Expected error message in claim result message")
	assert.Equal(t, claim.ActionUpgrade, c.Result.Action)
	assert.Equal(t, claim.StatusFailure, c.Result.Status)
	assert.Empty(t, c.Outputs)
}
