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

// Package logs contains implementation of API for getting logs out of
// Kubernetes cluster workloads.
package logs

import (
	"context"

	controllerv1beta1 "github.com/percona-platform/dbaas-api/gen/controller"
	"github.com/percona-platform/dbaas-controller/service/k8sclient"
	"github.com/pkg/errors"
	"golang.org/x/text/message"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LogsService struct {
	p                *message.Printer
	defaultLogSource LogsSource
	logSources       []LogsSource
}

type LogsSource interface {
	GetLogs(ctx context.Context, client *k8sclient.K8sClient, clusterName string) ([]*controllerv1beta1.Logs, error)
}

type AllLogsSource struct{}

func (a AllLogsSource) GetLogs(ctx context.Context, client *k8sclient.K8sClient, clusterName string) ([]*controllerv1beta1.Logs, error) {
	pods, err := client.GetClusterPods(ctx, clusterName)
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrap(err, "failed to get pods").Error())
	}
	// Every pod has at least one contaier, set cap to that value.
	response := make([]*controllerv1beta1.Logs, 0, len(pods.Items))
	// Get all logs from all pod's containers.
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			logs, err := client.GetLogs(ctx, pod.Name, container.Name)
			if err != nil {
				return nil, status.Error(codes.Internal, errors.Wrap(err, "failed to get logs").Error())
			}
			response = append(response, &controllerv1beta1.Logs{
				Pod:       pod.Name,
				Container: container.Name,
				Logs:      logs,
			})
		}
	}

	// TODO Get all events from all pods.
	return response, nil
}

func NewService(p *message.Printer) *LogsService {
	return &LogsService{
		p:                p,
		defaultLogSource: AllLogsSource{},
		logSources:       []LogsSource{},
	}
}

// GetLogs first tries to get logs and events only from failing pods/containers.
// If no such logs/events are found, it returns all logs and events.
func (s *LogsService) GetLogs(ctx context.Context, req *controllerv1beta1.GetLogsRequest) (*controllerv1beta1.GetLogsResponse, error) {
	client, ok := ctx.Value("k8sclient").(*k8sclient.K8sClient)
	if !ok {
		return nil, status.Error(codes.Internal, "failed to get k8s client")
	}

	response := []*controllerv1beta1.Logs{}
	for _, logSource := range s.logSources {
		logs, err := logSource.GetLogs(ctx, client, req.ClusterName)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to get logs")
		}
		response = append(response, logs...)
	}
	if len(response) == 0 {
		logs, err := s.defaultLogSource.GetLogs(ctx, client, req.ClusterName)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to get logs")
		}
		response = append(response, logs...)
	}

	return &controllerv1beta1.GetLogsResponse{
		Logs: response,
	}, nil
}
