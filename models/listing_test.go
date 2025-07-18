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
				PostTown:          "London",
				ShortenedPostCode: "W1",
				Region:            "South East",
				PropertyType:      "apartment",
				PriceInCents:      10000000,
			},
			wantErr: false,
		},
		{
			name: "missing post town",
			listing: &Listing{
				ShortenedPostCode: "W1",
				Region:            "South East",
				PropertyType:      "apartment",
				PriceInCents:      10000000,
			},
			wantErr: true,
			errMsg:  "post town is required",
		},
		{
			name: "missing shortened post code",
			listing: &Listing{
				PostTown:     "London",
				Region:       "South East",
				PropertyType: "apartment",
				PriceInCents: 10000000,
			},
			wantErr: true,
			errMsg:  "shortened post code is required",
		},
		{
			name: "missing region",
			listing: &Listing{
				PostTown:          "London",
				ShortenedPostCode: "W1",
				PropertyType:      "apartment",
				PriceInCents:      10000000,
			},
			wantErr: true,
			errMsg:  "region is required",
		},
		{
			name: "missing property type",
			listing: &Listing{
				PostTown:          "London",
				ShortenedPostCode: "W1",
				Region:            "South East",
				PriceInCents:      10000000,
			},
			wantErr: true,
			errMsg:  "property type is required",
		},
		{
			name: "invalid price",
			listing: &Listing{
				PostTown:          "London",
				ShortenedPostCode: "W1",
				Region:            "South East",
				PropertyType:      "apartment",
				PriceInCents:      0,
			},
			wantErr: true,
			errMsg:  "price must be greater than 0",
		},
		{
			name: "negative price",
			listing: &Listing{
				PostTown:          "London",
				ShortenedPostCode: "W1",
				Region:            "South East",
				PropertyType:      "apartment",
				PriceInCents:      -1000,
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
		PostTown:          "London",
		ShortenedPostCode: "W1",
		Region:            "South East",
		PropertyType:      "apartment",
		PriceInCents:      10000000,
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
			PostTown:          "London",
			ShortenedPostCode: "W1",
			Region:            "South East",
			PropertyType:      "apartment",
			PriceInCents:      10000000,
		},
		{
			PostTown:          "Manchester",
			ShortenedPostCode: "M1",
			Region:            "North West",
			PropertyType:      "house",
			PriceInCents:      20000000,
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
		PostTown:          "London",
		ShortenedPostCode: "W1",
		Region:            "South East",
		PropertyType:      "apartment",
		PriceInCents:      10000000,
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
				ID:                listing.ID,
				PostTown:          "Manchester",
				ShortenedPostCode: "M1",
				Region:            "North West",
				PropertyType:      "house",
				PriceInCents:      20000000,
			},
			wantErr: false,
		},
		{
			name: "non-existing listing",
			listing: &Listing{
				ID:                999,
				PostTown:          "London",
				ShortenedPostCode: "W1",
				Region:            "South East",
				PropertyType:      "apartment",
				PriceInCents:      10000000,
			},
			wantErr: true,
		},
		{
			name: "invalid update - missing post town",
			listing: &Listing{
				ID:                listing.ID,
				ShortenedPostCode: "W1",
				Region:            "South East",
				PropertyType:      "apartment",
				PriceInCents:      10000000,
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
				// Verify the update
				updated, err := repo.GetByID(context.Background(), tt.listing.ID)
				assert.NoError(t, err)
				assert.Equal(t, tt.listing.PostTown, updated.PostTown)
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
		PostTown:          "London",
		ShortenedPostCode: "W1",
		Region:            "South East",
		PropertyType:      "apartment",
		PriceInCents:      10000000,
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
				// Verify deletion
				_, err := repo.GetByID(context.Background(), tt.id)
				assert.Error(t, err)
			}
		})
	}
}

