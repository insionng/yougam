package models

import (
	"fmt"
)

type IssueMark struct {
	Id      int64
	Uid     int64  `xorm:"index"` //举报人
	Rid     int64  `xorm:"index"` //被举报项id
	Ctype   int64  `xorm:"index"` // 1为topic类型  -1是comment类型
	Title   string `xorm:"index"` //举报标题
	Content string `xorm:"text"`  //举报内容
}

//SetIssueMark 设置举报标记
func SetIssueMark(uid int64, rid int64, ctype int64, content string) (int64, error) {
	rptm := &IssueMark{Uid: uid, Rid: rid, Ctype: ctype, Content: content}
	rows, err := Engine.Insert(rptm)
	return rows, err
}

func IsIssueMark(uid int64, rid int64, ctype int64) bool {

	rptm := &IssueMark{}

	if has, err := Engine.Where("uid=? and rid=? and ctype=?", uid, rid, ctype).Get(rptm); err != nil {
		fmt.Println(err)
		return false
	} else {
		if has {
			if rptm.Uid == uid {
				return true
			} else {
				return false
			}

		} else {
			return false
		}
	}

}

func GetIssue(rid int64) (*IssueMark, error) {

	rm := &IssueMark{}

	has, err := Engine.Id(rid).Get(rm)
	if has {
		return rm, err
	} else {

		return nil, err
	}
}

func GetIssues(offset int, limit int, field string, ctype int64) (*[]*IssueMark, error) {

	rms := new([]*IssueMark)

	if ctype == 0 { //同时查询两种类型举报项
		err := Engine.Limit(limit, offset).Desc(field).Find(rms)
		return rms, err

	} else if (ctype == 1) || (ctype == -1) { //1为topic类型 -1是comment类型 按指定类型查询举报项目

		err := Engine.Where("ctype=?", ctype).Limit(limit, offset).Desc(field).Find(rms)
		return rms, err
	} else {

		return nil, nil
	}
}
