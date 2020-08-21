package logs

import (
	"fmt"
	"github.com/getsentry/raven-go"
)

func InitSentry(dsn string) {
	err := raven.SetDSN(dsn)
	if err != nil {
		fmt.Print("连接sentry失败,", err)
	}
}
