package event

import(
	"farm/db"
	"farm/cache"
	"time"
)

type EventHandler struct {
	dbHandler *db.DBHandler
	cacheHandler *cache.CACHEHandler
	Version string //当前服务端版本号
	ticker *time.Ticker
	isExistTicker bool
}

func CreateEventHandler()*EventHandler{
	 eventHandler:=new(EventHandler)
	 eventHandler.cacheHandler = cache.CreateCACHEHandler() 
	 eventHandler.dbHandler = db.CreateDBHandler() 
	 eventHandler.createTicker(5*time.Minute)
	 eventHandler.Version = "1.0.0.0"
	 return eventHandler
}






