// +build linux
// +build amd64

package daemon

import (
	"context"
	"fmt"
	"github.com/ZZINNO/go-zhiliao-common/util"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var DaemonCtl *Daemon
var PidLock *FileLock

type Daemon struct {
	PidFile  string
	GlaceFul int
}

func InitDemon(port int, pidPerfix string, glaceful int) {
	DaemonCtl = &Daemon{
		PidFile:  fmt.Sprintf(pidPerfix+"_%d"+".pid", port),
		GlaceFul: glaceful,
	}
}

func (d *Daemon) createPid(pid string) bool {
	_, err := os.OpenFile(d.PidFile, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Printf("PID文件创建失败:%v \n", err)
		return false
	}
	return true
}

func (d *Daemon) pidStatus(pid string) bool {
	if d.createPid(pid) != true {
		return false
	}
	pidlock := NewFileLock(pid)
	//下面是为文件加锁
	if err := pidlock.Lock(); err != nil {
		fmt.Printf("pid测试加锁失败，文件仍被占用:%v \n", err)
		return false
	}
	if err := pidlock.Unlock(); err != nil {
		return false
	}
	return true
}

//通过Fork执行自己
func (d *Daemon) ForkRun() error {
	if d.pidStatus(d.PidFile) != true {
		return errors.New("加锁失败")
	}
	// 将命令行参数中执行文件路径转换成可用路径
	filePath, _ := filepath.Abs(os.Args[0])
	cmd := exec.Command(filePath, os.Args[1:]...)
	// 将其他命令传入生成出的进程
	cmd.Stdin = os.Stdin // 给新进程设置文件描述符，可以重定向到文件中
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//进程组ID
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd.Start() // 开始执行新进程，不等待新进程退出
}

func (d *Daemon) SigParmSwitch(sig string, daemon string) bool {
	switch sig {
	case "start":
		{
			if daemon == "true" {
				//如果是daemon == true 那么就禁用stdout
				if os.Getppid() != 1 {
					if err := d.ForkRun(); err != nil {
						return true
					} else {
						os.Exit(-1)
					}
				}
			}
		}
	case "reload":
		{
			d.SignalProcess(syscall.SIGUSR2)
			return true
		}
	case "stop":
		{
			d.SignalProcess(syscall.SIGKILL)
			return true
		}
	}
	return false
}

//处理信号
func (d *Daemon) SignalProcess(signal os.Signal) {
	b, err := ioutil.ReadFile(d.PidFile)
	if err != nil {
		fmt.Println("pid 文件不存在")
		os.Exit(-1)
	}
	//选取PID
	p, err := os.FindProcess(util.S(string(b)).Int())
	if err != nil {
		fmt.Println("进程不存在")
		os.Exit(-1)
	}
	p.Signal(signal)
}

func (d *Daemon) WritePid() {
	//创建文件
	if d.createPid(d.PidFile) != true {
		os.Exit(-1)
	}
	//创建PID锁
	PidLock = NewFileLock(d.PidFile)
	//上锁
	if err := PidLock.Lock(); err != nil {
		fmt.Printf("PID加锁失败:%v", err)
		os.Exit(-1)
	}
	fmt.Printf("程序真实PID: %d \n", syscall.Getpid())
	if err := ioutil.WriteFile(d.PidFile, []byte(fmt.Sprintf("%d", syscall.Getpid())), os.ModePerm); err != nil {
		fmt.Printf("PID文件写入错误:%v", err)
		os.Exit(-1)
	}
}

//处理Http
func (d *Daemon) ZeroDown(srv http.Server) {
	sleepTime := time.Second * time.Duration(d.GlaceFul)
	sig := make(chan os.Signal)
	d.WritePid()
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			fmt.Printf("gin服务错误:%v", err)
		}
	}()
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGUSR2)
	for {
		select {
		case s := <-sig:
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT:
				fmt.Println("收到退出信号:%s", s.String())
				os.Exit(0)
			case syscall.SIGUSR2:
				fmt.Println("收到重载信号:%s", s.String())
				ctx, cancel := context.WithTimeout(context.Background(), sleepTime)
				//解锁
				if err := PidLock.Unlock(); err != nil {
					fmt.Printf("解锁失败,信号不被执行: %v \n", err)
					return
				}
				if err := d.ForkRun(); err != nil {
					fmt.Printf("重启服务失败: %v \n", err)
					break
				}
				go func(c context.Context, cancel context.CancelFunc) {
					defer cancel()
					if err := srv.Shutdown(c); err != nil {
						fmt.Printf("关闭出错:%v\n", err)
					}
				}(ctx, cancel)
				//关闭服务的上下文处理
				//等待300s
				time.Sleep(sleepTime)
				fmt.Println("处理完成")
				<-ctx.Done()
				os.Exit(-1)
				return
			}
		}
	}
}
