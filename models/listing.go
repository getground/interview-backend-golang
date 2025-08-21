package models

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// Region represents valid UK regions
type Region string

const (
	RegionNorthWest Region = "North West"
	RegionLondon    Region = "London"
	RegionNorthEast Region = "North East"
	RegionSouthWest Region = "South West"
	RegionSouthEast Region = "South East"
	RegionMidlands  Region = "Midlands"
	RegionScotland  Region = "Scotland"
	RegionWales     Region = "Wales"
)

// PropertyType represents valid property types
type PropertyType string

const (
	PropertyTypeApartment    PropertyType = "apartment"
	PropertyTypeDetached     PropertyType = "detached"
	PropertyTypeSemiDetached PropertyType = "semi-detached"
	PropertyTypeTerraced     PropertyType = "terraced"
	PropertyTypeEndTerrace   PropertyType = "end-terrace"
)

// AddressDetails represents the address information for a listing
type AddressDetails struct {
	AddressLine1      string `json:"addressLine1"`
	AddressLine2      string `json:"addressLine2"`
	City              string `json:"city"`
	Postcode          string `json:"postcode"`
	ShortenedPostcode string `json:"shortenedPostcode"`
	Country           string `json:"country"`
	Region            Region `json:"region"`
}

// Photo represents a property photo
type Photo struct {
	OriginalURL  string `json:"originalURL"`
	StandardURL  string `json:"standardURL"`
	ThumbnailURL string `json:"thumbnailURL"`
	MimeType     string `json:"mimeType"`
}

// Listing represents a property listing
type Listing struct {
	ID                         int64          `json:"id"`
	AddressDetails             AddressDetails `json:"addressDetails"`
	Bedrooms                   int            `json:"bedrooms"`
	Bathrooms                  int            `json:"bathrooms"`
	Description                string         `json:"description"`
	GrossYield                 float64        `json:"grossYield"`
	IsCashOnly                 bool           `json:"isCashOnly"`
	IsCompany                  bool           `json:"isCompany"`
	IsNewBuild                 bool           `json:"isNewBuild"`
	IsShareSale                bool           `json:"isShareSale"`
	IsTenanted                 bool           `json:"isTenanted"`
	MadeVisibleAt              *string        `json:"madeVisibleAt"`
	EstimatedDepositInCents    int64          `json:"estimatedDepositInCents"`
	MinimumDepositInCents      int64          `json:"minimumDepositInCents"`
	Photos                     []Photo        `json:"photos"`
	PriceInCents               int64          `json:"priceInCents"`
	PropertyType               PropertyType   `json:"propertyType"`
	MonthlyRentalIncomeInCents int64          `json:"monthlyRentalIncomeInCents"`
	SizeSqFt                   int            `json:"sizeSqFt"`
}

// ListingResponse represents the top-level response structure
type ListingResponse struct {
	Type        string       `json:"type"`
	Listing     *Listing     `json:"listing"`
	Development *Development `json:"development"`
}

// Development represents a property development (can be null)
type Development struct {
	// Add development fields as needed
	// For now, this is a placeholder since the schema shows development as null
}

// ListingRepository interface defines the operations for listing data
type ListingRepository interface {
	Create(ctx context.Context, listing *Listing) error
	GetByID(ctx context.Context, id int64) (*Listing, error)
	GetAll(ctx context.Context) ([]*Listing, error)
	Update(ctx context.Context, listing *Listing) error
	Delete(ctx context.Context, id int64) error
	GetByRegion(ctx context.Context, region string) ([]*Listing, error)
	GetByPropertyType(ctx context.Context, propertyType string) ([]*Listing, error)
	GetFeatured(ctx context.Context) ([]*Listing, error)
	SearchByCity(ctx context.Context, city string) ([]*Listing, error)
	GetByPriceRange(ctx context.Context, minPrice, maxPrice int64) ([]*Listing, error)
	GetByBedroomRange(ctx context.Context, minBedrooms, maxBedrooms int) ([]*Listing, error)
	GetByBathroomRange(ctx context.Context, minBathrooms, maxBathrooms int) ([]*Listing, error)
}

// ListingRepositoryImpl implements the ListingRepository interface
type ListingRepositoryImpl struct {
	data   map[int64]*Listing
	mu     sync.RWMutex
	nextID int64
}

// NewListingRepository creates a new listing repository
func NewListingRepository() ListingRepository {
	repo := &ListingRepositoryImpl{
		data:   make(map[int64]*Listing),
		nextID: 1,
	}

	// Add some sample data for testing
	repo.addSampleData()

	return repo
}

