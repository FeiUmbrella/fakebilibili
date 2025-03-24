package limiter

import (
	"golang.org/x/time/rate"
	"sync"
	"time"
)

//*** 构造限流器 ***//

type Limiters struct {
	limiters *sync.Map // 可并发使用
}

type Limiter struct {
	limiter *rate.Limiter // 限流器
	lastGet time.Time     // 上次获取令牌时间
	key     string        // 邮箱
}

var globalLimiters = &Limiters{
	limiters: &sync.Map{},
}

var once sync.Once // 并发情况可以让某个函数只执行一次

// NewLimiter 返回一个限流器
func NewLimiter(r rate.Limit, b int, key string) *Limiter {
	once.Do(func() { // 启动定时清理map中过期限流器
		go globalLimiters.ClearLimiter()
	})
	keyLimiter := globalLimiters.getLimiter(r, b, key)
	return keyLimiter
}

// getLimiter 返回一个限流器 r:向桶中放token的速率 b:令牌桶的大小 key:可以对服务的id/ip地址进行限制，这里是邮件地址
func (ls *Limiters) getLimiter(r rate.Limit, b int, key string) *Limiter {
	limiter, ok := ls.limiters.Load(key) // 如果该邮件地址已经存在限流器
	if ok {
		return limiter.(*Limiter)
	}

	// 该邮件地址不存在对应限流器,就创建一个并保存
	l := &Limiter{
		limiter: rate.NewLimiter(r, b),
		lastGet: time.Now(),
		key:     key,
	}
	ls.limiters.Store(key, l) // 保存在map中
	return l
}

// ClearLimiter 清除map中超过1分钟的限流器
func (ls *Limiters) ClearLimiter() {
	for {
		time.Sleep(1 * time.Minute)
		ls.limiters.Range(func(key, value interface{}) bool {
			// 超过1分钟
			if time.Now().Sub(value.(*Limiter).lastGet) > time.Minute {
				ls.limiters.Delete(key)
			}
			return true
		})
	}
}

// Allow 取出一个令牌，如果桶内没有令牌则返回false
func (l *Limiter) Allow() bool {
	l.lastGet = time.Now()
	return l.limiter.Allow()
}
