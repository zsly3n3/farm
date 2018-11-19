package event

import (
	"farm/cache"
	"farm/conf"
	"farm/datastruct"
	"farm/db"
	"farm/tools"
	"time"
	//"farm/log"
)

type EventHandler struct {
	dbHandler     *db.DBHandler
	cacheHandler  *cache.CACHEHandler
	Version       string //当前服务端版本号
	ticker        *time.Ticker
	isExistTicker bool
	plants        []datastruct.Plant
	animals       map[datastruct.AnimalType]map[int]datastruct.Animal //按动物类型划分再根据动物编号 int为Number
	petbars       map[datastruct.AnimalType]datastruct.PetbarData
	soils         map[int]datastruct.SoilData //key为土地id
}

func CreateEventHandler() *EventHandler {
	eventHandler := new(EventHandler)
	eventHandler.Version = conf.Common.Version
	eventHandler.soils, eventHandler.petbars = tools.GetSoildInfo()
	eventHandler.cacheHandler = cache.CreateCACHEHandler()
	eventHandler.dbHandler = db.CreateDBHandler()
	eventHandler.createTicker(5 * time.Minute)
	eventHandler.plants = eventHandler.dbHandler.GetPlantsSlice()
	eventHandler.animals = eventHandler.dbHandler.GetAnimalsMap()
	return eventHandler
}
