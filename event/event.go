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
	plants map[int64]datastruct.Plant
	animals map[int64]datastruct.Animal
	petbars []datastruct.PetbarData
	soils []datastruct.SoilData
	
}

func CreateEventHandler()*EventHandler{
	 eventHandler:=new(EventHandler)
	 eventHandler.cacheHandler = cache.CreateCACHEHandler() 
	 eventHandler.dbHandler = db.CreateDBHandler() 
	 eventHandler.createTicker(5*time.Minute)
	 eventHandler.Version = "1.0.0.0"
	 eventHandler.plants = eventHandler.dbHandler.GetPlantsMap()
	 eventHandler.animals = eventHandler.dbHandler.GetAnimalsMap()
	 eventHandler.soils,eventHandler.petbars=tools.GetSoildInfo()
	 return eventHandler
}






