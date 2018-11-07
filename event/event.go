package event

import(
	"farm/db"
	"farm/cache"
)

type EventHandler struct {
	dbHandler *db.DBHandler
	cacheHandler *cache.CACHEHandler
	Version string //当前服务端版本号
}

func CreateEventHandler()*EventHandler{
	 eventHandler:=new(EventHandler)
	 eventHandler.cacheHandler = cache.CreateCACHEHandler() 
	 eventHandler.dbHandler = db.CreateDBHandler() 
	 eventHandler.Version = "1.0"
	 return eventHandler
}