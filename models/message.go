package models

type Message struct {
	Id       int64
	Key      string `xorm:"index"`
	Uid      int64  `xorm:"index"` //发送者之用户ID
	Sender   string `xorm:"index"` //发送者之用户名称
	Avatar   string `xorm:"index"` //发送者之头像
	Receiver string `xorm:"index"` //接收者之用户名称
	Content  string `xorm:"index"` //搜索功能需要索引
	Created  int64  `xorm:"index"`
}

func DelMessage(id int64) error {
	if row, err := Engine.Id(id).Delete(new(Message)); (err != nil) || (row <= 0) {
		return err
	} else {
		return nil
	}
}

func GetMessages(offset, limit int, field string) (*[]*Message, error) {
	msgs := new([]*Message)
	err := Engine.Limit(limit, offset).Desc(field).Find(msgs)
	return msgs, err
}

func GetMessagesViaReceiver(offset, limit int, username, field string) (*[]*Message, error) {
	msgs := new([]*Message)
	err := Engine.Where("receiver=?", username).Limit(limit, offset).Desc(field).Find(msgs)
	return msgs, err
}

func GetMessagesViaReceiverWithSender(offset, limit int, receiver, sender, field string) (*[]*Message, error) {
	msgs := new([]*Message)
	var err error
	if field == "asc" {
		err = Engine.Where("receiver=? and sender=?", receiver, sender).Limit(limit, offset).Find(msgs)
	} else {
		err = Engine.Where("receiver=? and sender=?", receiver, sender).Limit(limit, offset).Desc(field).Find(msgs)
	}
	return msgs, err
}

func GetMessagesViaKey(offset, limit int, key, field string) (*[]*Message, error) {
	msgs := new([]*Message)
	err := Engine.Where("key=?", key).Limit(limit, offset).Desc(field).Find(msgs)
	return msgs, err
}

func GetMessage(id int64) (*Message, error) {

	m := &Message{}

	has, err := Engine.Id(id).Get(m)
	if has {
		return m, err
	} else {
		return nil, err
	}
}

func GetMessageViaKey(key string) (*Message, error) {

	m := &Message{Key: key}

	has, err := Engine.Get(m)
	if has {
		return m, err
	} else {

		return nil, err
	}
}

func PostMessage(m *Message) (int64, error) {
	// 创建 Session 对象
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	messageid := int64(0)
	if _, err := sess.Insert(m); err == nil && m != nil {
		messageid = m.Id
	} else {
		// 发生错误时进行回滚
		sess.Rollback()
		return -1, err
	}

	// 提交事务
	return messageid, sess.Commit()
}

func PutMessage(mid int64, msg *Message) (int64, error) {
	//覆盖式更新
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	messageid := int64(0)
	if row, err := sess.Update(msg, &Message{Id: mid}); err != nil || row <= 0 {
		sess.Rollback()
		return -1, err
	} else {
		messageid = msg.Id
	}

	// 提交事务
	return messageid, sess.Commit()
}

func PutMessageViaKey(key string, msg *Message) (int64, error) {
	//覆盖式更新
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	messageid := int64(0)
	if row, err := sess.Update(msg, &Message{Key: key}); err != nil || row <= 0 {
		sess.Rollback()
		return -1, err
	} else {
		messageid = msg.Id
	}

	// 提交事务
	return messageid, sess.Commit()
}
