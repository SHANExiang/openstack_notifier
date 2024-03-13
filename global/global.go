package global

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"sincerecloud.com/openstack_notifier/configs"
)


var (
	CONF          configs.Server
	TRANSPORT     *http.Transport
	VIPER         *viper.Viper
	LOG           *zap.Logger
)

