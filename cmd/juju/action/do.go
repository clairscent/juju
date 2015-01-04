// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package action

import (
	"fmt"
	"regexp"

	"github.com/juju/cmd"
	"github.com/juju/errors"
	"github.com/juju/names"
	yaml "gopkg.in/yaml.v1"
	"launchpad.net/gnuflag"

	"github.com/juju/juju/apiserver/params"
)

// DoCommand enqueues an Action for running on the given unit with given
// params
type DoCommand struct {
	ActionCommandBase
	unitTag    names.UnitTag
	actionName string
	paramsYAML cmd.FileVar
	out        cmd.Output
}

const doDoc = `
Queue an Action for execution on a given unit, with a given set of params.
Displays the ID of the Action for use with 'juju kill', 'juju status', etc.

Params are validated according to the charm for the unit's service.  The 
valid params can be seen using "juju action defined <service>".  Params must
be in a yaml file which is passed with the --params flag.

Examples:

$ juju do mysql/3 backup 
action: <UUID>

$ juju status <UUID>
result:
  status: success
  file:
    size: 873.2
    units: GB
    name: foo.sql

$ juju do mysql/3 backup --params parameters.yml
...
`

// actionNameRule describes the format an action name must match to be valid.
var actionNameRule = regexp.MustCompile("^[a-z](?:[a-z-]*[a-z])?$")

// SetFlags offers an option for YAML output.
func (c *DoCommand) SetFlags(f *gnuflag.FlagSet) {
	c.out.AddFlags(f, "smart", cmd.DefaultFormatters)
	f.Var(&c.paramsYAML, "params", "path to yaml-formatted params file")
}

func (c *DoCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "do",
		Args:    "<unit> <action name>",
		Purpose: "WIP: queue an action for execution",
		Doc:     doDoc,
	}
}

// Init gets the unit tag, and checks for other correct args.
func (c *DoCommand) Init(args []string) error {
	switch len(args) {
	case 0:
		return errors.New("no unit specified")
	case 1:
		return errors.New("no action specified")
	case 2:
		unitName := args[0]
		if !names.IsValidUnit(unitName) {
			return errors.Errorf("invalid unit name %q", unitName)
		}
		actionName := args[1]
		if valid := actionNameRule.MatchString(actionName); !valid {
			return fmt.Errorf("invalid action name %q", actionName)
		}
		c.unitTag = names.NewUnitTag(unitName)
		c.actionName = actionName
		return nil
	default:
		return cmd.CheckEmpty(args[1:])
	}
}

func (c *DoCommand) Run(ctx *cmd.Context) error {
	api, err := c.NewActionAPIClient()
	if err != nil {
		return err
	}
	defer api.Close()

	actionParams := map[string]interface{}{}

	if c.paramsYAML.Path != "" {
		b, err := c.paramsYAML.Read(ctx)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(b, &actionParams)

		conformantParams, err := conform(actionParams)
		if err != nil {
			return err
		}

		betterParams, ok := conformantParams.(map[string]interface{})
		if !ok {
			return errors.New("params must contain a YAML map with string keys")
		}

		actionParams = betterParams
	}

	actionParam := params.Actions{
		Actions: []params.Action{{
			Receiver:   c.unitTag.String(),
			Name:       c.actionName,
			Parameters: actionParams,
		}},
	}

	results, err := api.Enqueue(actionParam)
	if err != nil {
		return err
	}
	if len(results.Results) != 1 {
		return errors.New("illegal number of results returned")
	}

	result := results.Results[0]

	if result.Error != nil {
		return result.Error
	}

	if result.Action == nil {
		return errors.New("action failed to enqueue")
	}

	tag, err := names.ParseActionTag(result.Action.Tag)
	if err != nil {
		return err
	}

	output := map[string]string{"Action queued with id": tag.Id()}
	return c.out.Write(ctx, output)
}
