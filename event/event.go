package event

import(
	"farm/db"
	"farm/cache"
	"time"
	"farm/datastruct"
	"farm/tools"
	//"farm/log"
)

type EventHandler struct {
	dbHandler *db.DBHandler
	cacheHandler *cache.CACHEHandler
	Version string //当前服务端版本号
	ticker *time.Ticker
	isExistTicker bool
	plants []datastruct.Plant
	animals map[datastruct.AnimalType]map[int]datastruct.Animal//按动物类型划分
	petbars map[datastruct.AnimalType]datastruct.PetbarData
	soils map[int]datastruct.SoilData//key为土地id
}

func CreateEventHandler()*EventHandler{
	 eventHandler:=new(EventHandler)
	 eventHandler.Version = "1.0.0.0"
	 eventHandler.soils,eventHandler.petbars=tools.GetSoildInfo()
	 eventHandler.cacheHandler = cache.CreateCACHEHandler() 
	 eventHandler.dbHandler = db.CreateDBHandler() 
	 eventHandler.createTicker(5*time.Minute)
	 eventHandler.plants = eventHandler.dbHandler.GetPlantsSlice()
	 eventHandler.animals = eventHandler.dbHandler.GetAnimalsMap()
	 return eventHandler
}







