package db

import (
	"farm/conf"
	"farm/datastruct"
	"farm/log"
	"farm/tools"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/robfig/cron"
)

type DBHandler struct {
	mysqlEngine *xorm.Engine
}

func CreateDBHandler() *DBHandler {
	dbHandler := new(DBHandler)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", conf.Server.DB_UserName, conf.Server.DB_Pwd, conf.Server.DB_IP, conf.Server.DB_Name)
	engine, err := xorm.NewEngine("mysql", dsn)
	errhandle(err)
	err = engine.Ping()
	errhandle(err)
	//日志打印SQL
	engine.ShowSQL(true)
	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(1)
	resetDB(engine)
	initData(engine)
	dbHandler.mysqlEngine = engine
	go timerTask(engine)
	return dbHandler
}

func resetDB(engine *xorm.Engine) {
	user := &datastruct.UserInfo{}
	player := &datastruct.PlayerInfo{}
	perm := &datastruct.Permission{}
	plants := &datastruct.Plant{}
	plantClass := &datastruct.PlantClass{}
	animal := &datastruct.Animal{}
	animalClass := &datastruct.AnimalClass{}

	soil1 := &datastruct.Soil1{}
	soil2 := &datastruct.Soil2{}
	soil3 := &datastruct.Soil3{}
	soil4 := &datastruct.Soil4{}
	soil5 := &datastruct.Soil5{}
	petbar1 := &datastruct.Petbar1{}
	petbar2 := &datastruct.Petbar2{}
	petbar3 := &datastruct.Petbar3{}
	petbar4 := &datastruct.Petbar4{}

	speedup := &datastruct.PlayerSpeedUp{}
	rewardStamina := &datastruct.RewardStamina{}
	inviteInfo := &datastruct.InviteInfo{}

	var err error
	if conf.Common.Mode == conf.Debug {
		err = engine.DropTables(user, player, perm, plants, plantClass, animal, animalClass, soil1, soil2, soil3, soil4, soil5, petbar1, petbar2, petbar3, petbar4, speedup, rewardStamina, inviteInfo)
		errhandle(err)
	}
	err = engine.CreateTables(user, player, perm, plants, plantClass, animal, animalClass, soil1, soil2, soil3, soil4, soil5, petbar1, petbar2, petbar3, petbar4, speedup, rewardStamina, inviteInfo)
	errhandle(err)
}

func initData(engine *xorm.Engine) {
	execStr := fmt.Sprintf("ALTER DATABASE %s CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;", conf.Server.DB_Name)
	_, err := engine.Exec(execStr)
	errhandle(err)
	_, err = engine.Exec("ALTER TABLE user_info CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	errhandle(err)
	_, err = engine.Exec("ALTER TABLE user_info CHANGE nick_name nick_name VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	errhandle(err)
	createPermissionData(engine)
	createPlantClass(engine)
	createPlant(engine)
	createAnimalClass(engine)
	createAnimal(engine)
}

func createAnimalClass(engine *xorm.Engine) {
	a := datastruct.AnimalClass{
		Id:   datastruct.Sea,
		Desc: "海",
	}
	b := datastruct.AnimalClass{
		Id:   datastruct.Land,
		Desc: "陆",
	}
	c := datastruct.AnimalClass{
		Id:   datastruct.Fly,
		Desc: "空",
	}
	d := datastruct.AnimalClass{
		Id:   datastruct.Deity,
		Desc: "神",
	}
	data := make([]datastruct.AnimalClass, 0)
	data = append(data, a)
	data = append(data, b)
	data = append(data, c)
	data = append(data, d)
	_, err := engine.Insert(&data)
	errhandle(err)
}

func createPlantClass(engine *xorm.Engine) {
	a := datastruct.PlantClass{
		Desc: "普通类植物",
	}
	b := datastruct.PlantClass{
		Desc: "仙类植物",
	}
	data := make([]datastruct.PlantClass, 0)
	data = append(data, a)
	data = append(data, b)
	_, err := engine.Insert(&data)
	errhandle(err)
}

func createPlant(engine *xorm.Engine) {
	data := tools.GetPlantsInfo()
	_, err := engine.Insert(&data)
	errhandle(err)
}

func createAnimal(engine *xorm.Engine) {
	data := tools.GetAnimalInfo()
	_, err := engine.Insert(&data)
	errhandle(err)
}

func createPermissionData(engine *xorm.Engine) {
	a := datastruct.Permission{
		Name: "游客",
	}
	b := datastruct.Permission{
		Name: "普通玩家",
	}
	// c:= datastruct.Permission{
	// 	Name:"会员",
	// }
	data := make([]datastruct.Permission, 0)
	data = append(data, a)
	data = append(data, b)
	_, err := engine.Insert(&data)
	errhandle(err)
}

func errhandle(err error) {
	if err != nil {
		log.Fatal("db error is %v", err.Error())
	}
}

/*定时任务*/
func timerTask(engine *xorm.Engine) {
	c := cron.New()
	spec := "0 0 0 * * ?"
	c.AddFunc(spec, func() {
		table0 := "reward_stamina"
		truncate := "truncate table "
		engine.Exec(truncate + table0)
	})
	c.Start()
}
