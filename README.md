# go 公共方法

go_common 提取了golang开发常用的方法和业务逻辑

### 方法列表

-  [配置类(go_common/config)](/config)
    - viper读取json配置,根据配置文件触发变更事件回调函数
    - go-redis封装,提供了redis常用操作以及初始化方法
- [日志类(go_common/logs)](/logs)
    - logrus的封装
    - ranven-go的封装(sentry)
- [工具类(go_common/util)](/util)
    - json和map以及string互转工具
    - string和各种类型互转工具
    - map深度复制
    - 删除指定的slice索引
-  [中间件(go_common/middleware)](/middleware)
    -  post和put表单提交去重
    -  casbin校验api路由权限
    -  csrf 中间件
    -  jwt中间件
    -  流控中间件

### 安装

#### 使用git submodule

##### 1. 先去需要引入这个包的项目里面执行

```bash
git submodule add https://github.com/ZZINNO/go-zhiliao-common.git common
# common 是克隆下来在项目的路径,可以是其他名字或者路径
```

##### 2. 然后去项目路径编辑 `go.mod`

```bash
go.mod
#里面加入
require github.com/ZZINNO/go-zhiliao-common  v0.0.0-incompatible
replace github.com/ZZINNO/go-zhiliao-common => ./common
```

#### 使用gomod & go get

##### 1. 配置环境变量

```bash
GOPROXY=https://goproxy.cn,direct
```

##### 2. go get

```bash
#go get拉取
go get github.com/ZZINNO/go-zhiliao-common@master
#master可以是指定的git commit hash
```



### 引用示例

对于gitsubmodule方法和gomod 方法都是一样

 ```go

import (
	"fmt"
	"github.com/ZZINNO/go-zhiliao-common/config"
	"github.com/ZZINNO/go-zhiliao-common/logs"
)
}
 ```

### 更新

如果git仓库，维护者更新了库，项目需要同步最新的库代码

#### 使用git submodule

##### 1. 项目路径执行submodule update

```bash
git submodule update --remote
```

#### 使用gomod & go get(需要上面提到的环境变量)

```bash
#更新到最新版
go get -u github.com/ZZINNO/go-zhiliao-common
#如果要更新到指定版本
go get -u github.com/ZZINNO/go-zhiliao-common@{版本具体的commit hash}
```



### 采用gitsubmodule 和 gomod的区别

采用gitsubmodule 会在本地留一个文件夹，然后包引入也是采用本地包引入的方法(确保common在项目路径内)

采用gomod方法，包会在 `$GOPATH/pkg/` 里面

**如果目标环境是和git仓库无法通信，只能使用 gitsubmodule  ，然后把项目文件推送到服务器编译**