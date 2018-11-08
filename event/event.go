package event

import(
	"farm/db"
	"farm/cache"
	"time"
	"farm/datastruct"
	"farm/log"
)

type EventHandler struct {
	dbHandler *db.DBHandler
	cacheHandler *cache.CACHEHandler
	Version string //当前服务端版本号
	ticker *time.Ticker
	isExistTicker bool
	Plants map[int64]*datastruct.Plant 
}

func CreateEventHandler()*EventHandler{
	 eventHandler:=new(EventHandler)
	 eventHandler.cacheHandler = cache.CreateCACHEHandler() 
	 eventHandler.dbHandler = db.CreateDBHandler() 
	 eventHandler.createTicker(5*time.Minute)
	 eventHandler.Version = "1.0.0.0"
	 eventHandler.Plants = eventHandler.dbHandler.GetPlantsMap()
	 for k,v:=range eventHandler.Plants{
		 log.Debug("key:%v , value:%v",k,v)
	 }
	 return eventHandler
}






