package datastorage

import (
	"bytes"
	"encoding/json"
)

// Announce the model we will use
// ip地址对应
type SprinkleRecord struct {
	ID             uint   `json:"id" gorm:"primaryKey; not null; auto_increment"`
	IPAddress      string `json:"ip_address" gorm:"type:varchar(20); not null" validate:"required" label:"访问者IP"`
	EtherumAddress string `json:"etherum_address" gorm:"type:varchar(50)" validate:"required,min=40,max=42" label:"以太坊地址"`
	Sprinkles      uint   `json:"sprinkles"`
	// 次数 {"xxx", "xxx", "xxx"}
	SprinkledAddrs string `json:"sprinkle_addrs" gorm:"type:varchar(200)"`
	CreateAt       int64  `json:"create_at" gorm:"autoCreateTime"`
}

// func (s *SprinkleRecord) SprinkledToString() (sprinkles string, err error) {

// 	return
// }

func (s *SprinkleRecord) SprinkledToSlice() (sprinkles []string, err error) {
	sprinkles = make([]string, 0)
	byts := bytes.NewBufferString(s.SprinkledAddrs)
	err = json.NewDecoder(byts).Decode(&sprinkles)

	return
}
