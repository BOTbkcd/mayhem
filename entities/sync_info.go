package entities

import (
	"gorm.io/gorm"
)

type SyncInfo struct {
	gorm.Model
	Key      string
	Token    string
	BoardURL string
}

func FetchSyncInfo() []SyncInfo {
	var info []SyncInfo
	DB.Model(&SyncInfo{}).Find(&info)

	if len(info) == 0 {
		info := SyncInfo{}
		return []SyncInfo{info}
	}

	return info
}

func (s SyncInfo) Save() SyncInfo {
	DB.Save(&s)
	return s
}