func TestListingRepository_GetByRegion(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create test listings
	listings := []*Listing{
		{
			PostTown:          "London",
			ShortenedPostCode: "W1",
			Region:            "South East",
			PropertyType:      "apartment",
			PriceInCents:      10000000,
		},
		{
			PostTown:          "Manchester",
			ShortenedPostCode: "M1",
			Region:            "North West",
			PropertyType:      "house",
			PriceInCents:      20000000,
		},
		{
			PostTown:          "Birmingham",
			ShortenedPostCode: "B1",
			Region:            "South East",
			PropertyType:      "apartment",
			PriceInCents:      15000000,
		},
	}

	for _, listing := range listings {
		err := repo.Create(context.Background(), listing)
		require.NoError(t, err)
	}

	tests := []struct {
		name     string
		region   string
		expected int
	}{
		{
			name:     "South East region",
			region:   "South East",
			expected: 2,
		},
		{
			name:     "North West region",
			region:   "North West",
			expected: 1,
		},
		{
			name:     "non-existing region",
			region:   "Non Existing",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByRegion(context.Background(), tt.region)
			assert.NoError(t, err)
			assert.Len(t, result, tt.expected)
		})
	}
}

func TestListingRepository_GetByPropertyType(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create test listings
	listings := []*Listing{
		{
			PostTown:          "London",
			ShortenedPostCode: "W1",
			Region:            "South East",
			PropertyType:      "apartment",
			PriceInCents:      10000000,
		},
		{
			PostTown:          "Manchester",
			ShortenedPostCode: "M1",
			Region:            "North West",
			PropertyType:      "house",
			PriceInCents:      20000000,
		},
		{
			PostTown:          "Birmingham",
			ShortenedPostCode: "B1",
			Region:            "South East",
			PropertyType:      "apartment",
			PriceInCents:      15000000,
		},
	}

	for _, listing := range listings {
		err := repo.Create(context.Background(), listing)
		require.NoError(t, err)
	}

	tests := []struct {
		name         string
		propertyType string
		expected     int
	}{
		{
			name:         "apartment type",
			propertyType: "apartment",
			expected:     2,
		},
		{
			name:         "house type",
			propertyType: "house",
			expected:     1,
		},
		{
			name:         "non-existing type",
			propertyType: "non-existing",
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByPropertyType(context.Background(), tt.propertyType)
			assert.NoError(t, err)
			assert.Len(t, result, tt.expected)
		})
	}
}

func TestListingRepository_GetFeatured(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create test listings
	listings := []*Listing{
		{
			PostTown:          "London",
			ShortenedPostCode: "W1",
			Region:            "South East",
			PropertyType:      "apartment",
			PriceInCents:      10000000,
			IsFeatured:        true,
		},
		{
			PostTown:          "Manchester",
			ShortenedPostCode: "M1",
			Region:            "North West",
			PropertyType:      "house",
			PriceInCents:      20000000,
			IsFeatured:        false,
		},
		{
			PostTown:          "Birmingham",
			ShortenedPostCode: "B1",
			Region:            "South East",
			PropertyType:      "apartment",
			PriceInCents:      15000000,
			IsFeatured:        true,
		},
	}

	for _, listing := range listings {
		err := repo.Create(context.Background(), listing)
		require.NoError(t, err)
	}

	result, err := repo.GetFeatured(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	for _, listing := range result {
		assert.True(t, listing.IsFeatured)
	}
}

func TestListingRepository_SearchByCity(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create test listings
	listings := []*Listing{
		{
			PostTown:          "London",
			ShortenedPostCode: "W1",
			Region:            "South East",
			PropertyType:      "apartment",
			PriceInCents:      10000000,
		},
		{
			PostTown:          "Manchester",
			ShortenedPostCode: "M1",
			Region:            "North West",
			PropertyType:      "house",
			PriceInCents:      20000000,
		},
		{
			PostTown:          "Birmingham",
			ShortenedPostCode: "B1",
			Region:            "South East",
			PropertyType:      "apartment",
			PriceInCents:      15000000,
		},
	}

	for _, listing := range listings {
		err := repo.Create(context.Background(), listing)
		require.NoError(t, err)
	}

	tests := []struct {
		name     string
		city     string
		expected int
	}{
		{
			name:     "exact match",
			city:     "London",
			expected: 1,
		},
		{
			name:     "case insensitive",
			city:     "london",
			expected: 1,
		},
		{
			name:     "partial match",
			city:     "Lon",
			expected: 1,
		},
		{
			name:     "non-existing city",
			city:     "NonExisting",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.SearchByCity(context.Background(), tt.city)
			assert.NoError(t, err)
			assert.Len(t, result, tt.expected)
		})
	}
}

