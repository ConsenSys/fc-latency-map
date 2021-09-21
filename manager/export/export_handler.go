package export

import (
	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
)

type ExportHandler struct {
	Service *ExportService
}

func NewExportHandler() *ExportHandler {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}

	mSer := NewExportServiceImpl(conf, &dbMgr)

	return &ExportHandler{
		Service: &mSer,
	}
}

func (h *ExportHandler) Export(fn string) {
	(*h.Service).export(fn)
}