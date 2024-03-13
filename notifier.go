package main

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"sincerecloud.com/openstack_notifier/global"
	"sincerecloud.com/openstack_notifier/log"
	"sincerecloud.com/openstack_notifier/services"
)


func init() {
	global.AssignCONF()
	//global.Viper()
	global.LOG = log.NewZap()
	zap.ReplaceGlobals(global.LOG)
}

func main() {
	services.RunServer()

	mux := http.NewServeMux()
	mux.HandleFunc("/actuator/health", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Health is ok.")
	})
	global.LOG.Info(fmt.Sprintf("Server is starting, port is %d.", global.CONF.App.HttpPort))
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", global.CONF.App.HttpPort), mux); err != nil {
			global.LOG.Error(fmt.Sprintf("Server status error，%s", err.Error()))
		}
	}()

	ExitServer()
}

func ExitServer() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		global.LOG.Info("Shutdown Server ...")
		for _, consumer := range services.Consumers {
			consumer.Close()
		}
	}
	global.LOG.Info("Server exited")
	os.Exit(0)
}
