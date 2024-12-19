package databricks

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/databricks/databricks-sdk-go"
	"github.com/databricks/databricks-sdk-go/service/compute"
	"github.com/elliotchance/orderedmap/v2"
)

// ClusterInfo contains only the desired fields from the cluster details.
type ClusterInfo struct {
	Profile                string
	AutoterminationMinutes int               `json:"autotermination_minutes"`
	ClusterID              string            `json:"cluster_id"`
	ClusterName            string            `json:"cluster_name"`
	DataSecurityMode       string            `json:"data_security_mode"`
	DriverNodeTypeID       string            `json:"driver_node_type_id"`
	SparkEnvVars           map[string]string `json:"spark_env_vars"`
	StartTime              int64             `json:"start_time"`
	State                  string            `json:"state"`
}

// Struct
type DatabricksConnection struct {
	// Public
	ProfileClusters map[string]*orderedmap.OrderedMap[string, string]

	// Private
	profiles   []string
	workspaces map[string]*databricks.WorkspaceClient
	mu         *sync.Mutex
}

// Constructor
func NewDatabricksConnection(profiles []string) *DatabricksConnection {
	return &DatabricksConnection{
		profiles:        profiles,
		workspaces:      make(map[string]*databricks.WorkspaceClient),
		ProfileClusters: make(map[string]*orderedmap.OrderedMap[string, string]),
		mu:              &sync.Mutex{},
	}
}

/*
Setters
*/

// Set workspaces
func (dc *DatabricksConnection) SetWorkspaces() {
	for _, profile := range dc.profiles {
		dc.workspaces[profile] = databricks.Must(databricks.NewWorkspaceClient(&databricks.Config{
			Profile: profile,
		}))
	}
}

// Set all cluster names and IDs
func (dc *DatabricksConnection) SetClusters() {
	var wg sync.WaitGroup
	for _, profile := range dc.profiles {
		wg.Add(1)
		go func(profile string) {
			defer wg.Done()

			all, err := dc.workspaces[profile].Clusters.ListAll(context.Background(), compute.ListClustersRequest{})
			if err != nil {
				log.Panicf("Failed to fetch clusters using profile '%s': %v\n\n[Suggestion: check your VPN]", profile, err)
			}

			nameToIDMap := orderedmap.NewOrderedMap[string, string]()
			for _, c := range all {
				if strings.HasPrefix(c.ClusterName, "job-") {
					continue
				}
				nameToIDMap.Set(c.ClusterName, c.ClusterId)
			}

			dc.mu.Lock()
			dc.ProfileClusters[profile] = nameToIDMap
			dc.mu.Unlock()
		}(profile)
	}
	wg.Wait()
}

// GetClusterDetails fetches detailed information about a specific cluster.
func (dc *DatabricksConnection) GetClusterDetails(profile *string, clusterID string) (*ClusterInfo, error) {
	details, err := dc.workspaces[*profile].Clusters.Get(context.Background(), compute.GetClusterRequest{ClusterId: clusterID})
	if err != nil {
		return &ClusterInfo{}, err // this is not pickung up correct profile
	}

	clusterInfo := &ClusterInfo{
		Profile:                *profile,
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
func (ci *ClusterInfo) Format() string {
	jsonData, err := json.MarshalIndent(ci, "", "    ")
	if err != nil {
		log.Fatalf("Failed to marshal cluster details: %v", err)
	}

	return string(jsonData)
}
