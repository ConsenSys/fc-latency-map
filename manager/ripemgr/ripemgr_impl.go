package ripemgr

import (
	"fmt"
	"math"
	"net"
	"strings"
	"time"

	atlas "github.com/keltia/ripe-atlas"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ConsenSys/fc-latency-map/manager/addresses"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type RipeMgrImpl struct {
	conf viper.Viper
	c    atlas.Client
}

const (
	startTimeDelay                = 10
	delayBetweenBatchMeasurements = 600
	// atlas not permitted run more than 100 concurrent measurements
	maxRunningConcurrentMeasurements = 100.0
)

func NewRipeImpl(conf *viper.Viper) (RipeMgr, error) {
	cfgs := []atlas.Config{}
	cfgs = append(cfgs, atlas.Config{
		APIKey: conf.GetString("RIPE_API_KEY"),
	})
	c, err := atlas.NewClient(cfgs...)
	if err != nil {
		log.Println("Connecting to Ripe Atlas API", err)

		return nil, err
	}
	ver := atlas.GetVersion()
	log.Println("api version", ver)

	return &RipeMgrImpl{
		c:    *c,
		conf: *conf,
	}, nil
}

func (rMgr *RipeMgrImpl) GetProbes(opts map[string]string) ([]atlas.Probe, error) {
	log.WithFields(log.Fields{
		"filters": opts,
	}).Info("started get probes from Atlas Ripe")
	probes, err := rMgr.c.GetProbes(opts)
	if err != nil {
		return nil, err
	}

	return probes, nil
}

func (rMgr *RipeMgrImpl) GetMeasurement(measurementID int) (*atlas.Measurement, error) {
	measurement, err := rMgr.c.GetMeasurement(measurementID)
	if err != nil {
		log.WithFields(log.Fields{
			"err":         err,
			"measurement": measurementID,
		}).Error("get measurements results")
		return nil, err
	}

	return measurement, nil
}

func (rMgr *RipeMgrImpl) GetMeasurementResults(measurementID int) ([]atlas.MeasurementResult, error) {
	log.WithFields(log.Fields{
		"measurement": measurementID,
	}).Warn("get measurements results")
	measurementResult, err := rMgr.c.GetResults(measurementID)
	if err != nil {
		log.WithFields(log.Fields{
			"err":         err,
			"measurement": measurementID,
		}).Error("get measurements results")
		return nil, err
	}

	return measurementResult.Results, nil
}

func (rMgr *RipeMgrImpl) CreateMeasurements(miners []*models.Miner, probeIDs string, t int) ([]*atlas.Measurement, error) {
	if len(miners) == 0 {
		log.WithFields(log.Fields{
			"msg": "miners are missing",
		}).Warn("Create Measurements")
		return nil, nil
	}
	if probeIDs == "" {
		log.WithFields(log.Fields{
			"msg": "probeIDs are missing",
		}).Warn("Create Measurements")
		return nil, nil
	}
	probes := []atlas.ProbeSet{
		{
			Type:      "probes",
			Value:     probeIDs,
			Requested: rMgr.getRequestedProbes(probeIDs),
		},
	}

	return rMgr.createTraceroute(miners, probes, t)
}

func (rMgr *RipeMgrImpl) getRequestedProbes(probeIDs string) int {
	requestedProbes := rMgr.conf.GetInt("RIPE_REQUESTED_PROBES")
	if requestedProbes == 0 {
		return len(strings.Split(probeIDs, ","))
	}
	return requestedProbes
}

func (rMgr *RipeMgrImpl) createTraceroute(miners []*models.Miner, probes []atlas.ProbeSet, t int) ([]*atlas.Measurement, error) {
	var d []atlas.Definition

	isOneOff := rMgr.conf.GetBool("RIPE_ONE_OFF")

	for _, miner := range miners {
		if miner.IP == "" || (miner.Latitude == 0 && miner.Longitude == 0) {
			continue
		}
		d = rMgr.getDefinitions(miner, d)
	}

	mr := rMgr.getMeasurementRequest(d, isOneOff, probes, t)

	if !isOneOff {
		runningTime := rMgr.conf.GetInt("RIPE_TRACEROUTE_RUNNING_TIME")
		mr.StopTime = mr.StartTime + runningTime
	}

	p, err := rMgr.c.Traceroute(mr)
	if err != nil {
		log.WithFields(log.Fields{
			"msg": mr,
			"err": err.Error(),
		}).Error("Create traceroute")

		return nil, err
	}

	log.WithFields(log.Fields{
		"id":          p,
		"isOneOff":    isOneOff,
		"measurement": fmt.Sprintf("%#v", d),
	}).Info("create newMeasurement")

	var measurement []*atlas.Measurement
	for _, v := range p.Measurements {
		measurement = append(measurement, &atlas.Measurement{
			ID:        v,
			StopTime:  mr.StopTime,
			StartTime: mr.StartTime,
			IsOneoff:  mr.IsOneoff,
		})
	}

	return measurement, err
}

func (rMgr *RipeMgrImpl) getMeasurementRequest(d []atlas.Definition, isOneOff bool, probes []atlas.ProbeSet, t int) *atlas.MeasurementRequest {
	startTime := int(time.Now().Unix()) + startTimeDelay

	if t >= maxRunningConcurrentMeasurements {
		startTime += int(delayBetweenBatchMeasurements * math.Ceil(float64(t)/maxRunningConcurrentMeasurements))
	}
	return &atlas.MeasurementRequest{
		Definitions: d,
		StartTime:   startTime,
		IsOneoff:    isOneOff,
		Probes:      probes,
	}
}

func (rMgr *RipeMgrImpl) getDefinitions(miner *models.Miner, d []atlas.Definition) []atlas.Definition {
	isOneOff := rMgr.conf.GetBool("RIPE_ONE_OFF")
	packets := rMgr.conf.GetInt("RIPE_PACKETS")
	tracerouteInterval := rMgr.conf.GetInt("RIPE_TRACEROUTE_INTERVAL")

	for _, ip := range strings.Split(miner.IP, ",") {
		ipAdd := net.ParseIP(ip)
		if ipAdd.IsPrivate() {
			continue
		}

		const defaultSize = 48
		const hops = 32
		definition := atlas.Definition{
			Description: fmt.Sprintf("%s traceroute to %s", miner.Address, ip),
			AF:          addresses.GetIPVersion(ipAdd),
			Target:      ip,
			Packets:     packets,
			Tags: []string{
				miner.Address,
			},
			Size:                  defaultSize,
			Protocol:              "TCP",
			Type:                  "traceroute",
			FirstHop:              hops,
			MaxHops:               hops,
			Port:                  miner.Port,
			DestinationOptionSize: 0,
			HopByHopOptionSize:    0,
			DontFragment:          false,
		}
		if !isOneOff {
			definition.Interval = tracerouteInterval
		}
		d = append(d, definition)
	}

	return d
}
