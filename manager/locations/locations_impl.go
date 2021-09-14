package locations

import (
	"log"

	"github.com/spf13/viper"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type LocationServiceImpl struct {
	Conf  *viper.Viper
	DbMgr *db.DatabaseMgr
}

func NewLocationServiceImpl(conf *viper.Viper, dbMgr *db.DatabaseMgr) LocationService {
	return &LocationServiceImpl{
		Conf:  conf,
		DbMgr: dbMgr,
	}
}

func (srv *LocationServiceImpl) GetLocations() []*models.Location {
	var locsList = []*models.Location{}
	(*srv.DbMgr).GetDb().Find(&locsList)
	for _, location := range locsList {
		log.Printf("ID:%d - Country code: %s\n", location.ID, location.Country)
	}
	return locsList
}

func (srv *LocationServiceImpl) GetLocation(location models.Location) models.Location {
	if err := (*srv.DbMgr).GetDb().Where(location).First(&location).Error; err != nil {
		return models.Location{}
	}
	return location
}

func (srv *LocationServiceImpl) AddLocation(newLocation models.Location) models.Location {
	var location = models.Location{}
	(*srv.DbMgr).GetDb().Where(&newLocation).First(&location)
	if location == (models.Location{}) {
		(*srv.DbMgr).GetDb().Create(&newLocation) 
		log.Printf("New location, ID:%d - Country code: %s\n", newLocation.ID, newLocation.Country)
	} else {
		log.Printf("Location already exists, ID:%d\n", location.ID)
	}
	return location
}

func (srv *LocationServiceImpl) DeleteLocation(location models.Location) bool {
	location = srv.GetLocation(location)
	if (location == models.Location{}) {
		log.Printf("Unable to find location %s\n", location.Country)
	} else {
		(*srv.DbMgr).GetDb().Delete(&location)
		log.Printf("Location %d deleted\n", location.ID)
	}
	
	return true
}