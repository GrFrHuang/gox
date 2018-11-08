package pressure_valve

//The Pressure-Valve use go's limiter and hystrix.
//
// Limiter limits the frequency of occurrence, and USES the algorithm of token pool.
// This pool starts out with b, full of b tokens, and then fills it with r tokens per second.
// Since there is a maximum of b tokens in the token pool, at most b events can be allowed at one time, and one event costs one token.
// If the token runs out, goroutine will block.
//
// Unlike the current limiter's protection mechanism for dependent services, a fuse is used to protect against avalanche effects in order to ensure the normal operation of its own services from accessing the dependent services when the dependent services have failed
// The fuse has three states:
// Off state: service is normal and a failure rate statistics is maintained. When the failure rate reaches the threshold, go to the open state
// Open state: service exception, call fallback function, after a period of time, enter semi-open state
// Half-open: attempt to restore service with a failure rate higher than the threshold, enter the open state, lower than the threshold, and enter the closed state

import (
	"github.com/afex/hystrix-go/hystrix"
	"golang.org/x/time/rate"
	"sync"
	"github.com/GrFrHuang/gox/log"
	"os"
	"context"
	"fmt"
)

const (
	ConstLimit     = 1000
	ConstTimeout   = 1000
	ConstThreshold = 10
)

type PressureValve struct {
	rw          *sync.RWMutex
	limiter     *rate.Limiter
	hystrix     *hystrix.CircuitBreaker
	openSync    bool // Open hystrix sync mode whether or not.
	poolSize    int  // The tokens pool size.
	sleepWindow int
	errorChan   chan error
}

// Purge all hystrix and reset hystrix's error counter.
func (pv *PressureValve) Flush() {
	pv.rw.Lock()
	defer pv.rw.Unlock()
	hystrix.Flush()
}

// Use this function in gateway middleware to limit http request flow.
func (pv *PressureValve) HaveError(err error) {
	testErr := fmt.Errorf("test_data")
	defer func() {
		if r := recover(); r != nil {
			log.Error(r)
		}
	}()
	pv.errorChan <- testErr
	_ = <-pv.errorChan
	// Ensure channel is open . Add err into error channel.
	pv.errorChan <- err
}

// Use this function in gateway middleware to limit http request flow.
func (pv *PressureValve) FlowFilter() (returnErr error) {
	pv.rw.RLock()
	defer pv.rw.RUnlock()
	if !pv.hystrix.IsOpen() {
		ctx := context.Background()
		err := pv.limiter.Wait(ctx)
		if err != nil {
			log.Error(err)
			os.Exit(-1)
		}
	}
	returnErr = hystrix.Do(pv.hystrix.Name, func() error {
		// todo 超时
		return nil
	}, nil)
	return returnErr
}

// Timeout period in millisecond.
// Limit is how many request occur per second.
// threshold is the threshold of the number of requests, which is used to calculate the percentage of the threshold.
// openSync is which whether open hystrix sync mode(Go,Do,).
func NewPressureValve(limit, timeout, threshold int, openSync bool) *PressureValve {
	if limit <= 0 {
		limit = ConstLimit
	}
	poolSize := limit * 10
	if threshold <= 0 || threshold <= limit {
		threshold = ConstThreshold
	}
	if timeout <= 0 {
		timeout = ConstTimeout
	}
	sleepWindow := timeout * 10
	pv := &PressureValve{
		errorChan:   make(chan error, limit),
		rw:          &sync.RWMutex{},
		poolSize:    poolSize,
		openSync:    openSync,
		sleepWindow: sleepWindow,
	}
	// New a Limiter.
	pv.limiter = rate.NewLimiter(rate.Limit(limit), pv.poolSize) // 第一个参数为每秒发生多少次事件（每秒token pool新增多少个token），第二个参数是token pool max size.

	percent := int((float64(threshold) / float64(pv.poolSize)) * 100)
	if percent >= 100 {
		log.Error("limit or threshold's value error !")
		os.Exit(-1)
	}
	name := "GrFrHuangService"
	hystrix.ConfigureCommand(
		name, // 熔断器名字，可以用服务名称命名，一个名字对应一个熔断器，对应一份熔断策略
		hystrix.CommandConfig{
			Timeout:                timeout,        // 超时时间，单位是毫秒
			MaxConcurrentRequests:  pv.poolSize,    // 每秒最大并发数，超过并发返回错误
			RequestVolumeThreshold: threshold,      // 请求数量的阀值，用这些数量的请求来计算阀值的百分比
			ErrorPercentThreshold:  percent,        // 错误数量阀值百分比，达到该错误率，熔断器开启
			SleepWindow:            pv.sleepWindow, // 熔断器间隔多少时间尝试恢复为关闭，单位是毫秒
		},
	)
	// Get a circuit breaker.
	fuse, exist, err := hystrix.GetCircuit(name)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
	if !exist {
		log.Panic(fmt.Printf("熔断器%s不存在", name))
	}
	pv.hystrix = fuse
	// Start listen request error number.
	monitoring(pv)
	return pv
}

func monitoring(pv *PressureValve) {
	go func() {
		for {
			select {
			case err := <-pv.errorChan:
				hystrix.Go(pv.hystrix.Name, func() error {
					return err
				}, func(err error) error {
					log.Error(err)
					return nil
				})
			default:
				continue
			}
		}
	}()
}
