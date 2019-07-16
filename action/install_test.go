package action

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/deislabs/cnab-go/claim"
	"github.com/deislabs/cnab-go/driver"

	"github.com/stretchr/testify/assert"
)

// makes sure Install implements Action interface
var _ Action = &Install{}

func TestInstall_Run(t *testing.T) {
	out := ioutil.Discard

	c := newClaim()
	inst := &Install{Driver: &driver.DebugDriver{}}
	assert.NoError(t, inst.Run(c, mockSet, out))

	c = newClaim()
	inst = &Install{Driver: &mockDriver{Error: errors.New("I always fail")}}
	assert.Error(t, inst.Run(c, mockSet, out))

	c = newClaim()
	inst = &Install{Driver: &mockDriver{shouldHandle: true, Error: errors.New("I always fail")}}
	assert.Error(t, inst.Run(c, mockSet, out))

	c = newClaim()
	inst = &Install{Driver: &mockDriver{
		shouldHandle: true,
		Result: driver.OperationResult{
			Outputs: map[string]string{
				"/tmp/some/path": "SOME CONTENT",
			},
		},
		Error: nil,
	}}
	assert.NoError(t, inst.Run(c, mockSet, out))
	assert.Equal(t, claim.StatusSuccess, c.Result.Status)
	assert.Equal(t, "install", c.Result.Action)
	assert.Equal(t, map[string]string{"some-output": "SOME CONTENT"}, c.Outputs)
}
