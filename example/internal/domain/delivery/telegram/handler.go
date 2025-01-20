package telegram

import (
	cancelCmd "github.com/woyow/setupper/example/internal/domain/delivery/telegram/cancel-cmd"
	defaultHandler "github.com/woyow/setupper/example/internal/domain/delivery/telegram/common"

	menuCmd "github.com/woyow/setupper/example/internal/domain/delivery/telegram/menu-cmd"
	startCmd "github.com/woyow/setupper/example/internal/domain/delivery/telegram/start-cmd"
	
	"time"

	"github.com/woyow/setupper/example/internal/domain/domain/service"
	
	"github.com/woyow/setupper/pkg/telegram"

	"github.com/sirupsen/logrus"
)

const (
	domainLoggingKey   = "domain"
	domainLoggingValue = "your-tg"
)

func InitHandler(tg *telegram.Telegram, service *service.Service, log *logrus.Logger) {

	{
		cancelCmd.InitHandler(tg, nil, log)
		menuCmd.InitHandler(tg, service.MenuCmd, log)
		defaultHandler.InitHandler(tg, nil, log)
		startCmd.InitHandler(tg, nil, log)
	}

	go func() {
		for {
			if err := tg.Run(); err != nil {
				log.WithFields(logrus.Fields{
					domainLoggingKey: domainLoggingValue,
				}).Error("InitHandler - tg.Run error: ", err.Error())
			}
			<-time.After(1 * time.Second)
		}
	}()
}
