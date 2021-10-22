package model

import "github.com/LittleBenx86/Benlog/internal/global/consts"

type Author struct {
	ID        uint      `json:"id" gorm:"primaryKey;<-;"` // natural key, meaningless
	Name      string    `json:"name" gorm:""`
	Authority Authority `json:"authority" gorm:"type:int;"`
	Role      string    `json:"-" gorm:"-"`
}

func (a *Author) UpdateRoleByAuthority() {
	switch a.Authority {
	case ADMINISTRATOR:
		a.Role = "admin"
	case METRICS_MEMBER:
		a.Role = "metrics_member"
	default:
		a.Role = consts.USER_ANONYMOUS_NAME
	}
}
