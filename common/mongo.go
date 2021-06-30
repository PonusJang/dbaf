package common

import (
	"gopkg.in/mgo.v2"
)

var (
	mgoSession *mgo.Session
	dataBase   = "dbaf"
)

func GetSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial("10.10.8.6:27017")
		if err != nil {
			panic(err) //直接终止程序运行
		}
	}
	//最大连接池默认为4096
	return mgoSession.Clone()
}

//公共方法，获取collection对象

func WitchCollection(collection string, s func(*mgo.Collection) error) error {
	session := GetSession()
	defer session.Close()
	c := session.DB(dataBase).C(collection)
	return s(c)
}
