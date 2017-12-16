package core

import (
  "time"
  "net/http"
  "fmt"
  "strings"

  "github.com/pgmonzon/Yangee/config"

  "gopkg.in/mgo.v2"
)

var mgoSession  *mgo.Session

func GetMongoSession() (*mgo.Session, error, int) {
  if mgoSession == nil {
    var err error

    mongoDBDialInfo := &mgo.DialInfo{
      Addrs:    []string{config.DB_Host},
      Timeout:  config.DB_Timeout * time.Second,
      Database: config.DB_Name,
      //Username: config.DB_User,
      //Password: config.DB_Pass,
    }

    mgoSession, err = mgo.DialWithInfo(mongoDBDialInfo)
    if err != nil {
      s := []string{"INTERNAL_SERVER_ERROR:", err.Error()}
      return nil, fmt.Errorf(strings.Join(s, " ")), http.StatusInternalServerError
    }
    mgoSession.SetMode(mgo.Monotonic, true)
  }

  return mgoSession.Copy(), nil, http.StatusOK
}
