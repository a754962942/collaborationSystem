package repo

import (
	"context"
	"github.com/a754962942/project-user/internal/data/member"
	"github.com/a754962942/project-user/internal/database"
)

type MemberRepo interface {
	GetMemberByEmail(ctx context.Context, email string) (bool, error)
	GetMemberByAccount(ctx context.Context, account string) (bool, error)
	GetMemberByMobile(ctx context.Context, mobile string) (bool, error)
	SaveMember(conn database.DbConn, ctx context.Context, mem *member.Member) error
	FindMember(ctx context.Context, account string, pwd string) (*member.Member, error)
}
