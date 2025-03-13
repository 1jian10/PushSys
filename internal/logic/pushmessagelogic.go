package logic

import (
	"context"
	"github.com/google/uuid"
	"puhser/internal/context"
	"puhser/proto/push"
	"puhser/route"
	"puhser/test"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.Context
	logx.Logger
}

func NewPushMessageLogic(ctx context.Context, svcCtx *svc.Context) *PushMessageLogic {
	return &PushMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PushMessageLogic) PushMessage(in *push.PushMessageReq) (*push.PushMessageResp, error) {
	uid := uuid.New().String()
	msg := &route.Message{
		UUId:       uid,
		PayLoad:    in.PayLoad,
		EncodeType: in.EncodeType,
	}
	for _, id := range in.UserId {
		value, ok := client.SMap.Load(strconv.FormatInt(id, 10))
		if !ok {
			continue
		}
		c := value.(*route.Client)
		msg.Session = c.Session
		c.Send(msg)
	}

	return &push.PushMessageResp{}, nil
}