func TestListingRepository_GetByPriceRange(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create test listings
	listings := []*Listing{
		{
			PostTown:          "London",
			ShortenedPostCode: "W1",
			Region:            "South East",
			PropertyType:      "apartment",
			PriceInCents:      10000000,
		},
		{
			PostTown:          "Manchester",
			ShortenedPostCode: "M1",
			Region:            "North West",
			PropertyType:      "house",
			PriceInCents:      20000000,
		},
		{
			PostTown:          "Birmingham",
			ShortenedPostCode: "B1",
			Region:            "South East",
			PropertyType:      "apartment",
			PriceInCents:      15000000,
		},
	}

	for _, listing := range listings {
		err := repo.Create(context.Background(), listing)
		require.NoError(t, err)
	}

	tests := []struct {
		name     string
		minPrice int64
		maxPrice int64
		expected int
	}{
		{
			name:     "range 10M-20M",
			minPrice: 10000000,
			maxPrice: 20000000,
			expected: 3,
		},
		{
			name:     "range 10M-15M",
			minPrice: 10000000,
			maxPrice: 15000000,
			expected: 2,
		},
		{
			name:     "range 5M-10M",
			minPrice: 5000000,
			maxPrice: 10000000,
			expected: 1,
		},
		{
			name:     "no matches",
			minPrice: 5000000,
			maxPrice: 9000000,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByPriceRange(context.Background(), tt.minPrice, tt.maxPrice)
			assert.NoError(t, err)
			assert.Len(t, result, tt.expected)
		})
	}
}

func TestListingRepository_GetByBedroomRange(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create test listings
	listings := []*Listing{
		{
			PostTown:          "London",
			ShortenedPostCode: "W1",
			Region:            "South East",
			PropertyType:      "apartment",
			Bedrooms:          1,
			PriceInCents:      10000000,
		},
		{
			PostTown:          "Manchester",
			ShortenedPostCode: "M1",
			Region:            "North West",
			PropertyType:      "house",
			Bedrooms:          3,
			PriceInCents:      20000000,
		},
		{
			PostTown:          "Birmingham",
			ShortenedPostCode: "B1",
			Region:            "South East",
			PropertyType:      "apartment",
			Bedrooms:          2,
			PriceInCents:      15000000,
		},
	}

	for _, listing := range listings {
		err := repo.Create(context.Background(), listing)
		require.NoError(t, err)
	}

	tests := []struct {
		name        string
		minBedrooms int
		maxBedrooms int
		expected    int
	}{
		{
			name:        "range 1-3",
			minBedrooms: 1,
			maxBedrooms: 3,
			expected:    3,
		},
		{
			name:        "range 1-2",
			minBedrooms: 1,
			maxBedrooms: 2,
			expected:    2,
		},
		{
			name:        "range 2-3",
			minBedrooms: 2,
			maxBedrooms: 3,
			expected:    2,
		},
		{
			name:        "no matches",
			minBedrooms: 4,
			maxBedrooms: 5,
			expected:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByBedroomRange(context.Background(), tt.minBedrooms, tt.maxBedrooms)
			assert.NoError(t, err)
			assert.Len(t, result, tt.expected)
		})
	}
}

