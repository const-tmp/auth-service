package types

import (
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nullc4ts/bitmask_authz/access"
	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		Code       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();unique"`
		Name       string    `gorm:"size:32;unique" json:"name"`
		Password   []byte    `gorm:"size:60" json:"password"`
		TGName     string    `gorm:"size:128" json:"tg_name"`
		TGID       uint64    `gorm:"unique" json:"tg_id"`
		TGUserName string    `gorm:"size:32;unique" json:"tg_user_name"`
		//Access      access.Access `json:"access"`
		ParentID    uint `gorm:"default:null"`
		Parent      *User
		AccountID   uint
		Account     Account
		Blocked     bool
		Permissions []Permission `gorm:"many2many:user_permissions;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}
	Account struct {
		gorm.Model
		Code     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();unique"`
		Services []Service `gorm:"many2many:service_accounts;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		Name     string
	}
	Service struct {
		gorm.Model
		Name        string       `gorm:"unique"`
		Code        uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();unique"`
		Accounts    []Account    `gorm:"many2many:service_accounts;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		Permissions []Permission `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}
	Permission struct {
		gorm.Model
		ServiceID uint
		Name      string
		Access    access.Access
		Users     []User `gorm:"many2many:user_permissions;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}
	CustomClaims struct {
		jwt.RegisteredClaims
		Access access.Access `json:"access"`
	}
	AccessToken struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)