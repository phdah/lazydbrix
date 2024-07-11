package databricks

import (
    "context"
    "fmt"
    "log"
    "encoding/json"

    "github.com/databricks/databricks-sdk-go"
    "github.com/databricks/databricks-sdk-go/service/compute"
)

// ClusterInfo contains only the desired fields from the cluster details.
type ClusterInfo struct {
	AutoterminationMinutes int               `json:"autotermination_minutes"`
	ClusterID              string            `json:"cluster_id"`
	ClusterName            string            `json:"cluster_name"`
	DataSecurityMode       string            `json:"data_security_mode"`
	DriverNodeTypeID       string            `json:"driver_node_type_id"`
	SparkEnvVars           map[string]string `json:"spark_env_vars"`
	StartTime              int64             `json:"start_time"`
	State                  string            `json:"state"`
}

// GetClusterNames fetches all cluster names and IDs from the specified profile.
func GetClusterNames(profile string) (map[string]string, map[string]string, error) {
	w := databricks.Must(databricks.NewWorkspaceClient(&databricks.Config{
		Profile: profile,
	}))

	all, err := w.Clusters.ListAll(context.Background(), compute.ListClustersRequest{})
	if err != nil {
		return nil, nil, err
	}

	nameToIDMap := make(map[string]string)
	idToNameMap := make(map[string]string)
	for _, c := range all {
		nameToIDMap[c.ClusterName] = c.ClusterId
		idToNameMap[c.ClusterId] = c.ClusterName
	}

	return nameToIDMap, idToNameMap, nil
}

// GetClusterDetails fetches detailed information about a specific cluster.
func GetClusterDetails(profile, clusterID string) (*ClusterInfo, error) {
	w := databricks.Must(databricks.NewWorkspaceClient(&databricks.Config{
		Profile: profile,
	}))

	details, err := w.Clusters.Get(context.Background(), compute.GetClusterRequest{ClusterId: clusterID})
	if err != nil {
		return &ClusterInfo{}, err
	}

	clusterInfo := &ClusterInfo{
		AutoterminationMinutes: details.AutoterminationMinutes,
		ClusterID:              details.ClusterId,
		ClusterName:            details.ClusterName,
		DataSecurityMode:       string(details.DataSecurityMode),
		DriverNodeTypeID:       details.DataSecurityMode.String(),
		SparkEnvVars:           details.SparkEnvVars,
		StartTime:              details.StartTime,
		State:                  string(details.State),
	}

	return clusterInfo, nil
}

// DisplayClusterDetails formats and displays cluster details as JSON.
func DisplayClusterDetails(details *ClusterInfo) {
	jsonData, err := json.MarshalIndent(details, "", "    ")
	if err != nil {
		log.Fatalf("Failed to marshal cluster details: %v", err)
	}

	fmt.Println(string(jsonData))
}

