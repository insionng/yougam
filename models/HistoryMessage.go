package models

type HistoryMessage struct {
	Id       int64
	Key      string `xorm:"index"`
	Uid      int64  `xorm:"index"` //发送者之用户ID
	Sender   string `xorm:"index"` //发送者之用户名称
	Avatar   string `xorm:"index"` //发送者之头像
	Receiver string `xorm:"index"` //接收者之用户名称
	Content  string `xorm:"index"` //搜索功能需要索引
	Created  int64  `xorm:"index"`
}

func DelHistoryMessage(id int64) error {
	if row, err := Engine.Id(id).Delete(new(HistoryMessage)); (err != nil) || (row <= 0) {
		return err
	} else {
		return nil
	}
}

func GetHistoryMessages(offset, limit int, field string) (*[]*HistoryMessage, error) {
	msgs := new([]*HistoryMessage)
	err := Engine.Limit(limit, offset).Desc(field).Find(msgs)
	return msgs, err
}

func GetHistoryMessagesViaReceiver(offset, limit int, username, field string) (*[]*HistoryMessage, error) {
	msgs := new([]*HistoryMessage)
	err := Engine.Where("receiver=?", username).Limit(limit, offset).Desc(field).Find(msgs)
	return msgs, err
}

func GetHistoryMessagesViaReceiverWithSender(offset, limit int, receiver, sender, field string) (*[]*HistoryMessage, error) {
	msgs := new([]*HistoryMessage)
	var err error
	if field == "asc" {
		err = Engine.Where("receiver=? and sender=?", receiver, sender).Limit(limit, offset).Find(msgs)
	} else {
		err = Engine.Where("receiver=? and sender=?", receiver, sender).Limit(limit, offset).Desc(field).Find(msgs)
	}
	return msgs, err
}

func GetHistoryMessagesViaKey(offset, limit int, key, field string) (*[]*HistoryMessage, error) {
	msgs := new([]*HistoryMessage)
	err := Engine.Where("key=?", key).Limit(limit, offset).Desc(field).Find(msgs)
	return msgs, err
}

func GetHistoryMessage(id int64) (*HistoryMessage, error) {

	m := &HistoryMessage{}

	has, err := Engine.Id(id).Get(m)
	if has {
		return m, err
	} else {
		return nil, err
	}
}

func GetHistoryMessageViaKey(key string) (*HistoryMessage, error) {

	m := &HistoryMessage{Key: key}

	has, err := Engine.Get(m)
	if has {
		return m, err
	} else {

		return nil, err
	}
}

func PostHistoryMessage(m *HistoryMessage) (int64, error) {
	// 创建 Session 对象
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	HistoryMessageid := int64(0)
	if _, err := sess.Insert(m); err == nil && m != nil {
		HistoryMessageid = m.Id
	} else {
		// 发生错误时进行回滚
		sess.Rollback()
		return -1, err
	}

	// 提交事务
	return HistoryMessageid, sess.Commit()
}

func PutHistoryMessage(mid int64, msg *HistoryMessage) (int64, error) {
	//覆盖式更新
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	HistoryMessageid := int64(0)
	if row, err := sess.Update(msg, &HistoryMessage{Id: mid}); err != nil || row <= 0 {
		sess.Rollback()
		return -1, err
	} else {
		HistoryMessageid = msg.Id
	}

	// 提交事务
	return HistoryMessageid, sess.Commit()
}

func PutHistoryMessageViaKey(key string, msg *HistoryMessage) (int64, error) {
	//覆盖式更新
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	HistoryMessageid := int64(0)
	if row, err := sess.Update(msg, &HistoryMessage{Key: key}); err != nil || row <= 0 {
		sess.Rollback()
		return -1, err
	} else {
		HistoryMessageid = msg.Id
	}

	// 提交事务
	return HistoryMessageid, sess.Commit()
}