func TestListingRepository_GetByBathroomRange(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create test listings
	listings := []*Listing{
		{
			PostTown:          "London",
			ShortenedPostCode: "W1",
			Region:            "South East",
			PropertyType:      "apartment",
			Bathrooms:         1,
			PriceInCents:      10000000,
		},
		{
			PostTown:          "Manchester",
			ShortenedPostCode: "M1",
			Region:            "North West",
			PropertyType:      "house",
			Bathrooms:         3,
			PriceInCents:      20000000,
		},
		{
			PostTown:          "Birmingham",
			ShortenedPostCode: "B1",
			Region:            "South East",
			PropertyType:      "apartment",
			Bathrooms:         2,
			PriceInCents:      15000000,
		},
	}

	for _, listing := range listings {
		err := repo.Create(context.Background(), listing)
		require.NoError(t, err)
	}

	tests := []struct {
		name         string
		minBathrooms int
		maxBathrooms int
		expected     int
	}{
		{
			name:         "range 1-3",
			minBathrooms: 1,
			maxBathrooms: 3,
			expected:     3,
		},
		{
			name:         "range 1-2",
			minBathrooms: 1,
			maxBathrooms: 2,
			expected:     2,
		},
		{
			name:         "range 2-3",
			minBathrooms: 2,
			maxBathrooms: 3,
			expected:     2,
		},
		{
			name:         "no matches",
			minBathrooms: 4,
			maxBathrooms: 5,
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByBathroomRange(context.Background(), tt.minBathrooms, tt.maxBathrooms)
			assert.NoError(t, err)
			assert.Len(t, result, tt.expected)
		})
	}
}

func TestStringPtr(t *testing.T) {
	str := "test string"
	ptr := stringPtr(str)
	assert.NotNil(t, ptr)
	assert.Equal(t, str, *ptr)
}

func TestListingRepository_Concurrency(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Test concurrent reads and writes
	done := make(chan bool)

	// Start multiple goroutines for concurrent operations
	for i := 0; i < 10; i++ {
		go func(id int) {
			// Create a listing
			listing := &Listing{
				PostTown:          "London",
				ShortenedPostCode: "W1",
				Region:            "South East",
				PropertyType:      "apartment",
				PriceInCents:      10000000,
			}

			err := repo.Create(context.Background(), listing)
			assert.NoError(t, err)

			// Read the listing
			_, err = repo.GetByID(context.Background(), listing.ID)
			assert.NoError(t, err)

			// Update the listing
			listing.PostTown = "Manchester"
			err = repo.Update(context.Background(), listing)
			assert.NoError(t, err)

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all listings were created
	listings, err := repo.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, listings, 10)
}

func TestListingRepository_UpdatePreservesMadeVisibleAt(t *testing.T) {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Create a listing with a specific MadeVisibleAt
	originalTime := "2023-01-01T00:00:00Z"
	listing := &Listing{
		PostTown:          "London",
		ShortenedPostCode: "W1",
		Region:            "South East",
		PropertyType:      "apartment",
		PriceInCents:      10000000,
		MadeVisibleAt:     &originalTime,
	}

	err := repo.Create(context.Background(), listing)
	require.NoError(t, err)

	// Update the listing without MadeVisibleAt
	updateListing := &Listing{
		ID:                listing.ID,
		PostTown:          "Manchester",
		ShortenedPostCode: "M1",
		Region:            "North West",
		PropertyType:      "house",
		PriceInCents:      20000000,
		// MadeVisibleAt is nil
	}

	err = repo.Update(context.Background(), updateListing)
	assert.NoError(t, err)

	// Verify MadeVisibleAt was preserved
	updated, err := repo.GetByID(context.Background(), listing.ID)
	assert.NoError(t, err)
	assert.Equal(t, originalTime, *updated.MadeVisibleAt)
	assert.Equal(t, "Manchester", updated.PostTown)
}
