package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/revel/revel"
)

var VisitorCount int64 = 0
var countLock sync.Mutex
var ReaderCounts = make(map[string]int)

func loadVisitorCount() {
	content, err := ioutil.ReadFile(BaseDir + "/visitors_count")
	if err != nil {
		revel.ERROR.Printf("Cannot load visitor count, %s", err.Error())
		VisitorCount = 0
	} else {
		VisitorCount, _ = strconv.ParseInt(string(content), 10, 64)
	}
}

func RecordVisit(c *revel.Controller) revel.Result {
	target := c.Request.URL.String()
	from := c.Request.Header.Get("X-Real-Ip")
	if from == "" {
		from = c.Request.RemoteAddr
	}
	go func(t string, f string) {
		today := time.Now().Format("Jan_02_2006")
		visitTime := time.Now().Format("15:04:05")
		countLock.Lock()
		defer countLock.Unlock()
		VisitorCount += 1
		msg := fmt.Sprintf("%d|%s|%s|%s\r\n", VisitorCount, t, f, visitTime)
		revel.INFO.Print(msg)
		if f, err := os.OpenFile(BaseDir+"/visitors_"+today,
			os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0644); err == nil {
			defer f.Close()
			if _, err = f.WriteString(msg); err != nil {
				revel.ERROR.Printf("Cannot log visitor info, %s\r\n", err.Error())
			}
		} else {
			revel.ERROR.Printf("Cannot log visitor info, %s\r\n", err.Error())
		}

		err := ioutil.WriteFile(BaseDir+"/visitors_count",
			[]byte(fmt.Sprintf("%d", VisitorCount)), 0644)
		if err != nil {
			revel.ERROR.Printf("Cannot update visitor info: %s", err.Error())
		}
	}(target, from)
	return nil
}
