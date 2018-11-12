package event

import(
	"farm/db"
	"farm/cache"
	"time"
	"farm/datastruct"
	"farm/tools"
	"farm/log"
)

type EventHandler struct {
	dbHandler *db.DBHandler
	cacheHandler *cache.CACHEHandler
	Version string //当前服务端版本号
	ticker *time.Ticker
	isExistTicker bool
	plants []datastruct.Plant
	animals map[datastruct.AnimalType]map[int]datastruct.Animal//按动物类型划分
	petbars []datastruct.PetbarData
	soils []datastruct.SoilData
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
	 for k,v := range eventHandler.animals{
		 log.Debug("k:%v , v:%v",k,v)
	 }
	 return eventHandler
}







