package iface

type Cluster interface {
	GetProfile() *string
	GetClusterID() *string
	GetClusterName() *string
}

type Selector interface {
	Cluster
	SetSelection(profile string, clusterID string, clusterName string)
}
