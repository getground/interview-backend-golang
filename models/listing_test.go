package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewListingRepository(t *testing.T) {
	repo := NewListingRepository()
	assert.NotNil(t, repo)

	// Test that sample data was loaded
	listings, err := repo.GetAll(context.Background())
	require.NoError(t, err)
	assert.Greater(t, len(listings), 0)
}

func TestListingRepository_Create(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	tests := []struct {
		name    string
		listing *Listing
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid listing",
			listing: &Listing{
				AddressDetails: AddressDetails{
					City:              "London",
					ShortenedPostcode: "W1",
					Region:            RegionSouthEast,
					Country:           "UK",
				},
				PropertyType: PropertyTypeApartment,
				PriceInCents: 10000000,
			},
			wantErr: false,
		},
		{
			name: "missing city",
			listing: &Listing{
				AddressDetails: AddressDetails{
					ShortenedPostcode: "W1",
					Region:            RegionSouthEast,
					Country:           "UK",
				},
				PropertyType: PropertyTypeApartment,
				PriceInCents: 10000000,
			},
			wantErr: true,
			errMsg:  "city is required",
		},
		{
			name: "missing shortened postcode",
			listing: &Listing{
				AddressDetails: AddressDetails{
					City:    "London",
					Region:  RegionSouthEast,
					Country: "UK",
				},
				PropertyType: PropertyTypeApartment,
				PriceInCents: 10000000,
			},
			wantErr: true,
			errMsg:  "shortened postcode is required",
		},
		{
			name: "missing region",
			listing: &Listing{
				AddressDetails: AddressDetails{
					City:              "London",
					ShortenedPostcode: "W1",
					Country:           "UK",
				},
				PropertyType: PropertyTypeApartment,
				PriceInCents: 10000000,
			},
			wantErr: true,
			errMsg:  "region is required",
		},
		{
			name: "missing property type",
			listing: &Listing{
				AddressDetails: AddressDetails{
					City:              "London",
					ShortenedPostcode: "W1",
					Region:            RegionSouthEast,
					Country:           "UK",
				},
				PriceInCents: 10000000,
			},
			wantErr: true,
			errMsg:  "property type is required",
		},
		{
			name: "invalid price",
			listing: &Listing{
				AddressDetails: AddressDetails{
					City:              "London",
					ShortenedPostcode: "W1",
					Region:            RegionSouthEast,
					Country:           "UK",
				},
				PropertyType: PropertyTypeApartment,
				PriceInCents: 0,
			},
			wantErr: true,
			errMsg:  "price must be greater than 0",
		},
		{
			name: "negative price",
			listing: &Listing{
				AddressDetails: AddressDetails{
					City:              "London",
					ShortenedPostcode: "W1",
					Region:            RegionSouthEast,
					Country:           "UK",
				},
				PropertyType: PropertyTypeApartment,
				PriceInCents: -1000,
			},
			wantErr: true,
			errMsg:  "price must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(context.Background(), tt.listing)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.listing.ID)
				assert.NotNil(t, tt.listing.MadeVisibleAt)
			}
		})
	}
}

func TestListingRepository_GetByID(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create a test listing
	listing := &Listing{
		AddressDetails: AddressDetails{
			City:              "London",
			ShortenedPostcode: "W1",
			Region:            RegionSouthEast,
			Country:           "UK",
		},
		PropertyType: PropertyTypeApartment,
		PriceInCents: 10000000,
	}
	err := repo.Create(context.Background(), listing)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "existing listing",
			id:      listing.ID,
			wantErr: false,
		},
		{
			name:    "non-existing listing",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByID(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.id, result.ID)
			}
		})
	}
}

