package redis

import (
	"strconv"
	"context"
	"time"
	"errors"
)

// 保证id的强一致性
// 基于redis实现,所有insert的地方都要上锁
func lock(tableName string, redis *Redis) (id int, err error) {
	for {
		value := ""
		value, err = redis.Get(tableName)
		if err != nil {
			return
		}
		// 有值代表当前已有其他goroutine锁定了这张表
		if value != "" {
			time.Sleep(time.Millisecond * 500)
			continue
		} else {
			// 没有值就获取这个锁,锁超时为2000毫秒
			err = redis.Set(tableName, "locked", 2000, PX)
			if err != nil {
				return
			}
			// 表第一次插入时
			err = redis.SetNX(tableName+"_Ids", 1)
			if err != nil {
				return
			}
			idStr, _ := redis.Get(tableName + "_Ids")
			id, _ = strconv.Atoi(idStr)
			if id == 0 {
				err = errors.New("nil id for table: " + tableName)
			}
			return
		}
	}
	return
}

// 删除锁,id计数器自增加一
func unLock(ctx context.Context, tableName string, redis *Redis) (err error) {
	select {
	case <-ctx.Done():
		_, err = redis.conn.Do("INCR", tableName+"_Ids")
		if err != nil {
			// todo 报警到管理后台
			return
		}
		err = redis.Delete(tableName)
	}
	return
}

func (redis *Redis) CreateId(tableName string) (cancel context.CancelFunc, id int, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err = unLock(ctx, tableName, redis)
		if err != nil {
			return
		}
	}()
	id, err = lock(tableName, redis)
	return
}

func (redis *Redis) End(cancel context.CancelFunc) {
	// 发送取消信号
	cancel()
}
