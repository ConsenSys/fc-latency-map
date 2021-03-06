package models

import (
	"gorm.io/gorm"
)

type Measurement struct {
	gorm.Model
	IsOneOff       bool
	MeasurementID  int `gorm:"index:idx_measurement_id,unique"`
	StartTime      int
	StopTime       int
	Status         string
	StatusStopTime int
}

type MeasurementResult struct {
	gorm.Model
	ProbeID              int         `gorm:"foreignKey:probe_id;index:idx_name,unique;index:idx_measurement_probe_id"`
	Probe                Probe       `gorm:"foreignkey:ProbeID;references:probe_id"`
	MeasurementID        int         `gorm:"foreignKey:measurement_id;index:idx_name,unique"`
	Measurement          Measurement `gorm:"foreignkey:MeasurementID;references:measurement_id"`
	MeasurementTimestamp int         `gorm:"index:idx_name,unique"`
	IP                   string      `gorm:"index:idx_name,unique;index:idx_measurement_ip"`
	MeasurementDate      string      `gorm:"index:idx_measurement_date"`
	Rtt                  float64
}
