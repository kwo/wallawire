package push

import (
	"context"
	"encoding/json"
	"time"

	"wallawire/logging"
	"wallawire/model"
)

func NewHeartbeatService(messageBus *PushMessenger, status *model.Status) *HeartbeatService {
	ctx, fnCancel := context.WithCancel(context.Background())
	return &HeartbeatService{
		ctx:        ctx,
		cancelFn:   fnCancel,
		messageBus: messageBus,
		status:     status,
	}
}

type HeartbeatService struct {
	ctx        context.Context
	cancelFn   context.CancelFunc
	messageBus *PushMessenger
	status     *model.Status
}

func (z *HeartbeatService) Start(interval time.Duration) {

	logger := logging.New(nil, "heartbeat")
	logger.Debug().Msg("starting...")

	ticker := time.NewTicker(interval)

Loop:
	for {
		select {
		case t := <-ticker.C:
			z.SendHeartbeat(t, "", "")
		case <-z.ctx.Done():
			break Loop
		}
	}

	ticker.Stop()
	logger.Debug().Msg("exiting")

}

func (z *HeartbeatService) Stop() {
	z.cancelFn()
}

func (z *HeartbeatService) SendHeartbeat(t time.Time, userID, sessionID string) {

	logger := logging.New(nil, "heartbeat")

	z.status.Populate(t.Truncate(time.Second))
	data, errData := json.Marshal(z.status)
	if errData != nil {
		logger.Warn().Err(errData).Msg("cannot serialize status")
		data = []byte("{}")
	}

	hb := model.PushMessage{
		Type: "heartbeat",
		Data: string(data),
	}

	count := z.messageBus.SendMessage(hb, userID, sessionID)
	logger.Debug().Int("count", count).Msg("heartbeat")

}
