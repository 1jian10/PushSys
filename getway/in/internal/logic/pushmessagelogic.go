package logic

import (
	"context"

	"InnerGetWay/internal/svc"
	"InnerGetWay/proto/in"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushMessageLogic {
	return &PushMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PushMessageLogic) PushMessage(in *in.PushMessageReq) (*in.PushMessageResp, error) {
	// todo: add your logic here and delete this line

	return &in.PushMessageResp{}, nil
}
