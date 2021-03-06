package export

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/file"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type ExportServiceImpl struct {
	Conf  *viper.Viper
	DBMgr db.DatabaseMgr
}

func newExportServiceImpl(conf *viper.Viper, dbMgr db.DatabaseMgr) Service {
	return &ExportServiceImpl{
		Conf:  conf,
		DBMgr: dbMgr,
	}
}

func (m *ExportServiceImpl) export() *[]string {
	files := []string{}
	dates := m.getDates()

	if len(dates) == 0 {
		log.Warn("No dates to generate exports")
	}

	for _, date := range dates {
		fn := fmt.Sprintf("export_%s.json", date)
		if file.IsUpdated(fn, date) {
			log.WithFields(
				map[string]interface{}{
					"file": fn,
				},
			).Info("file exists")
			continue
		}
		log.WithFields(
			map[string]interface{}{
				"file": fn,
			},
		).Info("generate file")

		measurements := m.getLatencyMeasurementsStored(date)
		log.Info("Measurements retrieved, start generating file")
		j := m.marshalJSON(measurements)
		file.Create(fn, j)
		files = append(files, fn)
	}
	log.Info("Export completed")
	return &files
}

func (m *ExportServiceImpl) marshalJSON(measurements *Result) []byte {
	fullJSON, err := json.MarshalIndent(measurements, "", "  ")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Create json data")

		return nil
	}
	return fullJSON
}

func (m *ExportServiceImpl) getLatencyMeasurementsStored(date string) *Result {
	results := &Result{Measurements: map[string]map[string][]*Miner{}}
	loc := m.getLocations()
	miners := m.getMiners()
	for _, l := range loc {
		log.WithFields(log.Fields{
			"iataCode": l.IataCode,
		}).Print("Get latency measurements")
		for _, miner := range miners {
			latency := &Miner{
				Address:  miner.Address,
				Measures: []*MeasureIP{},
			}
			if miner.IP == "" {
				continue
			}

			latency = m.appendLatency(latency, int(l.Model.ID), miner.IP, date)
			if len(latency.Measures) == 0 {
				continue
			}
			if _, found := results.Measurements[l.Country]; !found {
				results.Measurements[l.Country] = make(map[string][]*Miner)
			}
			results.Measurements[l.Country][l.IataCode] = append(results.Measurements[l.Country][l.IataCode], latency)
		}
	}

	results.Locations = loc
	results.Miners = miners
	results.Dates = m.getDates()

	return results
}

func (m *ExportServiceImpl) appendLatency(latency *Miner, locationID int, ip, date string) *Miner {
	for _, ip := range strings.Split(ip, ",") {
		meas := m.getMeasureResults(date, ip, locationID)
		if meas == nil || meas.IP == "" {
			continue
		}
		latency.Measures = append(latency.Measures, meas)
	}

	return latency
}

func (m *ExportServiceImpl) getMeasureResults(date, ip string, locationID int) *MeasureIP {
	var meas *MeasureIP
	where := &MeasureIP{
		IP:              ip,
		MeasurementDate: date,
	}

	dbc := m.DBMgr.GetDB()
	err := dbc.Model(&models.MeasurementResult{}).Select(
		"ip, avg(rtt) time_average").
		Where(where).
		Where("probe_id in (?)",
			dbc.Select("probe_id").
				Table("probes").
				Where("id in (?)",
					dbc.Select("probe_id").
						Table("locations_probes").
						Where("location_id in (?)", locationID))).
		Group("ip, measurement_date").
		First(&meas).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("measurement result")
		}

		return nil
	}

	return meas
}

func (m *ExportServiceImpl) getMiners() []*models.Miner {
	var miners []*models.Miner

	err := m.DBMgr.GetDB().Find(&miners).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("getMiners")

		return nil
	}

	return miners
}

func (m *ExportServiceImpl) getLocations() []*models.Location {
	var loc []*models.Location
	err := m.DBMgr.GetDB().
		Preload(clause.Associations).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "country"}}).
		Find(&loc).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GetAllLocations")

		return nil
	}

	return loc
}

func (m *ExportServiceImpl) getDates() []string {
	var dates []string
	m.DBMgr.GetDB().
		Model(&models.MeasurementResult{}).
		Distinct().
		Order(clause.OrderByColumn{Column: clause.Column{Name: "measurement_date"}}).
		Pluck("measurement_date", &dates)
	return dates
}