// addSampleData adds sample listings for testing
func (r *ListingRepositoryImpl) addSampleData() {
	sampleListings := []*Listing{
		{
			ID: 187,
			AddressDetails: AddressDetails{
				AddressLine1:      "5 Camden High Street",
				AddressLine2:      "",
				City:              "London",
				Postcode:          "N1 7AA",
				ShortenedPostcode: "N17",
				Country:           "UK",
				Region:            RegionLondon,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   1,
			Bathrooms:                  1,
			SizeSqFt:                   50,
			PriceInCents:               12500000,
			MinimumDepositInCents:      1000000,
			EstimatedDepositInCents:    3125000,
			MonthlyRentalIncomeInCents: 110000,
			IsTenanted:                 true,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "property",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/1b2b53fd-398b-4129-8f7d-c5932f90b3c3",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/1b2b53fd-398b-4129-8f7d-c5932f90b3c3_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/1b2b53fd-398b-4129-8f7d-c5932f90b3c3_thumbnail",
					MimeType:     "image/png",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/7c8d16b4-09d1-453b-8729-e7bfada38b2e",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/7c8d16b4-09d1-453b-8729-e7bfada38b2e_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/7c8d16b4-09d1-453b-8729-e7bfada38b2e_thumbnail",
					MimeType:     "image/png",
				},
			},
			GrossYield:    0.1056,
			MadeVisibleAt: nil,
		},
		{
			ID: 185,
			AddressDetails: AddressDetails{
				AddressLine1:      "2B fire blvd",
				AddressLine2:      "",
				City:              "Ashford",
				Postcode:          "",
				ShortenedPostcode: "ASF",
				Country:           "UK",
				Region:            RegionSouthEast,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   1,
			Bathrooms:                  1,
			SizeSqFt:                   2342,
			PriceInCents:               10000000,
			MinimumDepositInCents:      2550000,
			EstimatedDepositInCents:    2500000,
			MonthlyRentalIncomeInCents: 30000,
			IsTenanted:                 true,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "asdf",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/ba11810b-30fc-4061-b3f5-126e2aae0a95",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/ba11810b-30fc-4061-b3f5-126e2aae0a95_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/ba11810b-30fc-4061-b3f5-126e2aae0a95_thumbnail",
					MimeType:     "image/jpeg",
				},
			},
			GrossYield:    0.036,
			MadeVisibleAt: nil,
		},
		{
			ID: 79,
			AddressDetails: AddressDetails{
				AddressLine1:      "37 John Snow Drive",
				AddressLine2:      "",
				City:              "Wallington",
				Postcode:          "SM6 4ER",
				ShortenedPostcode: "SM6",
				Country:           "UK",
				Region:            RegionLondon,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   2,
			Bathrooms:                  1,
			SizeSqFt:                   300,
			PriceInCents:               10000000,
			MinimumDepositInCents:      1000000,
			EstimatedDepositInCents:    2500000,
			MonthlyRentalIncomeInCents: 60000,
			IsTenanted:                 true,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "test",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/dc1c52ca-1061-4673-a3ae-92bd9098189d",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/dc1c52ca-1061-4673-a3ae-92bd9098189d_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/dc1c52ca-1061-4673-a3ae-92bd9098189d_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/82540ded-ea98-4279-a57d-742c0495ae73",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/82540ded-ea98-4279-a57d-742c0495ae73_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/82540ded-ea98-4279-a57d-742c0495ae73_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/bb5ad33a-717c-4525-beea-7f811fb828ef",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/bb5ad33a-717c-4525-beea-7f811fb828ef_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/bb5ad33a-717c-4525-beea-7f811fb828ef_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/57ef8821-1c60-4569-9999-21c23851d284",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/57ef8821-1c60-4569-9999-21c23851d284_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/57ef8821-1c60-4569-9999-21c23851d284_thumbnail",
					MimeType:     "image/jpeg",
				},
			},
			GrossYield:    0.072,
			MadeVisibleAt: stringPtr("2023-02-01T16:42:09Z"),
		},
		{
			ID: 80,
			AddressDetails: AddressDetails{
				AddressLine1:      "87 Scotts way",
				AddressLine2:      "",
				City:              "Edinburgh",
				Postcode:          "EH12 5AA",
				ShortenedPostcode: "EH12",
				Country:           "UK",
				Region:            RegionScotland,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   0,
			Bathrooms:                  1,
			SizeSqFt:                   355,
			PriceInCents:               23456700,
			MinimumDepositInCents:      3000000,
			EstimatedDepositInCents:    18798136,
			MonthlyRentalIncomeInCents: 200000,
			IsTenanted:                 true,
			IsCashOnly:                 true,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                true,
			Description:                "Share Sale Test",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/32cdb9d4-f02e-4d4c-84a0-35318874c9c7",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/32cdb9d4-f02e-4d4c-84a0-35318874c9c7_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/32cdb9d4-f02e-4d4c-84a0-35318874c9c7_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/ada31aac-71b6-4e89-b75b-20cbef617b8c",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/ada31aac-71b6-4e89-b75b-20cbef617b8c_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/ada31aac-71b6-4e89-b75b-20cbef617b8c_thumbnail",
					MimeType:     "image/jpeg",
				},
			},
			GrossYield:    0.102316,
			MadeVisibleAt: stringPtr("2023-02-01T16:52:50Z"),
		},
		{
			ID: 81,
			AddressDetails: AddressDetails{
				AddressLine1:      "202 Grunge street",
				AddressLine2:      "",
				City:              "Manchester",
				Postcode:          "M1 1AA",
				ShortenedPostcode: "M1",
				Country:           "UK",
				Region:            RegionNorthWest,
			},
			PropertyType:               PropertyTypeSemiDetached,
			Bedrooms:                   3,
			Bathrooms:                  2,
			SizeSqFt:                   789,
			PriceInCents:               14400000,
			MinimumDepositInCents:      5200000,
			EstimatedDepositInCents:    3600000,
			MonthlyRentalIncomeInCents: 200000,
			IsTenanted:                 true,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  true,
			IsShareSale:                true,
			Description:                "GG Company test",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/44e55a81-c66e-48b6-abef-9b8cb3a5b674",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/44e55a81-c66e-48b6-abef-9b8cb3a5b674_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/44e55a81-c66e-48b6-abef-9b8cb3a5b674_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/f8fa432d-0306-4d2c-8d57-f33523372f26",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/f8fa432d-0306-4d2c-8d57-f33523372f26_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/f8fa432d-0306-4d2c-8d57-f33523372f26_thumbnail",
					MimeType:     "image/jpeg",
				},
			},
			GrossYield:    0.166667,
			MadeVisibleAt: stringPtr("2023-02-01T17:12:25Z"),
		},
		{
			ID: 82,
			AddressDetails: AddressDetails{
				AddressLine1:      "3 Knotting Lane",
				AddressLine2:      "",
				City:              "London",
				Postcode:          "W14 9AA",
				ShortenedPostcode: "W14",
				Country:           "UK",
				Region:            RegionLondon,
			},
			PropertyType:               PropertyTypeTerraced,
			Bedrooms:                   3,
			Bathrooms:                  3,
			SizeSqFt:                   300,
			PriceInCents:               100000000,
			MinimumDepositInCents:      20000000,
			EstimatedDepositInCents:    40827520,
			MonthlyRentalIncomeInCents: 850000,
			IsTenanted:                 false,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "Image test",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/58e9d12f-95fd-4516-9611-cc6f4c828959",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/58e9d12f-95fd-4516-9611-cc6f4c828959_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/58e9d12f-95fd-4516-9611-cc6f4c828959_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/f8647f62-a6a2-4bb5-a3ac-b22759af1804",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/f8647f62-a6a2-4bb5-a3ac-b22759af1804_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/f8647f62-a6a2-4bb5-a3ac-b22759af1804_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/0ccd6f6a-9390-41b8-8102-8a205ffe73cd",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/0ccd6f6a-9390-41b8-8102-8a205ffe73cd_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/0ccd6f6a-9390-41b8-8102-8a205ffe73cd_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/fb1a1f3b-313d-45a4-b1d3-1ee5c5f21895",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/fb1a1f3b-313d-45a4-b1d3-1ee5c5f21895_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/fb1a1f3b-313d-45a4-b1d3-1ee5c5f21895_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/62ed85f1-00ac-40bd-9607-780f5245ce1e",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/62ed85f1-00ac-40bd-9607-780f5245ce1e_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/62ed85f1-00ac-40bd-9607-780f5245ce1e_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/a39d1c3b-8a93-4d9c-bdb4-d6cda911f8e4",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/a39d1c3b-8a93-4d9c-bdb4-d6cda911f8e4_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/a39d1c3b-8a93-4d9c-bdb4-d6cda911f8e4_thumbnail",
					MimeType:     "image/jpeg",
				},
			},
			GrossYield:    0.102,
			MadeVisibleAt: stringPtr("2023-02-02T08:36:00Z"),
		},
		{
			ID: 68,
			AddressDetails: AddressDetails{
				AddressLine1:      "27 Hill Road",
				AddressLine2:      "",
				City:              "Sheffield",
				Postcode:          "S1 2AB",
				ShortenedPostcode: "S1",
				Country:           "UK",
				Region:            RegionNorthEast,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   1,
			Bathrooms:                  1,
			SizeSqFt:                   301,
			PriceInCents:               13875000,
			MinimumDepositInCents:      3468700,
			EstimatedDepositInCents:    3468750,
			MonthlyRentalIncomeInCents: 95100,
			IsTenanted:                 true,
			IsCashOnly:                 true,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "This property can be purchased via shares with 0% SDLT - this purchase option is only available with GetGround! \n\nWhen you buy shares, you usually pay stamp duty tax of 0.5% on the price you pay for the claims. As the property continues to be owned by the company no SDLT is payable.\n\nThis modern studio apartment is ideal for young professionals in the heart of Sheffield city centre. This flat is already tenanted generating a strong 8.2% yield for investors to benefit from. The rent is assured until the end of Q2 2025.\n\nBeing just a 10-minute walk from Sheffield Sheaf Street, you have fantastic access to all of the North's city hubs, with Manchester, Birmingham and Leeds all just a one-hour train journey away.",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/48d16039-5857-471e-a8ee-d427441e2604",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/48d16039-5857-471e-a8ee-d427441e2604_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/48d16039-5857-471e-a8ee-d427441e2604_thumbnail",
					MimeType:     "image/png",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/bd4a9d7b-9416-4d61-b550-0b01f2e7fa72",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/bd4a9d7b-9416-4d61-b550-0b01f2e7fa72_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/bd4a9d7b-9416-4d61-b550-0b01f2e7fa72_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/fff510fd-33d5-41a4-ab95-a2d5033551c7",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/fff510fd-33d5-41a4-ab95-a2d5033551c7_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/fff510fd-33d5-41a4-ab95-a2d5033551c7_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/7538ea96-08e3-4ce1-8c2c-dbadc24f7622",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/7538ea96-08e3-4ce1-8c2c-dbadc24f7622_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/7538ea96-08e3-4ce1-8c2c-dbadc24f7622_thumbnail",
					MimeType:     "image/jpeg",
				},
			},
			GrossYield:    0.0822486,
			MadeVisibleAt: stringPtr("2023-09-27T08:14:37Z"),
		},
		{
			ID: 66,
			AddressDetails: AddressDetails{
				AddressLine1:      "67 Kansas Street",
				AddressLine2:      "",
				City:              "Preston",
				Postcode:          "PR1 2TT",
				ShortenedPostcode: "PR1",
				Country:           "UK",
				Region:            RegionNorthWest,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   2,
			Bathrooms:                  1,
			SizeSqFt:                   686,
			PriceInCents:               3995000,
			MinimumDepositInCents:      3995000,
			EstimatedDepositInCents:    998750,
			MonthlyRentalIncomeInCents: 38000,
			IsTenanted:                 true,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "2 bed flat in Preston with a rear terrace space and large garden area.",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/5d4d5076-2cf3-4d09-8818-6a8fe0993530",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/5d4d5076-2cf3-4d09-8818-6a8fe0993530_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/5d4d5076-2cf3-4d09-8818-6a8fe0993530_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/cca8f249-5946-44d3-98f6-2b7f1ec0c286",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/cca8f249-5946-44d3-98f6-2b7f1ec0c286_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/cca8f249-5946-44d3-98f6-2b7f1ec0c286_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/5db35350-6efa-4c23-b1df-b8d2bdce32e3",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/5db35350-6efa-4c23-b1df-b8d2bdce32e3_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/5db35350-6efa-4c23-b1df-b8d2bdce32e3_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/bdc6a7da-4405-4721-a820-897f907cbcc7",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/bdc6a7da-4405-4721-a820-897f907cbcc7_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/bdc6a7da-4405-4721-a820-897f907cbcc7_thumbnail",
					MimeType:     "image/jpeg",
				},
			},
			GrossYield:    0.114143,
			MadeVisibleAt: stringPtr("2023-01-25T15:50:34Z"),
		},
		{
			ID: 71,
			AddressDetails: AddressDetails{
				AddressLine1:      "45 Larcher Street",
				AddressLine2:      "",
				City:              "Preston",
				Postcode:          "PR1 2AF",
				ShortenedPostcode: "PR1",
				Country:           "UK",
				Region:            RegionNorthWest,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   3,
			Bathrooms:                  1,
			SizeSqFt:                   840,
			PriceInCents:               88058000,
			MinimumDepositInCents:      7014500,
			EstimatedDepositInCents:    42397020,
			MonthlyRentalIncomeInCents: 875400,
			IsTenanted:                 false,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "Apt 170 sits on the eleventh floor and provides a great investment opportunity into a high specification apartment in a prime location that has been identified by local and national government as an area of great potential and has been extremely well funded in recent years. ",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/26e7679c-5ba9-4c9c-ba54-858b400e6863",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/26e7679c-5ba9-4c9c-ba54-858b400e6863_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/26e7679c-5ba9-4c9c-ba54-858b400e6863_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/5bee58eb-5bbd-4329-b884-eea1dad850e1",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/5bee58eb-5bbd-4329-b884-eea1dad850e1_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/5bee58eb-5bbd-4329-b884-eea1dad850e1_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/b9ecaccd-964e-4adb-aec6-863bef629166",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/b9ecaccd-964e-4adb-aec6-863bef629166_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/b9ecaccd-964e-4adb-aec6-863bef629166_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/19c9fd2f-67dc-4b87-855f-f023bbbb4e85",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/19c9fd2f-67dc-4b87-855f-f023bbbb4e85_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/19c9fd2f-67dc-4b87-855f-f023bbbb4e85_thumbnail",
					MimeType:     "image/jpeg",
				},
			},
			GrossYield:    0.119294,
			MadeVisibleAt: stringPtr("2023-01-25T16:11:54Z"),
		},
		{
			ID: 72,
			AddressDetails: AddressDetails{
				AddressLine1:      "3 Button Street",
				AddressLine2:      "",
				City:              "Preston",
				Postcode:          "PR1 2AF",
				ShortenedPostcode: "PR1",
				Country:           "UK",
				Region:            RegionNorthWest,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   2,
			Bathrooms:                  1,
			SizeSqFt:                   678,
			PriceInCents:               22429100,
			MinimumDepositInCents:      5607300,
			EstimatedDepositInCents:    9514036,
			MonthlyRentalIncomeInCents: 140200,
			IsTenanted:                 false,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "Apt 53 situated in Block B is a high specification apartment that sits on the third floor. The development itself encompasses a beautiful roof garden, residents lounge, concierge service, bike storage and a state of the art gym to complete this premium home in a well invested area, resulting in a rapidly developing cultural and economic landscape. \n\n\nTEST",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/5f3758df-98cb-427a-84f2-afc7ea6c40ca",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/5f3758df-98cb-427a-84f2-afc7ea6c40ca_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/5f3758df-98cb-427a-84f2-afc7ea6c40ca_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/35658a75-7f58-4549-8008-e2bbbdbdb20a",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/35658a75-7f58-4549-8008-e2bbbdbdb20a_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/35658a75-7f58-4549-8008-e2bbbdbdb20a_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/365fc59a-4b48-448e-a110-3b790df58db9",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/365fc59a-4b48-448e-a110-3b790df58db9_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/365fc59a-4b48-448e-a110-3b790df58db9_thumbnail",
					MimeType:     "image/jpeg",
				},
			},
			GrossYield:    0.0750097,
			MadeVisibleAt: stringPtr("2023-02-17T18:20:03Z"),
		},
		{
			ID: 105,
			AddressDetails: AddressDetails{
				AddressLine1:      "4 High Street",
				AddressLine2:      "",
				City:              "Wallington",
				Postcode:          "SM6 9AA",
				ShortenedPostcode: "SM6",
				Country:           "UK",
				Region:            RegionSouthEast,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   0,
			Bathrooms:                  1,
			SizeSqFt:                   286,
			PriceInCents:               52500000,
			MinimumDepositInCents:      2500000,
			EstimatedDepositInCents:    13125000,
			MonthlyRentalIncomeInCents: 350000,
			IsTenanted:                 true,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                true,
			Description:                "ertggr rtg trtg etgtrgrt",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/fa320e9c-cb0b-4c0d-be21-2b992ee6c126",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/fa320e9c-cb0b-4c0d-be21-2b992ee6c126_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/fa320e9c-cb0b-4c0d-be21-2b992ee6c126_thumbnail",
					MimeType:     "image/png",
				},
			},
			GrossYield:    0.08,
			MadeVisibleAt: stringPtr("2023-03-01T13:39:00Z"),
		},
		{
			ID: 106,
			AddressDetails: AddressDetails{
				AddressLine1:      "13 Cliff st",
				AddressLine2:      "",
				City:              "Dover",
				Postcode:          "",
				ShortenedPostcode: "DOV",
				Country:           "UK",
				Region:            RegionSouthEast,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   1,
			Bathrooms:                  2,
			SizeSqFt:                   123,
			PriceInCents:               12300000,
			MinimumDepositInCents:      1230000,
			EstimatedDepositInCents:    1500000,
			MonthlyRentalIncomeInCents: 80000,
			IsTenanted:                 true,
			IsCashOnly:                 true,
			IsNewBuild:                 false,
			IsCompany:                  true,
			IsShareSale:                true,
			Description:                "UPDATE",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/9db6de1b-88b9-4030-b944-bd15544a23ed",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/9db6de1b-88b9-4030-b944-bd15544a23ed_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/9db6de1b-88b9-4030-b944-bd15544a23ed_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/b71de69c-cfe3-41bf-8e80-a0a3602026fa",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/b71de69c-cfe3-41bf-8e80-a0a3602026fa_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/b71de69c-cfe3-41bf-8e80-a0a3602026fa_thumbnail",
					MimeType:     "image/jpeg",
				},
			},
			GrossYield:    0.0780488,
			MadeVisibleAt: nil,
		},
		{
			ID: 120,
			AddressDetails: AddressDetails{
				AddressLine1:      "77 Regata Street",
				AddressLine2:      "",
				City:              "Henley",
				Postcode:          "RE1 1AA",
				ShortenedPostcode: "RE1",
				Country:           "UK",
				Region:            RegionSouthEast,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   0,
			Bathrooms:                  1,
			SizeSqFt:                   400,
			PriceInCents:               100000000,
			MinimumDepositInCents:      10000000,
			EstimatedDepositInCents:    25000000,
			MonthlyRentalIncomeInCents: 700000,
			IsTenanted:                 true,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "test",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/86f3dadf-54bf-447c-a4fa-a8f067cf6aee",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/86f3dadf-54bf-447c-a4fa-a8f067cf6aee_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/86f3dadf-54bf-447c-a4fa-a8f067cf6aee_thumbnail",
					MimeType:     "image/jpeg",
				},
			},
			GrossYield:    0.084,
			MadeVisibleAt: stringPtr("2023-03-16T16:11:14Z"),
		},
		{
			ID: 91,
			AddressDetails: AddressDetails{
				AddressLine1:      "Animal Farm",
				AddressLine2:      "",
				City:              "Eastbourne",
				Postcode:          "BN20 1AA",
				ShortenedPostcode: "BN20",
				Country:           "UK",
				Region:            RegionSouthEast,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   0,
			Bathrooms:                  1,
			SizeSqFt:                   789,
			PriceInCents:               100000000,
			MinimumDepositInCents:      10000000,
			EstimatedDepositInCents:    25000000,
			MonthlyRentalIncomeInCents: 700000,
			IsTenanted:                 true,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/8816708f-9d0d-4c24-bac9-dad96fef53f4",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/8816708f-9d0d-4c24-bac9-dad96fef53f4_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/8816708f-9d0d-4c24-bac9-dad96fef53f4_thumbnail",
					MimeType:     "image/jpeg",
				},
			},
			GrossYield:    0.084,
			MadeVisibleAt: stringPtr("2023-02-17T17:47:45Z"),
		},
		{
			ID: 94,
			AddressDetails: AddressDetails{
				AddressLine1:      "12 Alfie Solomons Way",
				AddressLine2:      "",
				City:              "London",
				Postcode:          "N1 8LN",
				ShortenedPostcode: "N1",
				Country:           "UK",
				Region:            RegionLondon,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   2,
			Bathrooms:                  3,
			SizeSqFt:                   1000,
			PriceInCents:               125000,
			MinimumDepositInCents:      10000,
			EstimatedDepositInCents:    31250,
			MonthlyRentalIncomeInCents: 1100,
			IsTenanted:                 false,
			IsCashOnly:                 true,
			IsNewBuild:                 false,
			IsCompany:                  true,
			IsShareSale:                true,
			Description:                "Free text",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/bea385af-17c5-4369-bf94-be6b9570f4db",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/bea385af-17c5-4369-bf94-be6b9570f4db_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/bea385af-17c5-4369-bf94-be6b9570f4db_thumbnail",
					MimeType:     "image/png",
				},
			},
			GrossYield:    0.1056,
			MadeVisibleAt: stringPtr("2023-02-20T09:57:55Z"),
		},
		{
			ID: 97,
			AddressDetails: AddressDetails{
				AddressLine1:      "123 Rainbow Road",
				AddressLine2:      "",
				City:              "Brighton",
				Postcode:          "BN1 1AA",
				ShortenedPostcode: "BN1",
				Country:           "UK",
				Region:            RegionSouthEast,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   0,
			Bathrooms:                  1,
			SizeSqFt:                   678,
			PriceInCents:               50000000,
			MinimumDepositInCents:      10000000,
			EstimatedDepositInCents:    12000000,
			MonthlyRentalIncomeInCents: 400000,
			IsTenanted:                 true,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "0",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/c836d323-7548-47fd-82e8-f1529c87e51e",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/c836d323-7548-47fd-82e8-f1529c87e51e_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/c836d323-7548-47fd-82e8-f1529c87e51e_thumbnail",
					MimeType:     "image/jpeg",
				},
			},
			GrossYield:    0.096,
			MadeVisibleAt: nil,
		},
		{
			ID: 103,
			AddressDetails: AddressDetails{
				AddressLine1:      "4 Seaside Road",
				AddressLine2:      "",
				City:              "Whitstable",
				Postcode:          "CT5 1AB",
				ShortenedPostcode: "CT5",
				Country:           "UK",
				Region:            RegionSouthEast,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   0,
			Bathrooms:                  1,
			SizeSqFt:                   300,
			PriceInCents:               100000000,
			MinimumDepositInCents:      10000000,
			EstimatedDepositInCents:    12000000,
			MonthlyRentalIncomeInCents: 100000,
			IsTenanted:                 true,
			IsCashOnly:                 true,
			IsNewBuild:                 false,
			IsCompany:                  true,
			IsShareSale:                false,
			Description:                "66666",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/533f2688-23c2-4385-a272-0179f9d6b6db",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/533f2688-23c2-4385-a272-0179f9d6b6db_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/533f2688-23c2-4385-a272-0179f9d6b6db_thumbnail",
					MimeType:     "image/jpeg",
				},
			},
			GrossYield:    0.012,
			MadeVisibleAt: nil,
		},
		{
			ID: 143,
			AddressDetails: AddressDetails{
				AddressLine1:      "3 Chaucer Road",
				AddressLine2:      "",
				City:              "Canterbury",
				Postcode:          "CT1 3RF",
				ShortenedPostcode: "CT1",
				Country:           "UK",
				Region:            RegionSouthEast,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   0,
			Bathrooms:                  1,
			SizeSqFt:                   750,
			PriceInCents:               20000000,
			MinimumDepositInCents:      300000,
			EstimatedDepositInCents:    500000,
			MonthlyRentalIncomeInCents: 150000,
			IsTenanted:                 true,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "Test ",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/3a458a51-532c-456a-8814-5d488cac49f0",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/3a458a51-532c-456a-8814-5d488cac49f0_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/3a458a51-532c-456a-8814-5d488cac49f0_thumbnail",
					MimeType:     "image/png",
				},
			},
			GrossYield:    0.09,
			MadeVisibleAt: stringPtr("2023-03-27T11:01:38Z"),
		},
		{
			ID: 144,
			AddressDetails: AddressDetails{
				AddressLine1:      "115 Maidstone Road",
				AddressLine2:      "",
				City:              "Maidstone",
				Postcode:          "ME15 6AA",
				ShortenedPostcode: "MDS",
				Country:           "UK",
				Region:            RegionSouthEast,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   0,
			Bathrooms:                  1,
			SizeSqFt:                   40,
			PriceInCents:               15000000,
			MinimumDepositInCents:      1000000,
			EstimatedDepositInCents:    1500000,
			MonthlyRentalIncomeInCents: 40000,
			IsTenanted:                 true,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "Test",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/0ee871a1-ecb8-43c7-aab0-30128baec351",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/0ee871a1-ecb8-43c7-aab0-30128baec351_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/0ee871a1-ecb8-43c7-aab0-30128baec351_thumbnail",
					MimeType:     "image/png",
				},
			},
			GrossYield:    0.032,
			MadeVisibleAt: stringPtr("2023-03-27T11:04:16Z"),
		},
		{
			ID: 148,
			AddressDetails: AddressDetails{
				AddressLine1:      "3 Buckingham Palace Road",
				AddressLine2:      "",
				City:              "London",
				Postcode:          "",
				ShortenedPostcode: "W14 8FF",
				Country:           "UK",
				Region:            RegionLondon,
			},
			PropertyType:               PropertyTypeTerraced,
			Bedrooms:                   3,
			Bathrooms:                  1,
			SizeSqFt:                   32,
			PriceInCents:               35000000,
			MinimumDepositInCents:      10000000,
			EstimatedDepositInCents:    8750000,
			MonthlyRentalIncomeInCents: 400000,
			IsTenanted:                 true,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "test",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/70c29877-0654-49ed-8225-3189e81b7354",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/70c29877-0654-49ed-8225-3189e81b7354_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/70c29877-0654-49ed-8225-3189e81b7354_thumbnail",
					MimeType:     "image/png",
				},
			},
			GrossYield:    0.137143,
			MadeVisibleAt: stringPtr("2023-03-28T13:30:27Z"),
		},
		{
			ID: 145,
			AddressDetails: AddressDetails{
				AddressLine1:      "123 Main street",
				AddressLine2:      "W14 9AA",
				City:              "London",
				Postcode:          "",
				ShortenedPostcode: "W14",
				Country:           "UK",
				Region:            RegionLondon,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   0,
			Bathrooms:                  1,
			SizeSqFt:                   40,
			PriceInCents:               25000000,
			MinimumDepositInCents:      10000000,
			EstimatedDepositInCents:    6250000,
			MonthlyRentalIncomeInCents: 400000,
			IsTenanted:                 true,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "Classic house right on main street!",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/b73ac33d-b2de-4656-a92e-a7f4bde98b26",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/b73ac33d-b2de-4656-a92e-a7f4bde98b26_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/b73ac33d-b2de-4656-a92e-a7f4bde98b26_thumbnail",
					MimeType:     "image/png",
				},
			},
			GrossYield:    0.192,
			MadeVisibleAt: stringPtr("2023-03-27T12:36:10Z"),
		},
		{
			ID: 178,
			AddressDetails: AddressDetails{
				AddressLine1:      "331 High street",
				AddressLine2:      "",
				City:              "Spooky City",
				Postcode:          "",
				ShortenedPostcode: "SW2",
				Country:           "UK",
				Region:            RegionLondon,
			},
			PropertyType:               PropertyTypeTerraced,
			Bedrooms:                   3,
			Bathrooms:                  3,
			SizeSqFt:                   200,
			PriceInCents:               9999900,
			MinimumDepositInCents:      999900,
			EstimatedDepositInCents:    1050000,
			MonthlyRentalIncomeInCents: 99900,
			IsTenanted:                 true,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  true,
			IsShareSale:                true,
			Description:                "Boo! You've found our ghost listing for this Halloween season. You might have fallen for our trick, but using our buy-to-let marketplace is a treat. Scroll through our properties for some scarily good yields.",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/ab08a0af-0a78-4244-8f9e-985a517c85b4",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/ab08a0af-0a78-4244-8f9e-985a517c85b4_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/ab08a0af-0a78-4244-8f9e-985a517c85b4_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/e4df5172-4bde-4d7b-b11e-4ed8258c2bc6",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/e4df5172-4bde-4d7b-b11e-4ed8258c2bc6_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/e4df5172-4bde-4d7b-b11e-4ed8258c2bc6_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/4650d79c-196e-4a15-83ec-a146f8dbf64d",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/4650d79c-196e-4a15-83ec-a146f8dbf64d_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/4650d79c-196e-4a15-83ec-a146f8dbf64d_thumbnail",
					MimeType:     "image/jpeg",
				},
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/8779e7e5-4181-45b1-9d0f-f41e8a15dcd6",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/8779e7e5-4181-45b1-9d0f-f41e8a15dcd6_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/8779e7e5-4181-45b1-9d0f-f41e8a15dcd6_thumbnail",
					MimeType:     "image/jpeg",
				},
			},
			GrossYield:    0.119988,
			MadeVisibleAt: stringPtr("2023-10-13T13:15:07Z"),
		},
		{
			ID: 183,
			AddressDetails: AddressDetails{
				AddressLine1:      "1 Maidstone Road",
				AddressLine2:      "",
				City:              "Maidstone",
				Postcode:          "ME15 6AA",
				ShortenedPostcode: "ME15",
				Country:           "UK",
				Region:            RegionSouthEast,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   1,
			Bathrooms:                  1,
			SizeSqFt:                   1234,
			PriceInCents:               10000000,
			MinimumDepositInCents:      2500000,
			EstimatedDepositInCents:    2500000,
			MonthlyRentalIncomeInCents: 150000,
			IsTenanted:                 true,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "Apartment in central Maidstone",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/2889f267-fd95-43fd-8745-45e4e02f1e59",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/2889f267-fd95-43fd-8745-45e4e02f1e59_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/2889f267-fd95-43fd-8745-45e4e02f1e59_thumbnail",
					MimeType:     "image/jpeg",
				},
			},
			GrossYield:    0.18,
			MadeVisibleAt: nil,
		},
		{
			ID: 181,
			AddressDetails: AddressDetails{
				AddressLine1:      "5 Canterbury Road",
				AddressLine2:      "",
				City:              "Cantebury",
				Postcode:          "CT1 1AA",
				ShortenedPostcode: "CT1",
				Country:           "UK",
				Region:            RegionSouthEast,
			},
			PropertyType:               PropertyTypeApartment,
			Bedrooms:                   1,
			Bathrooms:                  1,
			SizeSqFt:                   566,
			PriceInCents:               25000000,
			MinimumDepositInCents:      5000000,
			EstimatedDepositInCents:    6500000,
			MonthlyRentalIncomeInCents: 150000,
			IsTenanted:                 true,
			IsCashOnly:                 false,
			IsNewBuild:                 false,
			IsCompany:                  false,
			IsShareSale:                false,
			Description:                "Very slick apartment in Canterbury",
			Photos: []Photo{
				{
					OriginalURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/6306344f-423f-4e67-be20-207bdb11eed8",
					StandardURL:  "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/6306344f-423f-4e67-be20-207bdb11eed8_standard",
					ThumbnailURL: "https://storage.googleapis.com/assets-terranova-qa-module-core/listings/6306344f-423f-4e67-be20-207bdb11eed8_thumbnail",
					MimeType:     "image/png",
				},
			},
			GrossYield:    0.072,
			MadeVisibleAt: nil,
		},
	}

	for _, listing := range sampleListings {
		r.data[listing.ID] = listing
		if listing.ID >= r.nextID {
			r.nextID = listing.ID + 1
		}
	}
}

