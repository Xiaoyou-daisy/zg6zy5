package srvhandler

import (
	"context"
	"fmt"
	"zg6zy5/apiway/inits"
	"zg6zy5/service/models"
	__ "zg6zy5/service/protobuf"
)

type Server struct {
	__.UnimplementedUserServer
}

func (s *Server) Login(_ context.Context, req *__.LoginRequest) (*__.LoginResponse, error) {

	//查询手机号存不存在
	var user models.Users
	err := inits.DB.Where("mobile LIKE ?", req.Mobile).First(&user).Error

	if err != nil {
		fmt.Errorf("手机号未验证")
	}

	if req.Status == 1 {
		get := inits.Client.Get(context.Background(), req.Mobile).Val()
		if get == "" {
			return &__.LoginResponse{Msg: "短信验证失败"}, nil
		}
	} else {
		if user.Password != req.Password {
			return &__.LoginResponse{Msg: "密码错误"}, nil
		}
	}

	return &__.LoginResponse{Msg: "登录成功"}, nil
}
