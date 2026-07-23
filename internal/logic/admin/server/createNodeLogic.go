package server

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/perfect-panel/server/internal/model/dto"
	"github.com/perfect-panel/server/internal/model/entity/node"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/pkg/logger"
	"github.com/perfect-panel/server/pkg/tool"
	"github.com/perfect-panel/server/pkg/xerr"
	"github.com/pkg/errors"
)

type CreateNodeLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateNodeLogic Create Node
func NewCreateNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNodeLogic {
	return &CreateNodeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateNodeLogic) CreateNode(req *dto.CreateNodeRequest) error {
	data := node.Node{
		Name:     req.Name,
		Tags:     tool.StringSliceToString(req.Tags),
		Enabled:  req.Enabled,
		Port:     req.Port,
		Address:  req.Address,
		ServerId: req.ServerId,
		Protocol: req.Protocol,
	}

	// Protocol configuration
	var protocol node.Protocol
	protocolBytes, err := json.Marshal(req.ProtocolConfig)
	if err != nil {
		l.Errorf("[CreateNode] Marshal ProtocolConfig Error: %v", err.Error())
		return errors.Wrapf(xerr.NewErrCodeMsg(xerr.InvalidParams, "protocol config marshal error"), "protocol config marshal error: %v", err)
	}
	if err := json.Unmarshal(protocolBytes, &protocol); err != nil {
		l.Errorf("[CreateNode] Unmarshal to Node Protocol Error: %v", err.Error())
		return errors.Wrapf(xerr.NewErrCodeMsg(xerr.InvalidParams, "protocol config parse error"), "protocol config parse error: %v", err)
	}

	// If protocol type is not set, use the protocol string
	if protocol.Type == "" {
		protocol.Type = req.Protocol
	}

	// VLESS Reality Key Generation
	if protocol.Type == "vless" {
		if protocol.Security == "reality" {
			if protocol.RealityPublicKey == "" {
				public, private, err := tool.Curve25519Genkey(false, "")
				if err != nil {
					l.Errorf("[CreateNode] Generate Reality Key Error: %v", err.Error())
					return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "generate reality key error: %v", err)
				}
				protocol.RealityPublicKey = public
				protocol.RealityPrivateKey = private
				protocol.RealityShortId = tool.GenerateShortID(private)
			}
			if protocol.RealityServerAddr == "" {
				protocol.RealityServerAddr = protocol.SNI
			}
			if protocol.RealityServerPort == 0 {
				protocol.RealityServerPort = 443
			}
		}
	}

	// ShadowSocks 2022 Key Generation
	if protocol.Type == "shadowsocks" {
		if strings.Contains(protocol.Cipher, "2022") {
			var length int
			switch protocol.Cipher {
			case "2022-blake3-aes-128-gcm":
				length = 16
			default:
				length = 32
			}
			if len(protocol.ServerKey) != length {
				protocol.ServerKey = tool.GenerateCipher(protocol.ServerKey, length)
			}
		}
	}

	err = data.MarshalProtocolConfig(protocol)
	if err != nil {
		l.Errorf("[CreateNode] Marshal Protocol Config Error: %v", err.Error())
		return errors.Wrapf(xerr.NewErrCodeMsg(xerr.InvalidParams, "protocol config marshal error"), "protocol config marshal error: %v", err)
	}

	err = l.svcCtx.Store.Node().InsertNode(l.ctx, &data)
	if err != nil {
		l.Errorw("[CreateNode] Insert Database Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "[CreateNode] Insert Database Error")
	}

	return nil
}
