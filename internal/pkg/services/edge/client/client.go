// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package client

import (
	"context"

	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-sdk-go/services/edge"
)

// APIClient is an interface that consolidates all client functionality to allow for mocking of the API client during testing.
type APIClient interface {
	CreateInstance(ctx context.Context, projectId, region string) edge.ApiCreateInstanceRequest
	DeleteInstance(ctx context.Context, projectId, region, instanceId string) edge.ApiDeleteInstanceRequest
	DeleteInstanceByName(ctx context.Context, projectId, region, instanceName string) edge.ApiDeleteInstanceByNameRequest
	GetInstance(ctx context.Context, projectId, region, instanceId string) edge.ApiGetInstanceRequest
	GetInstanceByName(ctx context.Context, projectId, region, instanceName string) edge.ApiGetInstanceByNameRequest
	ListInstances(ctx context.Context, projectId, region string) edge.ApiListInstancesRequest
	UpdateInstance(ctx context.Context, projectId, region, instanceId string) edge.ApiUpdateInstanceRequest
	UpdateInstanceByName(ctx context.Context, projectId, region, instanceName string) edge.ApiUpdateInstanceByNameRequest
	GetKubeconfigByInstanceId(ctx context.Context, projectId, region, instanceId string) edge.ApiGetKubeconfigByInstanceIdRequest
	GetKubeconfigByInstanceName(ctx context.Context, projectId, region, instanceName string) edge.ApiGetKubeconfigByInstanceNameRequest
	GetTokenByInstanceId(ctx context.Context, projectId, region, instanceId string) edge.ApiGetTokenByInstanceIdRequest
	GetTokenByInstanceName(ctx context.Context, projectId, region, instanceName string) edge.ApiGetTokenByInstanceNameRequest
	ListPlansProject(ctx context.Context, projectId string) edge.ApiListPlansProjectRequest
}

// ConfigureClient configures and returns a new API client for the Edge service.
func ConfigureClient(p *print.Printer, cliVersion string) (APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.EdgeCustomEndpointKey), false, edge.NewAPIClient)
}
