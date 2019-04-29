package model

type Ip struct {
    ID    uint64    `json:"id,string" gorm:"column:id;type:bigint unsigned AUTO_INCREMENT;primary_key"`
    Ip    string    `json:"ip" gorm:"column:ip;type:varchar(20) not null default ''"`
    Port    string    `json:"port" gorm:"column:port;type:varchar(20) not null default ''"`
	Type    string    `json:"type" gorm:"column:type;type:varchar(20) not null default ''"`
	Addr	string  `json:"addr" gorm:"column:addr;type:varchar(50) not null default ''"`
}

func (t *Ip) TableName() string {
    return "ip"
}