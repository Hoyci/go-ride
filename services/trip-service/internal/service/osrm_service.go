package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	tripTypes "go-ride/services/trip-service/pkg/types"
	"go-ride/shared/types"
	"io"
	"net/http"
)

type OSRMService struct{}

func NewOSRMService() *OSRMService {
	return &OSRMService{}
}

func (s *OSRMService) GetRoute(
	ctx context.Context,
	pickup, destination *types.Coordinate,
) (*tripTypes.OSRMApiResponse, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("failed to fetch from OSRM API: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response: %w", err)
	}

	var routeRes tripTypes.OSRMApiResponse
	if err := json.Unmarshal(body, &routeRes); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &routeRes, nil
}
