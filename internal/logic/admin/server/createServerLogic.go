package server

import (
	"context"

	"github.com/perfect-panel/server/internal/model/dto"
	"github.com/perfect-panel/server/internal/model/entity/node"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/pkg/ip"
	"github.com/perfect-panel/server/pkg/logger"
	"github.com/perfect-panel/server/pkg/xerr"
	"github.com/pkg/errors"
)

type CreateServerLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateServerLogic Create Server
func NewCreateServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateServerLogic {
	return &CreateServerLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateServerLogic) CreateServer(req *dto.CreateServerRequest) error {
	data := node.Server{
		Name:    req.Name,
		Country: req.Country,
		City:    req.City,
		Address: req.Address,
		Sort:    req.Sort,
	}

	if data.City == "" && data.Country == "" {
		result, err := ip.GetRegionByIp(req.Address)
		if err != nil {
			l.Errorf("[CreateServer] GetRegionByIp Error: %v", err.Error())
		} else {
			data.City = result.City
			data.Country = result.Country
		}
	}

	err := l.svcCtx.Store.Node().InsertServer(l.ctx, &data)
	if err != nil {
		l.Errorf("[CreateServer] Insert Server error: %v", err.Error())
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "insert server error: %v", err)
	}
	return nil
}