// stringPtr returns a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// Create adds a new listing to the repository
func (r *ListingRepositoryImpl) Create(ctx context.Context, listing *Listing) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if listing.AddressDetails.City == "" {
		return errors.New("city is required")
	}
	if listing.AddressDetails.ShortenedPostcode == "" {
		return errors.New("shortened postcode is required")
	}
	if listing.AddressDetails.Region == "" {
		return errors.New("region is required")
	}
	if listing.PropertyType == "" {
		return errors.New("property type is required")
	}
	if listing.PriceInCents <= 0 {
		return errors.New("price must be greater than 0")
	}

	listing.ID = r.nextID
	now := time.Now().Format(time.RFC3339)
	if listing.MadeVisibleAt == nil {
		listing.MadeVisibleAt = &now
	}
	r.data[listing.ID] = listing
	r.nextID++
	return nil
}

// GetByID retrieves a listing by its ID
func (r *ListingRepositoryImpl) GetByID(ctx context.Context, id int64) (*Listing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	listing, exists := r.data[id]
	if !exists {
		return nil, errors.Wrapf(errors.New("not found"), "listing not found with id: %d", id)
	}
	return listing, nil
}

// GetAll retrieves all listings
func (r *ListingRepositoryImpl) GetAll(ctx context.Context) ([]*Listing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	listings := make([]*Listing, 0, len(r.data))
	for _, listing := range r.data {
		listings = append(listings, listing)
	}
	return listings, nil
}

