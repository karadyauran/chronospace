package services

import (
	"context"
	"fmt"

	"googlemaps.github.io/maps"
)

type MapsService struct {
	client *maps.Client
}

func (s MapsService) ValidateLocation(ctx context.Context, location string) (bool, error) {
	r := &maps.GeocodingRequest{
		Address: location,
	}

	resp, err := s.client.Geocode(ctx, r)
	if err != nil {
		return false, err
	}

	// If we got any results, the location is valid
	if len(resp) > 0 {
		return true, nil
	}

	return false, nil
}

func NewMapsService(apiKey string) *MapsService {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		fmt.Print(err)
	}
	return &MapsService{client: client}
}

func (s *MapsService) SearchPlaces(query string) ([]maps.PlacesSearchResult, error) {
	r := &maps.TextSearchRequest{
		Query: query,
	}

	resp, err := s.client.TextSearch(context.Background(), r)
	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}
