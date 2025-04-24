package dao

import (
	"context"
	"database/sql"
	"errors"
	mysql2 "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = gorm.ErrDuplicatedKey
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (d *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := d.db.WithContext(ctx).Create(&u).Error
	// 判断返回错误是否是邮箱冲突，若是，就一层层传递到web层
	var mysqlErr *mysql2.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueIndexErrNo uint16 = 1062
		if mysqlErr.Number == uniqueIndexErrNo {
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (d *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (d *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}

func (d *UserDAO) Update(ctx context.Context, u User) error {
	return d.db.WithContext(ctx).Model(&u).
		Select("nickname", "birthday", "about_me", "utime").
		Where("id = ?", u.Id).
		Updates(map[string]interface{}{
			"nickname": u.Nickname,
			"birthday": u.Birthday,
			"about_me": u.AboutMe,
			"utime":    time.Now().UnixMilli(),
		}).Error
}

type User struct {
	Id       int64  `gorm:"primaryKey;autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	Birthday sql.NullInt64
	Nickname sql.NullString
	AboutMe  sql.NullString `gorm:"type:varchar(1024)"`
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}
