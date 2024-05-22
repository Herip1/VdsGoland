package entity

type Vds struct {
	Vid    uint64 `json:"uid" gorm:"primarykey"`
	Uid    uint64 // список пользователей
	Name   string //название машины
	Img    string
	Ram    uint64
	Hdd    uint64
	Core   uint8
	Status bool //вкл/выкл
}
