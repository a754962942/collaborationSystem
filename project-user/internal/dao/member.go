package dao

import (
	"context"
	"github.com/a754962942/project-user/internal/data/member"
	"github.com/a754962942/project-user/internal/database"
	"github.com/a754962942/project-user/internal/database/gorms"
)

type MemberDao struct {
	conn *gorms.GormConn
}

func NerMemberDao() *MemberDao {
	return &MemberDao{
		conn: gorms.New(),
	}
}
func (m *MemberDao) SaveMember(conn database.DbConn, ctx context.Context, mem *member.Member) error {
	m.conn = conn.(*gorms.GormConn)
	return m.conn.Tx(ctx).Create(mem).Error
}

func (m *MemberDao) GetMemberByAccount(ctx context.Context, account string) (bool, error) {
	var count int64

	err := m.conn.Default(ctx).Model(&member.Member{}).Where("account=?", account).Count(&count).Error
	return count > 0, err
}
func (m *MemberDao) GetMemberByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := m.conn.Default(ctx).Model(&member.Member{}).Where("email=?", email).Count(&count).Error
	return count > 0, err
}
func (m *MemberDao) GetMemberByMobile(ctx context.Context, mobile string) (bool, error) {
	var count int64
	err := m.conn.Default(ctx).Model(&member.Member{}).Where("mobile=?", mobile).Count(&count).Error
	return count > 0, err
}