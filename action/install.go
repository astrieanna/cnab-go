package action

import (
	"io"

	"github.com/deislabs/cnab-go/claim"
	"github.com/deislabs/cnab-go/credentials"
	"github.com/deislabs/cnab-go/driver"
)

// Install describes an installation action
type Install struct {
	Driver driver.Driver // Needs to be more than a string
}

// Run performs an installation and updates the Claim accordingly
func (i *Install) Run(c *claim.Claim, creds credentials.Set, w io.Writer) error {
	invocImage, err := selectInvocationImage(i.Driver, c)
	if err != nil {
		return err
	}

	op, err := opFromClaim(claim.ActionInstall, stateful, c, invocImage, creds, w)
	if err != nil {
		return err
	}

	opResult, err := i.Driver.Run(op)
	// Update outputs in claim
	c.Outputs = map[string]string{}
	for outputName, v := range c.Bundle.Outputs.Fields {
		if opResult.Outputs[v.Path] != "" {
			c.Outputs[outputName] = opResult.Outputs[v.Path]
		}
	}

	if err != nil {
		c.Update(claim.ActionInstall, claim.StatusFailure)
		c.Result.Message = err.Error()
		return err
	}
	c.Update(claim.ActionInstall, claim.StatusSuccess)
	return nil
}
