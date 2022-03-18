package log

import (
    "github.com/sirupsen/logrus"
    "testing"
    "time"
)

func Test_Log(t *testing.T) {
    Init(logrus.TraceLevel)

    //Log.Log(logrus.InfoLevel, "233")
    Log.Log(logrus.DebugLevel, "hehe")
    Log.Infof("%v hehe", "dabing")
    time.Sleep(2 * time.Second)
}
