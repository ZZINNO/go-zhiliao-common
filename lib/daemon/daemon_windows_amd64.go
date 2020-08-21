// +build windows
// +build amd64

package daemon

import (
	"fmt"
	"net/http"
)

var DaemonCtl *Daemon

type Daemon struct {
	PidFile string
}

func InitDemon(port int, pidPerfix string) {
	return
}

func (d *Daemon) SigParmSwitch(sig string, daemon string) bool {
	return false
}

func (d *Daemon) WritePid() {
	return
}

func (d *Daemon) ZeroDown(srv http.Server) {
	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("http启动错误: %v\n", err)
	}
}