func TestListingRepository_GetAll(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create multiple test listings
	listings := []*Listing{
		{
			AddressDetails: AddressDetails{
				City:              "London",
				ShortenedPostcode: "W1",
				Region:            RegionSouthEast,
				Country:           "UK",
			},
			PropertyType: PropertyTypeApartment,
			PriceInCents: 10000000,
		},
		{
			AddressDetails: AddressDetails{
				City:              "Manchester",
				ShortenedPostcode: "M1",
				Region:            RegionNorthWest,
				Country:           "UK",
			},
			PropertyType: PropertyTypeDetached,
			PriceInCents: 20000000,
		},
	}

	for _, listing := range listings {
		err := repo.Create(context.Background(), listing)
		require.NoError(t, err)
	}

	result, err := repo.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestListingRepository_Update(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create a test listing
	listing := &Listing{
		AddressDetails: AddressDetails{
			City:              "London",
			ShortenedPostcode: "W1",
			Region:            RegionSouthEast,
			Country:           "UK",
		},
		PropertyType: PropertyTypeApartment,
		PriceInCents: 10000000,
	}
	err := repo.Create(context.Background(), listing)
	require.NoError(t, err)

	tests := []struct {
		name    string
		listing *Listing
		wantErr bool
	}{
		{
			name: "valid update",
			listing: &Listing{
				ID: listing.ID,
				AddressDetails: AddressDetails{
					City:              "Manchester",
					ShortenedPostcode: "M1",
					Region:            RegionNorthWest,
					Country:           "UK",
				},
				PropertyType: PropertyTypeDetached,
				PriceInCents: 20000000,
			},
			wantErr: false,
		},
		{
			name: "non-existing listing",
			listing: &Listing{
				ID: 999,
				AddressDetails: AddressDetails{
					City:              "London",
					ShortenedPostcode: "W1",
					Region:            RegionSouthEast,
					Country:           "UK",
				},
				PropertyType: PropertyTypeApartment,
				PriceInCents: 10000000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(context.Background(), tt.listing)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListingRepository_Delete(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create a test listing
	listing := &Listing{
		AddressDetails: AddressDetails{
			City:              "London",
			ShortenedPostcode: "W1",
			Region:            RegionSouthEast,
			Country:           "UK",
		},
		PropertyType: PropertyTypeApartment,
		PriceInCents: 10000000,
	}
	err := repo.Create(context.Background(), listing)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "existing listing",
			id:      listing.ID,
			wantErr: false,
		},
		{
			name:    "non-existing listing",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListingRepository_GetByRegion(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create test listings in different regions
	listings := []*Listing{
		{
			AddressDetails: AddressDetails{
				City:              "London",
				ShortenedPostcode: "W1",
				Region:            RegionSouthEast,
				Country:           "UK",
			},
			PropertyType: PropertyTypeApartment,
			PriceInCents: 10000000,
		},
		{
			AddressDetails: AddressDetails{
				City:              "Manchester",
				ShortenedPostcode: "M1",
				Region:            RegionNorthWest,
				Country:           "UK",
			},
			PropertyType: PropertyTypeDetached,
			PriceInCents: 20000000,
		},
		{
			AddressDetails: AddressDetails{
				City:              "Birmingham",
				ShortenedPostcode: "B1",
				Region:            RegionMidlands,
				Country:           "UK",
			},
			PropertyType: PropertyTypeTerraced,
			PriceInCents: 15000000,
		},
	}

	for _, listing := range listings {
		err := repo.Create(context.Background(), listing)
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		region        string
		expectedCount int
	}{
		{
			name:          "South East region",
			region:        string(RegionSouthEast),
			expectedCount: 1,
		},
		{
			name:          "North West region",
			region:        string(RegionNorthWest),
			expectedCount: 1,
		},
		{
			name:          "Non-existing region",
			region:        "Non-existing",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByRegion(context.Background(), tt.region)
			assert.NoError(t, err)
			assert.Len(t, result, tt.expectedCount)
		})
	}
}

func TestListingRepository_GetByPropertyType(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create test listings with different property types
	listings := []*Listing{
		{
			AddressDetails: AddressDetails{
				City:              "London",
				ShortenedPostcode: "W1",
				Region:            RegionSouthEast,
				Country:           "UK",
			},
			PropertyType: PropertyTypeApartment,
			PriceInCents: 10000000,
		},
		{
			AddressDetails: AddressDetails{
				City:              "Manchester",
				ShortenedPostcode: "M1",
				Region:            RegionNorthWest,
				Country:           "UK",
			},
			PropertyType: PropertyTypeDetached,
			PriceInCents: 20000000,
		},
		{
			AddressDetails: AddressDetails{
				City:              "Birmingham",
				ShortenedPostcode: "B1",
				Region:            RegionMidlands,
				Country:           "UK",
			},
			PropertyType: PropertyTypeTerraced,
			PriceInCents: 15000000,
		},
	}

	for _, listing := range listings {
		err := repo.Create(context.Background(), listing)
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		propertyType  string
		expectedCount int
	}{
		{
			name:          "apartment",
			propertyType:  string(PropertyTypeApartment),
			expectedCount: 1,
		},
		{
			name:          "detached",
			propertyType:  string(PropertyTypeDetached),
			expectedCount: 1,
		},
		{
			name:          "non-existing type",
			propertyType:  "mansion",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByPropertyType(context.Background(), tt.propertyType)
			assert.NoError(t, err)
			assert.Len(t, result, tt.expectedCount)
		})
	}
}

func TestListingRepository_SearchByCity(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create test listings in different cities
	listings := []*Listing{
		{
			AddressDetails: AddressDetails{
				City:              "London",
				ShortenedPostcode: "W1",
				Region:            RegionSouthEast,
				Country:           "UK",
			},
			PropertyType: PropertyTypeApartment,
			PriceInCents: 10000000,
		},
		{
			AddressDetails: AddressDetails{
				City:              "Manchester",
				ShortenedPostcode: "M1",
				Region:            RegionNorthWest,
				Country:           "UK",
			},
			PropertyType: PropertyTypeDetached,
			PriceInCents: 20000000,
		},
		{
			AddressDetails: AddressDetails{
				City:              "Birmingham",
				ShortenedPostcode: "B1",
				Region:            RegionMidlands,
				Country:           "UK",
			},
			PropertyType: PropertyTypeTerraced,
			PriceInCents: 15000000,
		},
	}

	for _, listing := range listings {
		err := repo.Create(context.Background(), listing)
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		city          string
		expectedCount int
	}{
		{
			name:          "exact match - London",
			city:          "London",
			expectedCount: 1,
		},
		{
			name:          "partial match - man",
			city:          "man",
			expectedCount: 1,
		},
		{
			name:          "case insensitive - london",
			city:          "london",
			expectedCount: 1,
		},
		{
			name:          "no match",
			city:          "Edinburgh",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.SearchByCity(context.Background(), tt.city)
			assert.NoError(t, err)
			assert.Len(t, result, tt.expectedCount)
		})
	}
}

func TestListingRepository_GetByPriceRange(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create test listings with different prices
	listings := []*Listing{
		{
			AddressDetails: AddressDetails{
				City:              "London",
				ShortenedPostcode: "W1",
				Region:            RegionSouthEast,
				Country:           "UK",
			},
			PropertyType: PropertyTypeApartment,
			PriceInCents: 5000000,
		},
		{
			AddressDetails: AddressDetails{
				City:              "Manchester",
				ShortenedPostcode: "M1",
				Region:            RegionNorthWest,
				Country:           "UK",
			},
			PropertyType: PropertyTypeDetached,
			PriceInCents: 15000000,
		},
		{
			AddressDetails: AddressDetails{
				City:              "Birmingham",
				ShortenedPostcode: "B1",
				Region:            RegionMidlands,
				Country:           "UK",
			},
			PropertyType: PropertyTypeTerraced,
			PriceInCents: 25000000,
		},
	}

	for _, listing := range listings {
		err := repo.Create(context.Background(), listing)
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		minPrice      int64
		maxPrice      int64
		expectedCount int
	}{
		{
			name:          "low range",
			minPrice:      0,
			maxPrice:      10000000,
			expectedCount: 1,
		},
		{
			name:          "middle range",
			minPrice:      10000000,
			maxPrice:      20000000,
			expectedCount: 1,
		},
		{
			name:          "wide range",
			minPrice:      0,
			maxPrice:      30000000,
			expectedCount: 3,
		},
		{
			name:          "no matches",
			minPrice:      50000000,
			maxPrice:      100000000,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByPriceRange(context.Background(), tt.minPrice, tt.maxPrice)
			assert.NoError(t, err)
			assert.Len(t, result, tt.expectedCount)
		})
	}
}

func TestListingRepository_GetByBedroomRange(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create test listings with different bedroom counts
	listings := []*Listing{
		{
			AddressDetails: AddressDetails{
				City:              "London",
				ShortenedPostcode: "W1",
				Region:            RegionSouthEast,
				Country:           "UK",
			},
			PropertyType: PropertyTypeApartment,
			PriceInCents: 10000000,
			Bedrooms:     1,
		},
		{
			AddressDetails: AddressDetails{
				City:              "Manchester",
				ShortenedPostcode: "M1",
				Region:            RegionNorthWest,
				Country:           "UK",
			},
			PropertyType: PropertyTypeDetached,
			PriceInCents: 20000000,
			Bedrooms:     3,
		},
		{
			AddressDetails: AddressDetails{
				City:              "Birmingham",
				ShortenedPostcode: "B1",
				Region:            RegionMidlands,
				Country:           "UK",
			},
			PropertyType: PropertyTypeTerraced,
			PriceInCents: 15000000,
			Bedrooms:     2,
		},
	}

	for _, listing := range listings {
		err := repo.Create(context.Background(), listing)
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		minBedrooms   int
		maxBedrooms   int
		expectedCount int
	}{
		{
			name:          "1 bedroom",
			minBedrooms:   1,
			maxBedrooms:   1,
			expectedCount: 1,
		},
		{
			name:          "2-3 bedrooms",
			minBedrooms:   2,
			maxBedrooms:   3,
			expectedCount: 2,
		},
		{
			name:          "all bedrooms",
			minBedrooms:   0,
			maxBedrooms:   10,
			expectedCount: 3,
		},
		{
			name:          "no matches",
			minBedrooms:   5,
			maxBedrooms:   10,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByBedroomRange(context.Background(), tt.minBedrooms, tt.maxBedrooms)
			assert.NoError(t, err)
			assert.Len(t, result, tt.expectedCount)
		})
	}
}

func TestListingRepository_GetByBathroomRange(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create test listings with different bathroom counts
	listings := []*Listing{
		{
			AddressDetails: AddressDetails{
				City:              "London",
				ShortenedPostcode: "W1",
				Region:            RegionSouthEast,
				Country:           "UK",
			},
			PropertyType: PropertyTypeApartment,
			PriceInCents: 10000000,
			Bathrooms:    1,
		},
		{
			AddressDetails: AddressDetails{
				City:              "Manchester",
				ShortenedPostcode: "M1",
				Region:            RegionNorthWest,
				Country:           "UK",
			},
			PropertyType: PropertyTypeDetached,
			PriceInCents: 20000000,
			Bathrooms:    3,
		},
		{
			AddressDetails: AddressDetails{
				City:              "Birmingham",
				ShortenedPostcode: "B1",
				Region:            RegionMidlands,
				Country:           "UK",
			},
			PropertyType: PropertyTypeTerraced,
			PriceInCents: 15000000,
			Bathrooms:    2,
		},
	}

	for _, listing := range listings {
		err := repo.Create(context.Background(), listing)
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		minBathrooms  int
		maxBathrooms  int
		expectedCount int
	}{
		{
			name:          "1 bathroom",
			minBathrooms:  1,
			maxBathrooms:  1,
			expectedCount: 1,
		},
		{
			name:          "2-3 bathrooms",
			minBathrooms:  2,
			maxBathrooms:  3,
			expectedCount: 2,
		},
		{
			name:          "all bathrooms",
			minBathrooms:  0,
			maxBathrooms:  10,
			expectedCount: 3,
		},
		{
			name:          "no matches",
			minBathrooms:  5,
			maxBathrooms:  10,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByBathroomRange(context.Background(), tt.minBathrooms, tt.maxBathrooms)
			assert.NoError(t, err)
			assert.Len(t, result, tt.expectedCount)
		})
	}
}
