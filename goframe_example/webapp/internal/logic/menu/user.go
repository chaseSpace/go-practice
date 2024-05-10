package menu

import (
	"context"
	"webapp/internal/model/menu"
	"webapp/internal/service"
)

type sUser struct {
}

func init() {
	// Register方法仅在 gen service 后可调用
	service.RegisterUser(&sUser{})
}

func (s *sUser) GetUser(ctx context.Context, in *menu.GetUserInput) (*menu.GetUserOutput, error) {
	panic("")
}
