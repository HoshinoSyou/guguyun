package util

import (
	"sync"
	"time"
)

// 针对goroutine写的限制goroutine速率的下载限速
// 计划用于断点续传中

type speedLimitMutex struct {
	lock     sync.Mutex
	lastFunc time.Time
	duration time.Duration
	speed    int
}

// SpeedLimit 限制速率
func SpeedLimit(mutex *speedLimitMutex) {
	for {
		mutex.lock.Lock()
		// 获取上一次执行此函数到现在的时间间隔
		dur := time.Now().Sub(mutex.lastFunc)
		// 若时间间隔大于设定的限制间隔（速率），则将上次执行此函数的时间改为现在，并解锁返回
		if dur >= mutex.duration {
			mutex.lastFunc = time.Now()
			mutex.lock.Unlock()
			return
		}
		// 若小于，则休眠一个限制间隔，再将上次执行此函数的使劲啊改为现在，并解锁返回
		time.Sleep(mutex.duration)
		mutex.lastFunc = time.Now()
		mutex.lock.Unlock()
		return
	}
}

// SetLimitSpeed 设置限制速率
func SetLimitSpeed(mutex *speedLimitMutex, speed int) {
	mutex.speed = speed
	// 设置限制间隔 此处设置为一秒进行speed次
	mutex.duration = time.Second * time.Duration(1/speed)
}
