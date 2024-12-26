package databricks

import "sync"

type Cluster struct {
	mu          sync.Mutex
	Profile     string
	ClusterID   string
	ClusterName string
}

func NewCluster(profile string, clusterID string, clusterName string) *Cluster {
	return &Cluster{sync.Mutex{}, profile, clusterID, clusterName}
}

func (c *Cluster) GetProfile() *string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return &c.Profile
}

func (c *Cluster) GetClusterID() *string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return &c.ClusterID
}

func (c *Cluster) GetClusterName() *string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return &c.ClusterName
}
