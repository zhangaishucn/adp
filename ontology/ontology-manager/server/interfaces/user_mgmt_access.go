package interfaces

import (
	"context"
)

// UserMgmtAccess 定义用户管理相关的访问接口
type UserMgmtAccess interface {
	GetAccountNames(ctx context.Context, accountInfos []*AccountInfo) error
}
