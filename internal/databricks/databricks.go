package databricks

import (
	"context"

	"github.com/databricks/databricks-sdk-go"
	"github.com/databricks/databricks-sdk-go/service/compute"
)

// GetClusterNames fetches all cluster names from the specified profile.
func GetClusterNames(profile string) ([]string, error) {
	w := databricks.Must(databricks.NewWorkspaceClient(&databricks.Config{
		Profile: profile,
	}))

	all, err := w.Clusters.ListAll(context.Background(), compute.ListClustersRequest{})
	if err != nil {
		return nil, err
	}

	var clusterNames []string
	for _, c := range all {
		clusterNames = append(clusterNames, c.ClusterName)
	}

	return clusterNames, nil
}

