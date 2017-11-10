package output

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/olekukonko/tablewriter"
)

// ClusterServicePlanHeaders sets the appropriate headers on t for displaying
// ClusterServicePlans in t
func ClusterServicePlanHeaders(table *tablewriter.Table) {
	table.SetHeader([]string{
		"Name",
		"Description",
		"UUID",
		"Class Name",
		"Class UUID",
	})
}

// AppendClusterServicePlan appends plan to t by calling t.Append.
// Ensure that you've called ClusterServicePlanHeaders on t before you call this function
func AppendClusterServicePlan(table *tablewriter.Table, plan *v1beta1.ClusterServicePlan) {
	table.Append([]string{
		plan.Spec.ExternalName,
		plan.Spec.Description,
		plan.Name,
		plan.Spec.ClusterServiceBrokerName,
		plan.Spec.ClusterServiceClassRef.Name,
	})
}
