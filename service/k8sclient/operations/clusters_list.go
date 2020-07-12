// dbaas-controller
// Copyright (C) 2020 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package operations

import (
	"context"
	"encoding/json"

	pxc "github.com/percona/percona-xtradb-cluster-operator/pkg/apis/pxc/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/percona-platform/dbaas-controller/service/k8sclient/kubectl"
)

// Cluster contains information related to cluster.
type Cluster struct {
	Name   string
	Status string
}

// NewClusterList returns new object of ClusterList
func NewClusterList(kubeCtl *kubectl.KubeCtl) *ClusterList {
	return &ClusterList{
		kubeCtl: kubeCtl,
	}
}

// ClusterList contains all logic related to getting cluster list.
type ClusterList struct {
	kubeCtl *kubectl.KubeCtl
}

// GetClusters returns clusters list.
func (c *ClusterList) GetClusters(ctx context.Context) ([]Cluster, error) {
	perconaXtraDBClusters, err := c.getPerconaXtraDBClusters(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]Cluster, len(perconaXtraDBClusters))
	for i, cluster := range perconaXtraDBClusters {
		val := Cluster{
			Name:   cluster.Name,
			Status: string(cluster.Status.Status),
		}

		res[i] = val
	}
	return res, nil
}

func (c *ClusterList) getPerconaXtraDBClusters(ctx context.Context) ([]*pxc.PerconaXtraDBCluster, error) {
	stdout, err := c.kubeCtl.Get(ctx, clusterKind, "")
	if err != nil {
		return nil, err
	}

	var list meta.List
	if err := json.Unmarshal(stdout, &list); err != nil {
		return nil, err
	}

	res := make([]*pxc.PerconaXtraDBCluster, len(list.Items))
	for _, item := range list.Items {
		var cluster pxc.PerconaXtraDBCluster
		if err := json.Unmarshal(item.Raw, &cluster); err != nil {
			return nil, err
		}
		res = append(res, &cluster)
	}
	return res, nil
}