// Update updates an existing listing
func (r *ListingRepositoryImpl) Update(ctx context.Context, listing *Listing) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if listing.AddressDetails.City == "" {
		return errors.New("city is required")
	}
	if listing.AddressDetails.ShortenedPostcode == "" {
		return errors.New("shortened postcode is required")
	}
	if listing.AddressDetails.Region == "" {
		return errors.New("region is required")
	}
	if listing.PropertyType == "" {
		return errors.New("property type is required")
	}
	if listing.PriceInCents <= 0 {
		return errors.New("price must be greater than 0")
	}

	existing, exists := r.data[listing.ID]
	if !exists {
		return errors.Wrapf(errors.New("not found"), "listing not found with id: %d", listing.ID)
	}

	// Preserve the original MadeVisibleAt if it exists
	if existing.MadeVisibleAt != nil && listing.MadeVisibleAt == nil {
		listing.MadeVisibleAt = existing.MadeVisibleAt
	}

	r.data[listing.ID] = listing
	return nil
}

// Delete removes a listing by its ID
func (r *ListingRepositoryImpl) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.data[id]; !exists {
		return errors.Wrapf(errors.New("not found"), "listing not found with id: %d", id)
	}
	delete(r.data, id)
	return nil
}

// GetByRegion retrieves all listings in a specific region
func (r *ListingRepositoryImpl) GetByRegion(ctx context.Context, region string) ([]*Listing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	listings := make([]*Listing, 0)
	for _, listing := range r.data {
		if string(listing.AddressDetails.Region) == region {
			listings = append(listings, listing)
		}
	}
	return listings, nil
}

