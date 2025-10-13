# MCP Tools API Reference

Complete reference documentation for all 9 MCP tools provided by the Recreation Opportunities MCP Server.

## Table of Contents

- [search_parks](#search_parks)
- [get_park_details](#get_park_details)
- [get_park_alerts](#get_park_alerts)
- [search_campgrounds](#search_campgrounds)
- [search_recreation_areas](#search_recreation_areas)
- [get_facility_details](#get_facility_details)
- [get_weather](#get_weather)
- [get_weather_forecast](#get_weather_forecast)
- [list_activities](#list_activities)

---

## search_parks

Search for national parks by name, state, or activity.

### Input Schema

```json
{
  "query": "string (optional) - Search text",
  "state": "string (optional) - Two-letter state code (e.g., 'CA', 'CO')",
  "activity": "string (optional) - Activity name (e.g., 'Hiking', 'Camping')",
  "limit": "integer (optional, default: 10) - Maximum number of results"
}
```

### Output

Returns a list of parks with:
- Park name
- Park code (4-letter identifier)
- Description
- State(s)
- Designation (National Park, Monument, etc.)
- URL for more information

### Example Usage

**Query:** "Find parks in California with hiking"

**Input:**
```json
{
  "state": "CA",
  "activity": "Hiking",
  "limit": 5
}
```

**Response:**
```json
{
  "parks": [
    {
      "name": "Yosemite National Park",
      "parkCode": "yose",
      "states": "CA",
      "description": "...",
      "designation": "National Park",
      "url": "https://www.nps.gov/yose/index.htm"
    }
  ],
  "count": 5
}
```

### Common State Codes
- CA (California)
- CO (Colorado)
- AZ (Arizona)
- UT (Utah)
- WA (Washington)
- [Full list of US state codes](https://en.wikipedia.org/wiki/List_of_U.S._state_and_territory_abbreviations)

---

## get_park_details

Get comprehensive information about a specific national park.

### Input Schema

```json
{
  "park_code": "string (required) - Four-letter park code (e.g., 'yose')"
}
```

### Output

Returns detailed park information:
- Full name and description
- Operating hours (by season)
- Entrance fees and passes
- Weather information
- Contact information
- Directions and addresses
- Available activities
- Operating hours and seasons
- Images and multimedia

### Example Usage

**Query:** "Tell me about Yosemite National Park"

**Input:**
```json
{
  "park_code": "yose"
}
```

**Response:**
```json
{
  "name": "Yosemite National Park",
  "description": "...",
  "states": "CA",
  "latLong": "lat:37.84883288, long:-119.5571873",
  "activities": [
    {
      "name": "Hiking"
    },
    {
      "name": "Rock Climbing"
    }
  ],
  "entranceFees": [
    {
      "cost": "35.00",
      "description": "Vehicle permit - valid for 7 days",
      "title": "Yosemite National Park Entrance Fee - Private Vehicle"
    }
  ],
  "operatingHours": [...],
  "addresses": [...],
  "contacts": {...}
}
```

### Common Park Codes
- `yose` - Yosemite National Park
- `grca` - Grand Canyon National Park
- `yell` - Yellowstone National Park
- `zion` - Zion National Park
- `romo` - Rocky Mountain National Park
- `acad` - Acadia National Park
- `grsm` - Great Smoky Mountains National Park
- `olym` - Olympic National Park

---

## get_park_alerts

Get current alerts, closures, and important notices for a specific park.

### Input Schema

```json
{
  "park_code": "string (required) - Four-letter park code"
}
```

### Output

Returns list of alerts with:
- Alert title
- Description
- Category (Park Closure, Road Closure, Caution, Information, etc.)
- Date posted
- Effective dates

### Example Usage

**Query:** "Are there any alerts or closures at Grand Canyon?"

**Input:**
```json
{
  "park_code": "grca"
}
```

**Response:**
```json
{
  "alerts": [
    {
      "title": "North Rim Closed for Winter",
      "description": "The North Rim of Grand Canyon National Park is closed...",
      "category": "Park Closure",
      "lastIndexedDate": "2025-10-13",
      "parkCode": "grca"
    }
  ],
  "count": 1
}
```

### Alert Categories
- **Park Closure**: Entire park or sections closed
- **Road Closure**: Specific roads or trails closed
- **Caution**: Safety warnings, wildlife, weather
- **Information**: General updates, events
- **Danger**: Immediate hazards

---

## search_campgrounds

Search for campgrounds in national parks or recreation areas.

### Input Schema

```json
{
  "park_code": "string (optional) - NPS park code",
  "state": "string (optional) - Two-letter state code",
  "query": "string (optional) - Search text",
  "limit": "integer (optional, default: 10) - Maximum results"
}
```

### Output

Returns list of campgrounds with:
- Name
- Description
- Amenities (restrooms, water, RV hookups, etc.)
- Reservation information
- Contact information
- Accessibility information

### Example Usage

**Query:** "Find campgrounds in Yellowstone"

**Input:**
```json
{
  "park_code": "yell",
  "limit": 5
}
```

**Response:**
```json
{
  "campgrounds": [
    {
      "name": "Madison Campground",
      "description": "...",
      "amenities": {
        "toilets": ["Flush Toilets"],
        "showers": ["None"],
        "cellPhoneReception": "Yes",
        "internetConnectivity": "No"
      },
      "reservations": {
        "reservationInfo": "Reservations required",
        "reservationUrl": "https://www.recreation.gov/..."
      }
    }
  ],
  "count": 5
}
```

### Common Amenities
- Flush toilets / Vault toilets
- Potable water
- Showers
- RV hookups (electric, water, sewer)
- Dump station
- Amphitheater
- Camp store
- Cell phone reception
- Picnic tables
- Fire rings

---

## search_recreation_areas

Search for recreation areas on Recreation.gov across all federal lands.

### Input Schema

```json
{
  "query": "string (optional) - Search text",
  "state": "string (optional) - Two-letter state code",
  "activity": "string (optional) - Activity name",
  "limit": "integer (optional, default: 10) - Maximum results"
}
```

### Output

Returns list of recreation areas with:
- Name
- Description
- Managing agency (Forest Service, BLM, Army Corps of Engineers, etc.)
- Activities available
- Facilities
- Location information

### Example Usage

**Query:** "Show me recreation areas in Utah with mountain biking"

**Input:**
```json
{
  "state": "UT",
  "activity": "Mountain Biking",
  "limit": 5
}
```

**Response:**
```json
{
  "recreationAreas": [
    {
      "name": "Moab Recreation Area",
      "description": "...",
      "activities": ["Mountain Biking", "Hiking", "Camping"],
      "orgName": "Bureau of Land Management",
      "reservationInfo": "..."
    }
  ],
  "count": 5
}
```

### Managing Agencies
- Forest Service (USFS)
- Bureau of Land Management (BLM)
- Army Corps of Engineers
- Bureau of Reclamation
- Fish and Wildlife Service
- National Park Service

---

## get_facility_details

Get detailed information about a specific facility (campground, cabin, day-use area, etc.) from Recreation.gov.

### Input Schema

```json
{
  "facility_id": "string (required) - Recreation.gov facility ID"
}
```

### Output

Returns detailed facility information:
- Name and description
- Facility type
- Amenities
- Reservation requirements
- Contact information
- Directions
- Photos and media

### Example Usage

**Query:** "Get details for facility 234567"

**Input:**
```json
{
  "facility_id": "234567"
}
```

**Response:**
```json
{
  "name": "Mirror Lake Campground",
  "description": "...",
  "facilityType": "Campground",
  "amenities": [...],
  "reservationUrl": "https://www.recreation.gov/...",
  "phone": "555-1234",
  "adaAccess": "Accessible"
}
```

### Facility Types
- Campground
- Cabin/Lookout
- Day Use Area
- Permit Area
- Recreation Area
- Tour
- Ticket Facility

---

## get_weather

Get current weather conditions for a specific location.

### Input Schema

```json
{
  "latitude": "number (required) - Latitude coordinate",
  "longitude": "number (required) - Longitude coordinate",
  "units": "string (optional, default: 'imperial') - 'metric' or 'imperial'"
}
```

### Output

Returns current weather:
- Temperature (current, feels like)
- Weather conditions
- Humidity
- Wind speed and direction
- Atmospheric pressure
- Visibility
- Cloud cover
- Sunrise/sunset times

### Example Usage

**Query:** "What's the current weather at Yosemite?"

**Input:**
```json
{
  "latitude": 37.8651,
  "longitude": -119.5383,
  "units": "imperial"
}
```

**Response:**
```json
{
  "temperature": 72.5,
  "feelsLike": 71.3,
  "description": "Clear sky",
  "humidity": 45,
  "windSpeed": 5.2,
  "windDirection": 180,
  "pressure": 1013,
  "visibility": 10000,
  "clouds": 10,
  "sunrise": "2025-10-13T06:45:00Z",
  "sunset": "2025-10-13T18:30:00Z"
}
```

### Units
- **Imperial**: Fahrenheit, mph, inches
- **Metric**: Celsius, m/s, mm

### Getting Coordinates
You can get coordinates from:
- park details (latLong field)
- Google Maps (right-click → coordinates)
- GPS devices
- Coordinate lookup services

---

## get_weather_forecast

Get 5-day weather forecast for a specific location.

### Input Schema

```json
{
  "latitude": "number (required) - Latitude coordinate",
  "longitude": "number (required) - Longitude coordinate",
  "units": "string (optional, default: 'imperial') - 'metric' or 'imperial'"
}
```

### Output

Returns 5-day forecast with 3-hour intervals:
- Date and time
- Temperature (high/low)
- Weather conditions
- Precipitation probability
- Wind speed
- Humidity
- Cloud cover

### Example Usage

**Query:** "What's the weather forecast for Grand Canyon this week?"

**Input:**
```json
{
  "latitude": 36.0544,
  "longitude": -112.1401,
  "units": "imperial"
}
```

**Response:**
```json
{
  "forecast": [
    {
      "date": "2025-10-13T12:00:00Z",
      "temperature": 68.0,
      "tempMin": 52.0,
      "tempMax": 73.0,
      "description": "Clear sky",
      "precipitation": 0,
      "humidity": 35,
      "windSpeed": 8.5
    }
  ],
  "count": 40
}
```

### Forecast Details
- **Time intervals**: Every 3 hours for 5 days (40 data points)
- **Update frequency**: Every 10 minutes
- **Coverage**: Global

---

## list_activities

List all available activities across NPS and Recreation.gov systems.

### Input Schema

```json
{
  "source": "string (optional, default: 'all') - 'nps', 'recreation_gov', or 'all'",
  "limit": "integer (optional, default: 50) - Maximum results"
}
```

### Output

Returns list of activities with:
- Activity name
- Activity ID
- Source (NPS or Recreation.gov)

### Example Usage

**Query:** "What outdoor activities are available?"

**Input:**
```json
{
  "source": "all",
  "limit": 20
}
```

**Response:**
```json
{
  "activities": [
    {
      "name": "Hiking",
      "id": "BFF8C027-7C8F-480B-A5F8-CD8CE490BFBA",
      "source": "nps"
    },
    {
      "name": "Camping",
      "id": "9159DF0F-951D-4AAE-9987-CEB3CE2A9ADA",
      "source": "nps"
    },
    {
      "name": "Rock Climbing",
      "id": "5FF5B286-E9C3-440E-B5A4-B4B5BFB6F45F8",
      "source": "nps"
    }
  ],
  "count": 20
}
```

### Common Activities
- Hiking
- Camping
- Rock Climbing
- Wildlife Watching
- Photography
- Fishing
- Boating
- Kayaking
- Mountain Biking
- Horseback Riding
- Stargazing
- Snowshoeing
- Cross-Country Skiing
- Birdwatching
- Scenic Driving

### Using Activities
Use activity names from this list when searching parks or recreation areas to ensure matches.

---

## Error Handling

All tools return errors in MCP format:

```json
{
  "error": {
    "code": "error_code",
    "message": "Human-readable error message"
  }
}
```

### Common Error Codes

- `invalid_params` - Missing or invalid input parameters
- `api_error` - External API returned an error
- `rate_limit` - Rate limit exceeded on external API
- `not_found` - Resource not found
- `internal_error` - Server internal error
- `network_error` - Network connectivity issue

### Error Examples

**Invalid Parameter:**
```json
{
  "error": {
    "code": "invalid_params",
    "message": "park_code is required"
  }
}
```

**Rate Limit:**
```json
{
  "error": {
    "code": "rate_limit",
    "message": "Rate limit exceeded. Please try again in 60 seconds."
  }
}
```

**Not Found:**
```json
{
  "error": {
    "code": "not_found",
    "message": "Park with code 'xxxx' not found"
  }
}
```

---

## Best Practices

### 1. Use Specific Parameters
More specific queries return better results:
- ✅ Good: `{"state": "CA", "activity": "Hiking", "limit": 5}`
- ❌ Less optimal: `{"query": "parks"}`

### 2. Handle Empty Results
Not all queries will return results. Check the count field:
```json
{
  "parks": [],
  "count": 0
}
```

### 3. Combine Tools
Use multiple tools together for comprehensive information:
1. `search_parks` - Find parks
2. `get_park_details` - Get details
3. `get_weather` - Check weather
4. `get_park_alerts` - Check alerts

### 4. Cache Awareness
Results are cached for 1 hour. Repeated queries will be faster.

### 5. Coordinate Precision
Use at least 4 decimal places for coordinates:
- ✅ Good: `37.8651, -119.5383`
- ❌ Too imprecise: `37.9, -119.5`

---

## Rate Limits

### External API Limits
- **NPS API**: 1000 requests per hour
- **Recreation.gov**: 1000 requests per day (free tier)
- **OpenWeatherMap**: 1000 requests per day (free tier), 60/minute

### MCP Server
- No hard limits on MCP tool calls
- Caching reduces external API usage significantly
- Configure `MAX_REQUESTS_PER_MINUTE` if needed

---

## Additional Resources

- [National Park Service API Documentation](https://www.nps.gov/subjects/developer/api-documentation.htm)
- [Recreation.gov API Documentation](https://ridb.recreation.gov/docs)
- [OpenWeatherMap API Documentation](https://openweathermap.org/api)
- [MCP Protocol Specification](https://modelcontextprotocol.io/docs)

---

**Last Updated**: October 13, 2025
