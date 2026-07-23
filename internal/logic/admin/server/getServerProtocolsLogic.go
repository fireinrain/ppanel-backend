package server

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/internal/model/dto"
	"github.com/perfect-panel/server/internal/model/entity/node"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/pkg/logger"
	"github.com/perfect-panel/server/pkg/xerr"
	"github.com/pkg/errors"
)

type GetServerProtocolsLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get Server Protocols
func NewGetServerProtocolsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServerProtocolsLogic {
	return &GetServerProtocolsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetServerProtocolsLogic) GetServerProtocols(req *dto.GetServerProtocolsRequest) (resp *dto.GetServerProtocolsResponse, err error) {
	// find all nodes of this server
	_, nodes, err := l.svcCtx.Store.Node().FilterNodeList(l.ctx, &node.FilterNodeParams{
		Page:     1,
		Size:     1000,
		ServerId: []int64{req.Id},
		Search:   "",
	})
	if err != nil {
		l.Errorf("[GetServerProtocols] FilterNodeList Error: %s", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "[GetServerProtocols] FilterNodeList Error: %s", err.Error())
	}

	// aggregate protocols from nodes
	var protocols []dto.Protocol
	for _, nodeItem := range nodes {
		protocolConfig, err := nodeItem.UnmarshalProtocolConfig()
		if err != nil {
			l.Errorf("[GetServerProtocols] UnmarshalProtocolConfig Error: %s", err.Error())
			continue
		}
		if protocolConfig.Type != "" {
			var dtoProtocol dto.Protocol
			configBytes, marshalErr := json.Marshal(protocolConfig)
			if marshalErr == nil {
				_ = json.Unmarshal(configBytes, &dtoProtocol)
			}
			protocols = append(protocols, dtoProtocol)
		}
	}

	return &dto.GetServerProtocolsResponse{
		Protocols: protocols,
	}, nil
}
