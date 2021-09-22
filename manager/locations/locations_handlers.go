package locations

import (
	"strings"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type LocationHandler struct {
	LSer *LocationService
}

func NewLocationHandler() *LocationHandler {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}
	lSer := NewLocationServiceImpl(conf, &dbMgr)
	return &LocationHandler{
		LSer: &lSer,
	}
}

// GetLocationsHandler handle locations get cli command
func (mHdl *LocationHandler) GetLocations() {
	(*mHdl.LSer).GetLocations()
}

// AddLocationHandler handle location add cli command
func (mHdl *LocationHandler) AddLocation(airportCode string)  (models.Location, error) {
	airport, err := (*mHdl.LSer).FindAirport(airportCode)
	if err != nil {
		return models.Location{}, err
	}

	coords := strings.Split(airport.Coordinates, ", ")
	location := models.Location{
		Country: airport.IsoCountry,
		IataCode: airport.IataCode,
		Latitude:    coords[0],
		Longitude: coords[1],
	}
	
	location = (*mHdl.LSer).AddLocation(location)
	return location, nil
}

// DeleteLocation handle location delete cli command
func (mHdl *LocationHandler) DeleteLocation(countryCode string) {
	location := models.Location{
		Country: countryCode,
	}
	(*mHdl.LSer).DeleteLocation(location)
}