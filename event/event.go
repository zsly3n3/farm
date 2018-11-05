package event

import(
	"farm/db"
	"farm/cache"
)

type EventHandler struct {
	dbHandler *db.DBHandler
    cacheHandler *cache.CACHEHandler
}

func CreateEventHandler()*EventHandler{
	 eventHandler:=new(EventHandler)
	 eventHandler.cacheHandler = cache.CreateCACHEHandler()
	 eventHandler.dbHandler = db.CreateDBHandler()
	 return eventHandler
}