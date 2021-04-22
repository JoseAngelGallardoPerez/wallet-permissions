package rpc

import (
	"context"
	"log"
	"net/http"

	"github.com/twitchtv/twirp"

	"github.com/Confialink/wallet-permissions/internal/db/model"
	"github.com/Confialink/wallet-permissions/internal/service"
	pb "github.com/Confialink/wallet-permissions/rpc/permissions"
)

type Server struct {
	permissionsService *service.Permissions
}

func NewRPCServer(permissionsService *service.Permissions) *Server {
	return &Server{permissionsService: permissionsService}
}

func (s *Server) ListenAndServe(addr string) {
	twirpHandler := pb.NewPermissionCheckerServer(s, nil)
	mux := http.NewServeMux()
	mux.Handle(pb.PermissionCheckerPathPrefix, twirpHandler)
	go http.ListenAndServe(addr, mux)
	log.Printf("Listening and serving RPC on %s", addr)
}

//Check verifies if permission is granted
func (s *Server) Check(c context.Context, req *pb.PermissionReq) (*pb.PermissionResp, error) {
	perm, typedErr := s.permissionsService.Check(req.UserId, req.ActionKey)
	if nil != typedErr {
		return nil, typedErr
	}
	resp := &pb.PermissionResp{
		ActionKey: perm.ActionKey,
		UserId:    perm.UserId,
		IsAllowed: perm.IsAllowed,
	}
	return resp, nil
}

//CheckAll verifies if permissions are granted
func (s *Server) CheckAll(c context.Context, req *pb.PermissionsReq) (*pb.PermissionsResp, error) {
	perms, typedErr := s.permissionsService.CheckAll(req.UserId, req.ActionKeys)
	if nil != typedErr {
		return nil, typedErr
	}
	resp := &pb.PermissionsResp{Permissions: make([]*pb.PermissionResp, len(perms))}
	for i, perm := range perms {
		resp.Permissions[i] = &pb.PermissionResp{
			ActionKey: perm.ActionKey,
			UserId:    perm.UserId,
			IsAllowed: perm.IsAllowed,
		}
	}
	return resp, nil
}

// GetGroupsByIds returns groups by passed ids in request
func (s *Server) GetGroupsByIds(c context.Context, req *pb.GroupIdsReq) (*pb.GroupsResponse, error) {
	groups, err := s.permissionsService.GetGroupsByIds(req.Ids)
	if nil != err {
		return nil, twirp.InternalErrorWith(err)
	}
	resp := &pb.GroupsResponse{Groups: getPbGroups(groups)}
	return resp, nil
}

func getPbGroups(groups []*model.Group) []*pb.Group {
	pbGroups := make([]*pb.Group, len(groups))
	for i, v := range groups {
		pbGroups[i] = &pb.Group{
			Id:          v.GetId(),
			Name:        v.Name(),
			Description: v.Description(),
		}
	}
	return pbGroups
}
