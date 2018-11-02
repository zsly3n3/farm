package datastruct

import (
)


const NULLSTRING = ""
const NULLID = -1

type TestData struct {
	 UserName string
	 Avatar string
}



type CodeType int //错误码
const (
	NULLError CodeType = iota //无错误
	ParamError//参数错误,数据为空或者类型不对等
	LoginFailed//登录失败,如无此账号或者密码错误等
	JsonParseFailedFromPostBody//来自post请求中的Body解析json失败
	DBSessionGetError//xorm事务中Get方法执行出错
	DBSessionExecError//xorm事务中Exec方法执行出错
	DBSessionInsertError//xorm事务中Insert方法执行出错
	DBSessionCommitError//xorm事务中Commit方法执行出错
	DBSessionUpdateError//xorm事务中Update方法执行出错
)

type Platform int //平台
const (
	WX_Platform Platform = iota //微信平台
    PC_Platform //pc平台
)

type TestTable struct {
	Id    int       `xorm:"not null pk autoincr INT(11)"`
	Level int       `xorm:"INT(11) not null"`  //权限等级
    Desc  string    `xorm:"VARCHAR(32) not null"` //权限名称
}

