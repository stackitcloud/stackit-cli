package utils

import (
	"context"
	"fmt"
	"strings"

	sdkPostgreSQL "github.com/stackitcloud/stackit-sdk-go/services/postgresql"
)

type PostgreSQLClient interface {
	CreateInstance(ctx context.Context, projectId string) sdkPostgreSQL.ApiCreateInstanceRequest
	UpdateInstance(ctx context.Context, projectId, instanceId string) sdkPostgreSQL.ApiUpdateInstanceRequest
	GetOfferingsExecute(ctx context.Context, projectId string) (*sdkPostgreSQL.OfferingList, error)
}

func LoadPlanId(ctx context.Context, client PostgreSQLClient, projectId, planName, version string) (*string, error) {
	res, err := client.GetOfferingsExecute(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("get PostgreSQL offerings: %w", err)
	}

	availableVersions := ""
	availablePlanNames := ""
	isValidVersion := false
	for _, offer := range *res.Offerings {
		if !strings.EqualFold(*offer.Version, version) {
			availableVersions = fmt.Sprintf("%s\n- %s", availableVersions, *offer.Version)
			continue
		}
		isValidVersion = true

		for _, plan := range *offer.Plans {
			if plan.Name == nil {
				continue
			}
			if strings.EqualFold(*plan.Name, planName) && plan.Id != nil {
				return plan.Id, nil
			}
			availablePlanNames = fmt.Sprintf("%s\n- %s", availablePlanNames, *plan.Name)
		}
	}

	if !isValidVersion {
		return nil, fmt.Errorf("find version '%s', available versions are: %s", version, availableVersions)
	}
	return nil, fmt.Errorf("find plan_name '%s' for version %s, available names are: %s", planName, version, availablePlanNames)
}
