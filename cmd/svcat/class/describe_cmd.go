package class

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

// NewDescribeCmd builds a "svcat describe class" command
func NewDescribeCmd(cxt *command.Context) *cobra.Command {
	describeCmd := &describeCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:     "class NAME",
		Aliases: []string{"classes", "cl"},
		Short:   "Show details of a specific class",
		Example: `
  svcat describe class azure-mysqldb
  svcat describe class -uuid 997b8372-8dac-40ac-ae65-758b4a5075a5
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
	var class *v1beta1.ClusterServiceClass
	var err error
	if c.lookupByUUID {
		class, err = c.App.RetrieveClassByID(c.uuid)
	} else {
		class, err = c.App.RetrieveClassByName(c.name)
	}
	if err != nil {
		return err
	}

	output.WriteClassDetails(c.Output, class)

	plans, err := c.App.RetrievePlansByClass(class)
	if err != nil {
		return err
	}
	output.WriteAssociatedPlans(c.Output, plans)

	if c.traverse {
		broker, err := c.App.RetrieveBrokerByClass(class)
		if err != nil {
			return err
		}
		output.WriteParentBroker(c.Output, broker)
		output.WriteAssociatedPlans(c.Output, plans)
	}

	return nil
}
