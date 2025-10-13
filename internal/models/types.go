package models

import "time"

// Park represents a national park from the NPS API
type Park struct {
	ID             string           `json:"id"`
	ParkCode       string           `json:"parkCode"`
	Name           string           `json:"name"`
	FullName       string           `json:"fullName"`
	Description    string           `json:"description"`
	States         string           `json:"states"`
	Latitude       string           `json:"latitude"`
	Longitude      string           `json:"longitude"`
	URL            string           `json:"url"`
	Activities     []Activity       `json:"activities,omitempty"`
	Contacts       *Contact         `json:"contacts,omitempty"`
	EntranceFees   []EntranceFee    `json:"entranceFees,omitempty"`
	Hours          []OperatingHours `json:"operatingHours,omitempty"`
	Addresses      []Address        `json:"addresses,omitempty"`
	Images         []Image          `json:"images,omitempty"`
	WeatherInfo    string           `json:"weatherInfo,omitempty"`
	DirectionsInfo string           `json:"directionsInfo,omitempty"`
}

// Activity represents an activity available at a park or facility
type Activity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Contact represents contact information
type Contact struct {
	PhoneNumbers   []PhoneNumber  `json:"phoneNumbers,omitempty"`
	EmailAddresses []EmailAddress `json:"emailAddresses,omitempty"`
}

// PhoneNumber represents a phone contact
type PhoneNumber struct {
	PhoneNumber string `json:"phoneNumber"`
	Description string `json:"description"`
	Extension   string `json:"extension,omitempty"`
	Type        string `json:"type"`
}

// EmailAddress represents an email contact
type EmailAddress struct {
	EmailAddress string `json:"emailAddress"`
	Description  string `json:"description"`
}

// EntranceFee represents park entrance fee information
type EntranceFee struct {
	Cost        string `json:"cost"`
	Description string `json:"description"`
	Title       string `json:"title"`
}

// OperatingHours represents operating hours for a facility
type OperatingHours struct {
	Name          string            `json:"name,omitempty"`
	Description   string            `json:"description,omitempty"`
	StandardHours map[string]string `json:"standardHours,omitempty"`
	Exceptions    []HourException   `json:"exceptions,omitempty"`
}

// HourException represents an exception to standard hours
type HourException struct {
	Name           string            `json:"name"`
	StartDate      string            `json:"startDate"`
	EndDate        string            `json:"endDate"`
	ExceptionHours map[string]string `json:"exceptionHours"`
}

// Address represents a physical or mailing address
type Address struct {
	Type        string `json:"type"`
	Line1       string `json:"line1"`
	Line2       string `json:"line2,omitempty"`
	Line3       string `json:"line3,omitempty"`
	City        string `json:"city"`
	StateCode   string `json:"stateCode"`
	PostalCode  string `json:"postalCode"`
	CountryCode string `json:"countryCode"`
}

// Image represents an image associated with a park or facility
type Image struct {
	Credit  string `json:"credit"`
	Title   string `json:"title"`
	AltText string `json:"altText"`
	Caption string `json:"caption"`
	URL     string `json:"url"`
}

