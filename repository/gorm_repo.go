package repository

import (
	"errors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormUserRepository struct {
    db *gorm.DB
}

func NewGormUserRepository(dsn string) (*GormUserRepository, error) {
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    if err := db.AutoMigrate(&User{}); err != nil {
        return nil, err
    }

    return &GormUserRepository{db: db}, nil
}

func (r *GormUserRepository) Create(user User) (User, error) {
    if err := r.db.Create(&user).Error; err != nil {
        return User{}, err
    }
    return user, nil
}

func (r *GormUserRepository) FindByID(id int) (User, error) {
    var user User
    err := r.db.First(&user, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return User{}, ErrUserNotFound
        }
        return User{}, err
    }
    return user, nil
}

func (r *GormUserRepository) FindAll() ([]User, error) {
    var users []User
    if err := r.db.Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}

func (r *GormUserRepository) Update(user User) (User, error) {
    if err := r.db.Save(&user).Error; err != nil {
        return User{}, err
    }
    return user, nil
}

func (r *GormUserRepository) Delete(id int) error {
    result := r.db.Delete(&User{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return ErrUserNotFound
    }
    return nil
}
