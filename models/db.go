package models

import (
	"sync"
	"time"
	"github.com/insionng/yougam/libraries/go-xorm/xorm"
)

var (
	cache        = make([]interface{}, 0)
	mutex        sync.RWMutex
	isTimeout    bool
	maxCacheSize = 100
	firstTime    time.Time
	timeOut      = 20 * time.Second
)

func CacheInsert(orm *xorm.Engine, beans ...interface{}) (int64, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if len(cache) == 0 {
		firstTime = time.Now()
		isTimeout = false
	} else {
		isTimeout = time.Now().Sub(firstTime) > timeOut
	}

	cache = append(cache, beans...)
	if len(cache) >= maxCacheSize || isTimeout {
		affected, err := orm.Insert(&cache)
		cache = make([]interface{}, 0)
		return affected, err
	}
	return 0, nil
}

/*
type User struct {
	Id   int64
	Name string
}

func main() {
	orm, err := xorm.NewEngine("sqlite3", "./test.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	orm.ShowSQL = true
	err = orm.Sync(new(User))
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 200; i++ {
		users := []interface{}{
			User{Name: fmt.Sprintf("a%d", i)},
		}
		_, err = CacheInsert(orm, users...)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
*/

/*
import (
	"errors"
	"fmt"
	"github.com/insionng/yougam/helper"
	"time"
)

var (
	TStartTime time.Time
	TList      []interface{}
	TLimit     = 10000
	TTimes     = 2
)
*/
/*
	缓存TTimes 秒 内的所有插入请求,上限设为10000行
	当超过TTimes 秒 或达到上限,并且被插入的内容不为空即执行事务
*/
/*
func Insert(beans ...interface{}) (int64, error) {
	if beans != nil {
		if TStartTime.IsZero() {
			TStartTime = time.Now()
			fmt.Println("开始时间:", TStartTime)
		}

		//小于TTimes秒 or 小于限制
		if since := int(time.Since(TStartTime).Seconds()); (since < TTimes) || (len(TList) < TLimit) {

			fmt.Println("//小于TTimes秒 or 小于限制")
			if len(TList) < TLimit { //如果列表未满 可以无视超时把beans添加到TList

				fmt.Println("//如果列表未满 可以无视超时把beans添加到TList")
				for k, v := range beans {
					TList = append(TList, v)
					fmt.Println("K1#", k, ",append:", v)
				}
				beans = nil
				fmt.Println("TList = append(TList, v):", len(TList))
			}

			if len(TList) >= TLimit { //如果列表已满则执行
				fmt.Println("//如果列表已满则执行")
				return Execute(beans)
			} else if since := int(time.Since(TStartTime).Seconds()); since >= TTimes && (len(TList) > 0) { //另外若果列表未满 但超时 则执行
				fmt.Println("//另外若果列表未满 但超时 则执行")
				return Execute(beans)
			} else { //另外 小于TTimes秒 or 列表未满
				//return Insert(beans) //继续传递未处理beans
				//return -1, errors.New("Insert()发生错误,beans:" + fmt.Sprint(beans))
				fmt.Println("return Execute(beans) //继续传递未处理beans")
				return Execute(beans) //继续传递未处理beans
			}

		} else { //大于TTimes秒 or 大于限制 则执行 并传递未处理beans
			fmt.Println("//大于TTimes秒 or 大于限制 则执行 并传递未处理beans")
			return Execute(beans)
		}

	} else {
		fmt.Println("Insert Errors!")
		return -1, errors.New("Insert Errors!")
	}

}

func Execute(beans ...interface{}) (int64, error) {

	if TList != nil {
		// 创建 Session 对象
		fmt.Println("// 创建 Session 对象")
		sess := Engine.NewSession()
		defer sess.Close()
		// 启动事务
		fmt.Println("// 启动事务")
		if err := sess.Begin(); err != nil {
			fmt.Println("// 启动事务失败")
			return -1, err
		} else {
			fmt.Println("// 启动事务成功")
			if len(TList) < TLimit && (len(beans)-1) != 0 {
				fmt.Println("// len(TList) < TLimit && (len(beans)-1) != 0 ")
				for k, v := range beans {
					TList = append(TList, v)
					fmt.Println(">>>>>>>>>K2#", k, ",append:", v)
				}
			}

			//执行事务
			fmt.Println("时间:", helper.TimeSince(TStartTime), "len(TList):", len(TList))

			for _, v := range TList {
				if row, err := sess.Insert(v); err != nil || row <= 0 {
					fmt.Println("执行事务失败:", err)
					sess.Rollback()
					return -1, err
				}
			}

			// 提交事务
			fmt.Println("//提交事务")
			if err = sess.Commit(); err != nil {
				//事务提交失败后继续传递参数
				//不管参数beans是否为空,下一轮自会检查并处理
				fmt.Println("//事务提交失败后继续传递参数,err:", err)
				return Insert(beans)
				//return -1, err
			} else {
				//事务完成后
				fmt.Println("//事务完成后")
				if len(TList) >= TLimit && beans != nil {
					//重置变量
					fmt.Println("//重置变量AAA")
					TList = nil
					TStartTime = time.Now()
					return Insert(beans) //继续传递未处理 beans
				}
				//重置变量
				fmt.Println("//重置变量BBB")
				TList = nil
				TStartTime = time.Now()
			}
		}
	} else { //TList为nil
		fmt.Println(" //TList为nil start")
		if beans != nil {
			fmt.Println(" //beans != nil")

			TStartTime = time.Now()
			return Insert(beans) //继续传递未处理 beans
		}

		fmt.Println(" ////TList为nil end")
		return -1, errors.New("Execute() 参数不能都为空!")
	}

	return -1, nil
}
*/