// Alert represents a park alert or closure
type Alert struct {
	ID              string `json:"id"`
	ParkCode        string `json:"parkCode"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	URL             string `json:"url,omitempty"`
	LastIndexedDate string `json:"lastIndexedDate,omitempty"`
}

// Campground represents a campground facility
type Campground struct {
	ID                 string          `json:"id"`
	Name               string          `json:"name"`
	ParkCode           string          `json:"parkCode,omitempty"`
	Description        string          `json:"description"`
	Latitude           string          `json:"latitude,omitempty"`
	Longitude          string          `json:"longitude,omitempty"`
	Amenities          []string        `json:"amenities,omitempty"`
	Accessibility      []string        `json:"accessibility,omitempty"`
	Campsites          []Campsite      `json:"campsites,omitempty"`
	Contacts           *Contact        `json:"contacts,omitempty"`
	Hours              *OperatingHours `json:"operatingHours,omitempty"`
	Fees               []Fee           `json:"fees,omitempty"`
	DirectionsOverview string          `json:"directionsOverview,omitempty"`
	ReservationInfo    string          `json:"reservationInfo,omitempty"`
	ReservationURL     string          `json:"reservationUrl,omitempty"`
}

// Campsite represents an individual campsite
type Campsite struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type,omitempty"`
	Loop       string                 `json:"loop,omitempty"`
	Accessible bool                   `json:"accessible"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// Fee represents a fee structure
type Fee struct {
	Cost        string `json:"cost"`
	Description string `json:"description"`
	Title       string `json:"title"`
}

// RecreationArea represents a recreation area from Recreation.gov
type RecreationArea struct {
	RecAreaID             string               `json:"RecAreaID"`
	RecAreaName           string               `json:"RecAreaName"`
	RecAreaDescription    string               `json:"RecAreaDescription"`
	RecAreaFeeDescription string               `json:"RecAreaFeeDescription,omitempty"`
	RecAreaLatitude       float64              `json:"RecAreaLatitude,omitempty"`
	RecAreaLongitude      float64              `json:"RecAreaLongitude,omitempty"`
	StayLimit             string               `json:"StayLimit,omitempty"`
	Keywords              string               `json:"Keywords,omitempty"`
	Activities            []RecreationActivity `json:"ACTIVITY,omitempty"`
	Facilities            []string             `json:"FACILITYID,omitempty"`
	RecAreaPhone          string               `json:"RecAreaPhone,omitempty"`
	RecAreaEmail          string               `json:"RecAreaEmail,omitempty"`
}

// RecreationActivity represents an activity from Recreation.gov
type RecreationActivity struct {
	ActivityID          string `json:"ActivityID"`
	ActivityName        string `json:"ActivityName"`
	ActivityDescription string `json:"ActivityDescription,omitempty"`
}

// Facility represents a facility from Recreation.gov
type Facility struct {
	FacilityID              string               `json:"FacilityID"`
	FacilityName            string               `json:"FacilityName"`
	FacilityDescription     string               `json:"FacilityDescription"`
	FacilityTypeDescription string               `json:"FacilityTypeDescription,omitempty"`
	FacilityLatitude        float64              `json:"FacilityLatitude,omitempty"`
	FacilityLongitude       float64              `json:"FacilityLongitude,omitempty"`
	FacilityPhone           string               `json:"FacilityPhone,omitempty"`
	FacilityEmail           string               `json:"FacilityEmail,omitempty"`
	FacilityReservationURL  string               `json:"FacilityReservationURL,omitempty"`
	FacilityAdaAccess       string               `json:"FacilityAdaAccess,omitempty"`
	Activities              []RecreationActivity `json:"ACTIVITY,omitempty"`
	Addresses               []FacilityAddress    `json:"FACILITYADDRESS,omitempty"`
}

// FacilityAddress represents an address for a Recreation.gov facility
type FacilityAddress struct {
	FacilityAddressID      string `json:"FacilityAddressID"`
	FacilityAddressType    string `json:"FacilityAddressType"`
	FacilityStreetAddress1 string `json:"FacilityStreetAddress1,omitempty"`
	FacilityStreetAddress2 string `json:"FacilityStreetAddress2,omitempty"`
	FacilityStreetAddress3 string `json:"FacilityStreetAddress3,omitempty"`
	City                   string `json:"City"`
	PostalCode             string `json:"PostalCode"`
	StateCode              string `json:"AddressStateCode"`
	CountryCode            string `json:"AddressCountryCode"`
}

// Weather represents current weather conditions from OpenWeatherMap
type Weather struct {
	Temperature   float64 `json:"temp"`
	FeelsLike     float64 `json:"feels_like"`
	TempMin       float64 `json:"temp_min"`
	TempMax       float64 `json:"temp_max"`
	Pressure      int     `json:"pressure"`
	Humidity      int     `json:"humidity"`
	Conditions    string  `json:"conditions"`
	Description   string  `json:"description"`
	WindSpeed     float64 `json:"wind_speed"`
	WindDirection int     `json:"wind_deg,omitempty"`
	Clouds        int     `json:"clouds"`
	Visibility    int     `json:"visibility,omitempty"`
	Sunrise       int64   `json:"sunrise,omitempty"`
	Sunset        int64   `json:"sunset,omitempty"`
	Timezone      int     `json:"timezone,omitempty"`
}

// WeatherForecast represents a weather forecast from OpenWeatherMap
type WeatherForecast struct {
	Date        time.Time `json:"dt"`
	Temperature struct {
		Day   float64 `json:"day"`
		Min   float64 `json:"min"`
		Max   float64 `json:"max"`
		Night float64 `json:"night"`
		Eve   float64 `json:"eve"`
		Morn  float64 `json:"morn"`
	} `json:"temp"`
	FeelsLike struct {
		Day   float64 `json:"day"`
		Night float64 `json:"night"`
		Eve   float64 `json:"eve"`
		Morn  float64 `json:"morn"`
	} `json:"feels_like"`
	Pressure  int                `json:"pressure"`
	Humidity  int                `json:"humidity"`
	Weather   []WeatherCondition `json:"weather"`
	WindSpeed float64            `json:"speed"`
	Clouds    int                `json:"clouds"`
	Pop       float64            `json:"pop"`
	Rain      float64            `json:"rain,omitempty"`
	Snow      float64            `json:"snow,omitempty"`
}

// WeatherCondition represents a weather condition
type WeatherCondition struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// APIResponse represents a generic API response wrapper
type APIResponse struct {
	Total int         `json:"total,omitempty"`
	Data  interface{} `json:"data"`
	Limit int         `json:"limit,omitempty"`
	Start int         `json:"start,omitempty"`
	Error *APIError   `json:"error,omitempty"`
}

// APIError represents an error from an external API
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NPSResponse represents a response from the NPS API
type NPSResponse struct {
	Total string      `json:"total"`
	Data  interface{} `json:"data"`
	Limit string      `json:"limit"`
	Start string      `json:"start"`
}

// RecreationGovResponse represents a response from Recreation.gov API
type RecreationGovResponse struct {
	RECDATA  interface{} `json:"RECDATA"`
	METADATA struct {
		RESULTS struct {
			CurrentCount int `json:"CURRENT_COUNT"`
			TotalCount   int `json:"TOTAL_COUNT"`
		} `json:"RESULTS"`
	} `json:"METADATA"`
}
