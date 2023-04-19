package main

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/duiniwukenaihe/king-kf/router"
	"github.com/duiniwukenaihe/king-utils/common"
	"github.com/duiniwukenaihe/king-utils/common/log"
	"github.com/duiniwukenaihe/king-utils/common/rabbitmq"
	"github.com/duiniwukenaihe/king-utils/config"
	"github.com/duiniwukenaihe/king-utils/kit"
	_ "github.com/duiniwukenaihe/king-utils/middleware/Validator"
)

func main() {
	// Debug Mode
	gin.SetMode(config.Mode)
	g := gin.New()
	// 设置路由
	r := router.SetupRouter(kit.EnhanceGin(g))
	ginpprof.Wrap(r)
	// 通过消息中间件获取更新kubeConfig文件消息
	consumer := rabbitmq.Consumer{
		Address:      config.RabbitMQURL,
		ExchangeName: common.UpdateKubeConfig,
		Handler:      &rabbitmq.UpdateKubeConfig{},
	}
	go consumer.Run()
	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(config.Listen); err != nil {
		log.Fatalf("Listen error: %v", err)
	}
}
