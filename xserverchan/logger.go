package xserverchan

import (
	"fmt"
	"github.com/Aoi-hosizora/go-serverchan"
	"github.com/sirupsen/logrus"
)

type ServerchanLogrus struct {
	logger  *logrus.Logger
	logMode bool
}

func NewServerchanLogrus(logger *logrus.Logger, logMode bool) *ServerchanLogrus {
	return &ServerchanLogrus{logger: logger, logMode: logMode}
}

func (s *ServerchanLogrus) Log(sckey string, title string, code int32, err error) {
	if !s.logMode {
		return
	}

	sckey = serverchan.Mask(sckey)
	title = serverchan.Mask(title)

	if err != nil {
		if !serverchan.IsResponseError(err) {
			s.logger.Error(fmt.Sprintf("[Serverchan] Send to %s error: %v", sckey, err))
		} else {
			s.logger.WithFields(map[string]interface{}{
				"module":    "serverchan",
				"sckeyMask": sckey,
				"code":      code,
			}).Error(fmt.Sprintf("[Serverchan] Send to %s error: %v", sckey, err))
		}
	} else {
		s.logger.WithFields(map[string]interface{}{
			"module":    "serverchan",
			"sckeyMask": sckey,
			"titleMask": title,
			"code":      0,
		}).Info(fmt.Sprintf("[Serverchan] < | %s | %s", sckey, title))
	}
}
