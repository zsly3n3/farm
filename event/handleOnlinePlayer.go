package event

import(
	"time"
)

func (handle *EventHandler)createTicker(times time.Duration){
    if !handle.isExistTicker {
		handle.isExistTicker = true
		handle.ticker = time.NewTicker(times)
        go handle.selectTicker()
    }
}

// func (handle *EventHandler)stopTicker(){
//     if handle.ticker != nil{
// 	   handle.ticker.Stop() 
// 	   handle.isExistTicker = false
//     }
// }

func (handle *EventHandler)selectTicker(){
    for{
        select {
         case <-handle.ticker.C:
            handle.checkOnlinePlayer()
        }
    }
}

func (handle *EventHandler)checkOnlinePlayer(){

}