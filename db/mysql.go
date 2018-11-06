package db

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"farm/datastruct"
	"farm/log"
	_ "github.com/go-sql-driver/mysql"
)

const DB_IP = "localhost:3306"
const DB_Name = "farm"
const DB_UserName = "root"
const DB_Pwd = "Zsly3n@s"

type DBHandler struct {
	 mysqlEngine *xorm.Engine
}

func CreateDBHandler()*DBHandler{
	dbHandler:=new(DBHandler)
    dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",DB_UserName,DB_Pwd,DB_IP,DB_Name)
	engine, err:= xorm.NewEngine("mysql", dsn)
	errhandle(err)
	err=engine.Ping()
	errhandle(err)
	//日志打印SQL
    engine.ShowSQL(true)
	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(1)
    resetDB(engine)
    initData(engine)
	dbHandler.mysqlEngine = engine
    return dbHandler
}

func resetDB(engine *xorm.Engine){
	user:=&datastruct.UserInfo{}
	player:=&datastruct.PlayerInfo{}
	perm:=&datastruct.Permission{}
	err:=engine.DropTables(user,player,perm)
    errhandle(err)
	err=engine.CreateTables(user,player,perm)
    errhandle(err)
}

func initData(engine *xorm.Engine){
	data:=createPermissionData()
	_, err := engine.Insert(&data)
	errhandle(err)
}

func createPermissionData()[]datastruct.Permission{
	a:= datastruct.Permission{
		Name:"游客",
	}
	b:= datastruct.Permission{
		Name:"普通玩家",
	}
	// c:= datastruct.Permission{
	// 	Name:"会员",
	// }
	return []datastruct.Permission{a,b}
}

func errhandle(err error){
	if err != nil {
		log.Fatal("db error is %v", err.Error())
	}
}