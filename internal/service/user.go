package service

import (
	"github.com/Confialink/wallet-permissions/internal/srvdiscovery"
	"context"
	"net/http"

	"github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/patrickmn/go-cache"
)

type User struct {
	cache *cache.Cache
}

func NewUser() *User {
	return &User{}
}

func (u *User) GetByUID(uid string) (*users.User, error) {
	req := users.Request{UID: uid}
	client, err := u.getClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.GetByUID(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return resp.User, nil
}

func (u *User) GetByClassId(classId int64) (result []*users.User, err error) {
	req := users.Request{ClassId: uint64(classId)}
	client, err := u.getClient()
	if err != nil {
		return
	}

	resp, err := client.GetByAdministratorClassId(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	result = resp.Users

	return
}

func (u *User) getClient() (users.UserHandler, error) {
	usersUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNameUsers)
	if nil != err {
		return nil, err
	}
	return users.NewUserHandlerProtobufClient(usersUrl.String(), http.DefaultClient), nil
}
