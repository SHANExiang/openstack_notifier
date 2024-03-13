### 环境信息
1. go版本1.18；
2. openstack消费使用rabbitmq。

### 项目启动流程
1. 进入项目目录，执行 `go mod tidy` 更新下载项目所需包
2. 调整配置文件 `configs/app.yaml`下的ECS、zap、App等配置

### 开发流程
1. services层开发 `services`目录下添加对应实现业务文件

