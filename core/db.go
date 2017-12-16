package core

import (
  "time"

  "github.com/pgmonzon/ServiciosYng/config"

  "gopkg.in/mgo.v2"
)

var mgoSession  *mgo.Session

func GetMongoSession() *mgo.Session {
  if mgoSession == nil {
    var err error

    mongoDBDialInfo := &mgo.DialInfo{
      Addrs:    []string{config.DB_Host},
      Timeout:  60 * time.Second,
      //Database: config.DB_Name,
      //Username: config.DB_User,
      //Password: config.DB_Pass,
    }

    mgoSession, err = mgo.DialWithInfo(mongoDBDialInfo)
    if err != nil {
      panic(err)
    }
    mgoSession.SetMode(mgo.Monotonic, true)
  }

  return mgoSession.Copy()
}
