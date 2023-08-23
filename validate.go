package go_validate_user

import (
	"context"

	"github.com/bugfixes/go-bugfixes/logs"
	pb "github.com/todo-lists-app/protobufs/generated/id_checker/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IdCheckerServiceClientCreator interface {
	NewIdCheckerServiceClient() pb.IdCheckerServiceClient
}

type Validate struct {
	IdentityService string
	CTX             context.Context
	DevMode         bool
	Client          pb.IdCheckerServiceClient
}

type Checker interface {
	ValidateUser(accessToken, userId string) (bool, error)
}

func NewValidate(ctx context.Context, identityService string, devMode bool) *Validate {
	return &Validate{
		IdentityService: identityService,
		CTX:             ctx,
		DevMode:         devMode,
	}
}

func (v *Validate) GetClient() (*Validate, error) {
	conn, err := grpc.DialContext(v.CTX, v.IdentityService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, logs.Errorf("error dialing grpc: %v", err)
	}

	v.Client = pb.NewIdCheckerServiceClient(conn)
	return v, nil
}

func (v *Validate) ValidateUser(accessToken, userId string) (bool, error) {
	if v.DevMode {
		return true, nil
	}

	resp, err := v.Client.CheckId(v.CTX, &pb.CheckIdRequest{
		Id:          userId,
		AccessToken: accessToken,
	})
	if err != nil {
		_ = logs.Errorf("error checking id: %v, %s", err, userId)
		return false, nil
	}
	if !resp.GetIsValid() {
		return false, nil
	}

	return true, nil
}
