package plan

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/cmd/svcat/command"
	"github.com/Azure/service-catalog-cli/cmd/svcat/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/spf13/cobra"
)

type describeCmd struct {
	*command.Context
	traverse     bool
	lookupByUUID bool
	uuid         string
	name         string
}

// NewDescribeCmd builds a "svcat describe plan" command
func NewDescribeCmd(cxt *command.Context) *cobra.Command {
	describeCmd := &describeCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:     "plan NAME",
		Aliases: []string{"plans", "pl"},
		Short:   "Show details of a specific plan",
		Example: `
  svcat describe plan standard800
  svcat describe plan --uuid 08e4b43a-36bc-447e-a81f-8202b13e339c
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return describeCmd.run(args)
		},
	}
	cmd.Flags().BoolVarP(
		&describeCmd.traverse,
		"traverse",
		"t",
		false,
		"Whether or not to traverse from plan -> class -> broker",
	)
	cmd.Flags().BoolVarP(
		&describeCmd.lookupByUUID,
		"uuid",
		"u",
		false,
		"Whether or not to get the class by UUID (the default is by name)",
	)
	return cmd
}

func (c *describeCmd) run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("name or uuid is required")
	}

	if c.lookupByUUID {
		c.uuid = args[0]
	} else {
		c.name = args[0]
	}

	return c.describe()
}

func (c *describeCmd) describe() error {
	var plan *v1beta1.ClusterServicePlan
	var err error
	if c.lookupByUUID {
		plan, err = c.App.RetrievePlanByID(c.uuid)
	} else {
		plan, err = c.App.RetrievePlanByName(c.name)
	}
	if err != nil {
		return err
	}

	// Retrieve the class as well because plans don't have the external class name
	class, err := c.App.RetrieveClassByPlan(plan)
	if err != nil {
		return err
	}

	output.WritePlanDetails(c.Output, plan, class)

	instances, err := c.App.RetrieveInstancesByPlan(plan)
	if err != nil {
		return err
	}
	output.WriteAssociatedInstances(c.Output, instances)

	if c.traverse {
		broker, err := c.App.RetrieveBrokerByClass(class)
		if err != nil {
			return err
		}
		output.WriteParentClass(c.Output, class)
		output.WriteParentBroker(c.Output, broker)
	}

	return nil
}
