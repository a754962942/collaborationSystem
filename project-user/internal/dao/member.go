package dao

import (
	"context"
	"github.com/a754962942/project-user/internal/data/member"
	"github.com/a754962942/project-user/internal/database"
	"github.com/a754962942/project-user/internal/database/gorms"
	"gorm.io/gorm"
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
func (m *MemberDao) FindMember(ctx context.Context, account string, pwd string) (*member.Member, error) {
	mem := &member.Member{}
	err := m.conn.Default(ctx).Where("account=? and password=? and status=1", account, pwd).First(&mem).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return mem, err
}
func (m *MemberDao) FindMemberById(ctx context.Context, id int64) (mem *member.Member, err error) {
	err = m.conn.Default(ctx).Where("id=?", id).Find(&mem).Error
	return
}
