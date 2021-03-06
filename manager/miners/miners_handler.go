package miners

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/geomgr"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type MinerHandler struct {
	Conf viper.Viper
	MSer MinerService
}

func NewMinerHandler(conf *viper.Viper, mSer MinerService) *MinerHandler {
	return &MinerHandler{
		Conf: *conf,
		MSer: mSer,
	}
}

func BuildMinerHandlerInstance() *MinerHandler {
	conf := config.NewConfig()
	dbMgr, err := db.NewDatabaseMgrImpl(conf)
	if err != nil {
		panic("failed to connect database")
	}
	nodeURL := conf.GetString("FILECOIN_NODE_URL")
	fMgr, err := fmgr.NewFilecoinImpl(nodeURL)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}

	g := geomgr.NewGeoMgrImpl(conf)
	mSer := NewMinerServiceImpl(conf, dbMgr, fMgr, g)

	return NewMinerHandler(conf, mSer)
}

func (mHdl *MinerHandler) GetAllMiners() []*models.Miner {
	return mHdl.MSer.GetAllMiners()
}

func (mHdl *MinerHandler) MinersParseOffset(offset string) []*models.Miner {
	miners := make([]*models.Miner, 0)
	if strings.TrimSpace(offset) == "" {
		off := mHdl.Conf.GetInt("FILECOIN_BLOCKS_OFFSET")
		mHdl.MSer.ParseMinersByBlockOffset(off)
		return miners
	}
	off, err := strconv.Atoi(offset)
	if err != nil {
		log.Println("Error: provided offset is not a valid integer")
		return miners
	}
	miners = mHdl.MSer.ParseMinersByBlockOffset(off)
	return miners
}

func (mHdl *MinerHandler) MinersParseBlock(height int64) []*models.Miner {
	return mHdl.MSer.ParseMinersByBlockHeight(height)
}

func (mHdl *MinerHandler) MinersParseStateMarket() []*models.Miner {
	return mHdl.MSer.ParseMinersByStateMarket()
}
