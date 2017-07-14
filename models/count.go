package models

type Count struct {
	Id      int64
	Key     string `xorm:"index"`
	Value   int64  `xorm:"index"`
	Version int64  `xorm:"version"` //乐观锁
}
