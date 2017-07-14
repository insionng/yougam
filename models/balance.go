package models

import (
	"errors"
	"fmt"
	"time"
)

//时间 类型 数额 余额 描述
//Time Ctype Amount Balance Description
type Balance struct {
	Id          int64  `xorm:"index"`
	Uid         int64  `xorm:"index"`   //钱包拥有者
	Time        int64  `xorm:"index"`   //发生时间 为保证时间的正确性 此处使用时间戳
	Ctype       int64  `xorm:"index"`   //动作类型
	Amount      int64  `xorm:"index"`   //数额
	Balance     int64  `xorm:"index"`   //余额
	Description string `xorm:"text"`    //描述
	Version     int64  `xorm:"version"` //乐观锁
}

type Balancejuser struct {
	Balance `xorm:"extends"`
	User    `xorm:"extends"`
}

func GetBalancesByUidJoinUser(uid int64, ctype int64, offset int, limit int, field string) *[]*Balancejuser {
	var balance = new([]*Balancejuser)
	if field == "asc" {

		if uid == 0 { //uid为0则查询所有用户
			if ctype != 0 { //查询特定ctype
				Engine.Table("balance").Where("balance.ctype=?", ctype).Limit(limit, offset).Asc("balance.id").Join("LEFT", "user", "user.id = balance.uid").Find(balance)
			} else { //不限制ctype
				Engine.Table("balance").Limit(limit, offset).Asc("balance.id").Join("LEFT", "user", "user.id = balance.uid").Find(balance)
			}

		} else { //查询特定uid

			if ctype == 0 { //不限制ctype
				Engine.Table("balance").Where("balance.uid=?", uid).Limit(limit, offset).Asc("balance.id").Join("LEFT", "user", "user.id = balance.uid").Find(balance)

			} else { //查询特定ctype

				Engine.Table("balance").Where("balance.ctype=? and balance.uid=?", ctype, uid).Limit(limit, offset).Asc("balance.id").Join("LEFT", "user", "user.id = balance.uid").Find(balance)
			}
		}
	} else if field == "desc" {

		if uid == 0 {
			if ctype != 0 {
				Engine.Table("balance").Where("balance.ctype=?", ctype).Limit(limit, offset).Desc("balance."+field).Join("LEFT", "user", "user.id = balance.uid").Find(balance)
			} else {
				Engine.Table("balance").Limit(limit, offset).Desc("balance."+field).Join("LEFT", "user", "user.id = balance.uid").Find(balance)
			}

		} else {

			if ctype == 0 {
				Engine.Table("balance").Where("balance.uid=?", uid).Limit(limit, offset).Desc("balance."+field).Join("LEFT", "user", "user.id = balance.uid").Find(balance)

			} else {

				Engine.Table("balance").Where("balance.ctype=? and balance.uid=?", ctype, uid).Limit(limit, offset).Desc("balance."+field).Join("LEFT", "user", "user.id = balance.uid").Find(balance)
			}
		}
	}
	return balance
}

func GetBalancesByUid(uid int64, ctype int64, offset int, limit int, field string) *[]*Balance {
	var balance = new([]*Balance)
	if field == "asc" {

		if uid == 0 { //uid为0则查询所有用户
			if ctype != 0 { //查询特定ctype
				Engine.Table("balance").Where("balance.ctype=?", ctype).Limit(limit, offset).Asc("balance.id").Find(balance)
			} else { //不限制ctype
				Engine.Table("balance").Limit(limit, offset).Asc("balance.id").Find(balance)
			}

		} else { //查询特定uid

			if ctype == 0 { //不限制ctype
				Engine.Table("balance").Where("balance.uid=?", uid).Limit(limit, offset).Asc("balance.id").Find(balance)

			} else { //查询特定ctype

				Engine.Table("balance").Where("balance.ctype=? and balance.uid=?", ctype, uid).Limit(limit, offset).Asc("balance.id").Find(balance)
			}
		}
	} else {

		if uid == 0 {
			if ctype != 0 {
				Engine.Table("balance").Where("balance.ctype=?", ctype).Limit(limit, offset).Desc("balance." + field).Find(balance)
			} else {
				Engine.Table("balance").Limit(limit, offset).Desc("balance." + field).Find(balance)
			}

		} else {

			if ctype == 0 {
				Engine.Table("balance").Where("balance.uid=?", uid).Limit(limit, offset).Desc("balance." + field).Find(balance)

			} else {

				Engine.Table("balance").Where("balance.ctype=? and balance.uid=?", ctype, uid).Limit(limit, offset).Desc("balance." + field).Find(balance)
			}
		}
	}
	return balance
}

