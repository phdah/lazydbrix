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

// GetClusterNames fetches all cluster names and IDs from the specified profile.
func GetClusterNames(profile string) (*orderedmap.OrderedMap[string, string], *orderedmap.OrderedMap[string, string], error) {
    w := databricks.Must(databricks.NewWorkspaceClient(&databricks.Config{
        Profile: profile,
    }))

    all, err := w.Clusters.ListAll(context.Background(), compute.ListClustersRequest{})
    if err != nil {
        return nil, nil, err
    }

    nameToIDMap := orderedmap.NewOrderedMap[string, string]()
    idToNameMap := orderedmap.NewOrderedMap[string, string]()
    for _, c := range all {
        if strings.HasPrefix(c.ClusterName, "job-") {
            continue
        }
        nameToIDMap.Set(c.ClusterName, c.ClusterId)
        idToNameMap.Set(c.ClusterId, c.ClusterName)
    }

    return nameToIDMap, idToNameMap, nil
}

func GetAllEnvClusters(mu *sync.Mutex, profiles []string) (map[string]*orderedmap.OrderedMap[string, string]){
    profileClusters := make(map[string]*orderedmap.OrderedMap[string, string])
    var wg sync.WaitGroup
    for _, profile := range profiles {
        wg.Add(1)
        go func(profile string) {
            defer wg.Done()
            nameToIDMap, _, err := GetClusterNames(profile)
            if err != nil {
                log.Printf("Failed to fetch clusters for profile %s: %v", profile, err)
            }
            mu.Lock()
            profileClusters[profile] = nameToIDMap
            mu.Unlock()
        }(profile)
    }
    wg.Wait()

    return profileClusters
}

// GetClusterDetails fetches detailed information about a specific cluster.
func GetClusterDetails(profile *string, clusterID string) (*ClusterInfo, error) {
    w := databricks.Must(databricks.NewWorkspaceClient(&databricks.Config{
        Profile: *profile,
    }))

    details, err := w.Clusters.Get(context.Background(), compute.GetClusterRequest{ClusterId: clusterID})
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
func FormatClusterDetails(details *ClusterInfo) (string) {
    jsonData, err := json.MarshalIndent(details, "", "    ")
    if err != nil {
        log.Fatalf("Failed to marshal cluster details: %v", err)
    }

    return string(jsonData)
}

