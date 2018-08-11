package googlePlaces

import (
	"errors"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	"theAmazingCodeExample/app/config"
)

func GetValidAddress(addressText string) (addressName string, lat float64, long float64, err error) {

	//Create google maps client
	googleMapsClient, err := maps.NewClient(maps.WithAPIKey(config.GetConfig().GOOGLE_PLACES_API_KEY))
	if err != nil {
		return "", 0, 0, err
	}

	//Add request parameters
	request := maps.GeocodingRequest{
		Address:    addressText,
		Components: map[maps.Component]string{"country": "ar"},
	}

	//Get results
	result, err := googleMapsClient.Geocode(context.Background(), &request)
	if err != nil {
		return "", 0, 0, err
	}

	//Get the first prediction and return it's description
	if len(result) == 0 {
		return "", 0, 0, errors.New("Invalid address")
	} else {
		return result[0].FormattedAddress, result[0].Geometry.Location.Lat, result[0].Geometry.Location.Lng, nil
	}

}