func GetBalancesByUsername(username string, ctype int64, offset int, limit int, field string) *[]*Balance {
	usr, e := GetUserByUsername(username)
	if e != nil {
		return nil
	} else if usr != nil {
		return GetBalancesByUid(usr.Id, ctype, offset, limit, field)
	}
	return nil
}

func GetBalanceById(id int64) (*Balance, error) {

	balancey := &Balance{}
	has, err := Engine.Id(id).Get(balancey)
	if has {
		return balancey, err
	} else {

		return nil, err
	}
}

func GetBalanceByUid(uid int64) *Balance {
	balance := &Balance{}

	if has, err := Engine.Where("uid=?", uid).Get(balance); err != nil || has == false {
		return nil
	}
	return balance
}

func GetAllBalance() *[]*Balance {
	balances := &[]*Balance{}
	Engine.Desc("id").Find(balances)
	return balances
}

func SetAmountById(id, uid, ctype, amount int64, description string) error {
	balance := &Balance{}
	var err error
	if id > 0 { //大于0则更新指定id的数据
		if balance, err = GetBalanceById(id); err != nil {
			return errors.New(fmt.Sprintf("SetAmountById() GetBalance() Error:", err))
		} else if balance != nil {
			balance.Uid = uid
			balance.Ctype = ctype
			balance.Time = time.Now().Unix()
			balance.Amount = amount
			balance.Balance = balance.Balance + amount
			balance.Description = description
			return SetBalance(balance) //更新数据
		}
	} else { //等于0则插入新id数据
		balance.Id = 0
		balance.Uid = uid
		balance.Ctype = ctype
		balance.Time = time.Now().Unix()
		balance.Amount = amount
		balance.Balance = balance.Balance + amount
		balance.Description = description
		return SetBalance(balance) //插入数据
	}

	return errors.New("SetAmountById() Error")
}

func SetAmountByUid(uid, ctype, amount int64, description string) error {

	if balance := GetBalanceByUid(uid); balance != nil {
		balance.Ctype = ctype
		balance.Time = time.Now().Unix()
		balance.Amount = amount
		balance.Balance = balance.Balance + amount
		balance.Description = description
		return SetBalance(balance) //更新数据
	} else {
		balance := &Balance{}
		balance.Id = 0
		balance.Uid = uid
		balance.Ctype = ctype
		balance.Time = time.Now().Unix()
		balance.Amount = amount
		balance.Balance = balance.Balance + amount
		balance.Description = description
		return SetBalance(balance) //插入数据
	}
	return errors.New("SetAmountByUid() Error")

}

func SetBalance(balance *Balance) error {

	balancey := &Balance{}

	//id大于0则执行更新
	if balance.Id > 0 {

		if _, err := Engine.Id(balance.Id).Get(balancey); err != nil {
			return err
		} else if balancey != nil {

			balancey.Amount = balance.Amount
			balancey.Balance = balancey.Balance + balance.Amount
			balancey.Description = balance.Description
			balancey.Time = balance.Time
			balancey.Ctype = balance.Ctype

			if row, err := Engine.Id(balancey.Id).Cols("time,amount,balance,ctype,description").Update(balancey); err != nil || row == 0 {
				return err
			} else { //更新成功后把数据设到user表的balance字段
				return SetBalanceForUser(balance.Uid, balance.Balance)

			}
		}
	} else { //id小于等于0则执行插入 version=1

		if row, err := Engine.Insert(balance); err != nil || row == 0 {
			return err
		} else { //更新成功后把数据设到user表的balance字段

			return SetBalanceForUser(balance.Uid, balance.Balance)

		}

	}
	return nil

}

func DelBalance(bid int64) error {
	if row, err := Engine.Id(bid).Delete(new(Balance)); err != nil || row == 0 {
		return errors.New("DelBalance() Error!")
	} else {
		return nil
	}
}

func DelBalanceByRole(id int64, uid int64, role int64) error {
	allow := bool(false)
	if anz, err := GetBalanceById(id); err == nil && anz != nil {
		if anz.Uid == uid {
			allow = true
		} else if role < 0 {
			allow = true
		}
		if allow {
			if row, err := Engine.Id(id).Delete(new(Balance)); err != nil || row == 0 {
				return errors.New("row, err := Engine.Id(rid).Delete(new(Balance)) Error!")
			} else {
				return nil
			}
		} else {
			return errors.New("DelBalanceByRole() not allow!")
		}
	} else {
		return errors.New("DelBalanceByRole() GetBalanceById Error")
	}

}

func DelBalancesByUid(uid int64) error {
	balancey := &[]Balance{}
	if err := Engine.Where("uid=?", uid).Find(balancey); err == nil && balancey != nil {
		for _, v := range *balancey {
			if err := DelBalanceByRole(v.Id, v.Uid, -1000); err != nil {
				fmt.Println("DelBalanceByRole:", err)
			}
		}
		return nil
	} else {
		return err
	}
}
