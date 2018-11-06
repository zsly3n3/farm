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
	 eventHandler.cacheHandler = cache.CreateCACHEHandler() //test
	 eventHandler.dbHandler = db.CreateDBHandler() //test
	 return eventHandler
}