// GetByPropertyType retrieves all listings of a specific property type
func (r *ListingRepositoryImpl) GetByPropertyType(ctx context.Context, propertyType string) ([]*Listing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	listings := make([]*Listing, 0)
	for _, listing := range r.data {
		if string(listing.PropertyType) == propertyType {
			listings = append(listings, listing)
		}
	}
	return listings, nil
}

// GetFeatured retrieves all featured listings (deprecated - returns all listings)
func (r *ListingRepositoryImpl) GetFeatured(ctx context.Context) ([]*Listing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	listings := make([]*Listing, 0)
	for _, listing := range r.data {
		listings = append(listings, listing)
	}
	return listings, nil
}

// SearchByCity searches listings by city
func (r *ListingRepositoryImpl) SearchByCity(ctx context.Context, city string) ([]*Listing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	listings := make([]*Listing, 0)
	cityLower := strings.ToLower(city)
	for _, listing := range r.data {
		if strings.Contains(strings.ToLower(listing.AddressDetails.City), cityLower) {
			listings = append(listings, listing)
		}
	}
	return listings, nil
}

// GetByPriceRange retrieves listings within a price range
func (r *ListingRepositoryImpl) GetByPriceRange(ctx context.Context, minPrice, maxPrice int64) ([]*Listing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	listings := make([]*Listing, 0)
	for _, listing := range r.data {
		if listing.PriceInCents >= minPrice && listing.PriceInCents <= maxPrice {
			listings = append(listings, listing)
		}
	}
	return listings, nil
}

// GetByBedroomRange retrieves listings within a bedroom range
func (r *ListingRepositoryImpl) GetByBedroomRange(ctx context.Context, minBedrooms, maxBedrooms int) ([]*Listing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	listings := make([]*Listing, 0)
	for _, listing := range r.data {
		if listing.Bedrooms >= minBedrooms && listing.Bedrooms <= maxBedrooms {
			listings = append(listings, listing)
		}
	}
	return listings, nil
}

// GetByBathroomRange retrieves listings within a bathroom range
func (r *ListingRepositoryImpl) GetByBathroomRange(ctx context.Context, minBathrooms, maxBathrooms int) ([]*Listing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	listings := make([]*Listing, 0)
	for _, listing := range r.data {
		if listing.Bathrooms >= minBathrooms && listing.Bathrooms <= maxBathrooms {
			listings = append(listings, listing)
		}
	}
	return listings, nil
}
