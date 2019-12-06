package server

import (
	"context"
	"gim/logic/model"
	"gim/logic/service"
	"gim/public/imerror"
	"gim/public/logger"
	"gim/public/pb"
	"gim/public/util"
)

type LogicClientExtServer struct{}

// RegisterDevice 注册设备
func (*LogicClientExtServer) RegisterDevice(ctx context.Context, in *pb.RegisterDeviceReq) (*pb.RegisterDeviceResp, error) {
	device := model.Device{
		AppId:         in.AppId,
		Type:          in.Type,
		Brand:         in.Brand,
		Model:         in.Model,
		SystemVersion: in.SystemVersion,
		SDKVersion:    in.SdkVersion,
	}

	if device.AppId == 0 || device.Type == 0 || device.Brand == "" || device.Model == "" ||
		device.SystemVersion == "" || device.SDKVersion == "" {
		return nil, imerror.ErrBadRequest
	}

	id, err := service.DeviceService.Register(Context(), device)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return &pb.RegisterDeviceResp{DeviceId: id}, nil
}

// AddUser 添加用户
func (*LogicClientExtServer) AddUser(ctx context.Context, in *pb.AddUserReq) (*pb.AddUserResp, error) {
	appId, userId, _, err := util.GetCtxData(ctx)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.AddUserResp{}, err
	}
	user := model.User{
		AppId:     appId,
		UserId:    userId,
		Nickname:  in.User.Nickname,
		Sex:       in.User.Sex,
		AvatarUrl: in.User.AvatarUrl,
		Extra:     in.User.Extra,
	}

	return &pb.AddUserResp{}, service.UserService.Add(Context(), user)
}

// GetUser 获取用户信息
func (*LogicClientExtServer) GetUser(ctx context.Context, in *pb.GetUserReq) (*pb.GetUserResp, error) {
	appId, _, _, err := util.GetCtxData(ctx)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.GetUserResp{}, err
	}

	user, err := service.UserService.Get(Context(), appId, in.UserId)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.GetUserResp{}, nil
	}

	if user == nil {
		return nil, imerror.ErrUserNotExist
	}

	pbUser := pb.User{
		UserId:     user.UserId,
		Nickname:   user.Nickname,
		Sex:        user.Sex,
		AvatarUrl:  user.AvatarUrl,
		Extra:      user.Extra,
		CreateTime: user.CreateTime.Unix(),
		UpdateTime: user.UpdateTime.Unix(),
	}
	return &pb.GetUserResp{User: &pbUser}, nil
}

// SendMessage 发送消息
func (*LogicClientExtServer) SendMessage(ctx context.Context, in *pb.SendMessageReq) (*pb.SendMessageResp, error) {
	appId, userId, deviceId, err := util.GetCtxData(ctx)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	err = service.MessageService.Send(Context(), appId, userId, deviceId, *in)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return &pb.SendMessageResp{}, nil
}

// CreateGroup 创建群组
func (*LogicClientExtServer) CreateGroup(ctx context.Context, in *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {
	appId, _, _, err := util.GetCtxData(ctx)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.CreateGroupResp{}, err
	}

	var group = model.Group{
		AppId:        appId,
		GroupId:      in.Group.GroupId,
		Name:         in.Group.Name,
		Introduction: in.Group.Introduction,
		Type:         in.Group.Type,
		Extra:        in.Group.Extra,
	}
	err = service.GroupService.Create(Context(), group)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.CreateGroupResp{}, err
	}
	return &pb.CreateGroupResp{}, nil
}

// UpdateGroup 更新群组
func (*LogicClientExtServer) UpdateGroup(ctx context.Context, in *pb.UpdateGroupReq) (*pb.UpdateGroupResp, error) {
	appId, _, _, err := util.GetCtxData(ctx)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.UpdateGroupResp{}, err
	}

	var group = model.Group{
		AppId:        appId,
		GroupId:      in.Group.GroupId,
		Name:         in.Group.Name,
		Introduction: in.Group.Introduction,
		Type:         in.Group.Type,
		Extra:        in.Group.Extra,
	}
	err = service.GroupService.Update(Context(), group)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.UpdateGroupResp{}, err
	}
	return &pb.UpdateGroupResp{}, nil
}

