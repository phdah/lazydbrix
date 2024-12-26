package tui

import (
	"log"
	"sync"
)

type ClusterSelection struct {
	mu          sync.Mutex
	Profile     string
	ClusterID   string
	ClusterName string
}

func NewClusterSelection() *ClusterSelection {
	return &ClusterSelection{}
}

func (cs *ClusterSelection) GetProfile() *string {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return &cs.Profile
}

func (cs *ClusterSelection) GetClusterID() *string {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return &cs.ClusterID
}

func (cs *ClusterSelection) GetClusterName() *string {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return &cs.ClusterName
}

func (cs *ClusterSelection) SetSelection(profile string, clusterID string, clusterName string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	log.Printf("Cluster selected: %s (ID: %s, Profile: %s)", clusterName, clusterID, profile)
	cs.Profile = profile
	cs.ClusterID = clusterID
	cs.ClusterName = clusterName
}