// GetGroup 获取群组信息
func (*LogicClientExtServer) GetGroup(ctx context.Context, in *pb.GetGroupReq) (*pb.GetGroupResp, error) {
	appId, _, _, err := util.GetCtxData(ctx)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.GetGroupResp{}, err
	}

	group, err := service.GroupService.Get(Context(), appId, in.GroupId)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.GetGroupResp{}, err
	}

	if group == nil {
		return nil, imerror.ErrGroupNotExist
	}

	return &pb.GetGroupResp{
		Group: &pb.Group{
			GroupId:      group.GroupId,
			Name:         group.Name,
			Introduction: group.Introduction,
			UserMum:      group.UserNum,
			Type:         group.Type,
			Extra:        group.Extra,
			CreateTime:   group.CreateTime.Unix(),
			UpdateTime:   group.UpdateTime.Unix(),
		},
	}, nil
}

// GetUserGroups 获取用户加入的所有群组
func (*LogicClientExtServer) GetUserGroups(ctx context.Context, in *pb.GetUserGroupsReq) (*pb.GetUserGroupsResp, error) {
	appId, userId, _, err := util.GetCtxData(ctx)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.GetUserGroupsResp{}, err
	}

	groups, err := service.GroupUserService.ListByUserId(Context(), appId, userId)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.GetUserGroupsResp{}, err
	}
	pbGroups := make([]*pb.Group, 0, len(groups))
	for i := range groups {
		pbGroups = append(pbGroups, &pb.Group{
			GroupId:      groups[i].GroupId,
			Name:         groups[i].Name,
			Introduction: groups[i].Introduction,
			UserMum:      groups[i].UserNum,
			Type:         groups[i].Type,
			Extra:        groups[i].Extra,
			CreateTime:   groups[i].CreateTime.Unix(),
			UpdateTime:   groups[i].UpdateTime.Unix(),
		})
	}
	return &pb.GetUserGroupsResp{Groups: pbGroups}, err
}

// AddGroupMember 添加群组成员
func (*LogicClientExtServer) AddGroupMember(ctx context.Context, in *pb.AddGroupMemberReq) (*pb.AddGroupMemberResp, error) {
	appId, _, _, err := util.GetCtxData(ctx)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.AddGroupMemberResp{}, err
	}

	err = service.GroupService.AddUser(Context(), appId, in.GroupUser.GroupId, in.GroupUser.UserId, in.GroupUser.Label, in.GroupUser.Extra)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.AddGroupMemberResp{}, err
	}

	return &pb.AddGroupMemberResp{}, nil
}

// UpdateGroupMember 更新群组成员信息
func (*LogicClientExtServer) UpdateGroupMember(ctx context.Context, in *pb.UpdateGroupMemberReq) (*pb.UpdateGroupMemberResp, error) {
	appId, _, _, err := util.GetCtxData(ctx)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.UpdateGroupMemberResp{}, err
	}

	err = service.GroupService.UpdateUser(Context(), appId, in.GroupUser.GroupId, in.GroupUser.UserId, in.GroupUser.Label, in.GroupUser.Extra)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.UpdateGroupMemberResp{}, err
	}

	return &pb.UpdateGroupMemberResp{}, nil
}

// DeleteGroupMember 添加群组成员
func (*LogicClientExtServer) DeleteGroupMember(ctx context.Context, in *pb.DeleteGroupMemberReq) (*pb.DeleteGroupMemberResp, error) {
	appId, _, _, err := util.GetCtxData(ctx)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.DeleteGroupMemberResp{}, err
	}

	err = service.GroupService.DeleteUser(Context(), appId, in.GroupId, in.UserId)
	if err != nil {
		logger.Sugar.Error(err)
		return &pb.DeleteGroupMemberResp{}, err
	}

	return &pb.DeleteGroupMemberResp{}, nil
}