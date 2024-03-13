package publish

import (
	"context"
	"fmt"
	"math/rand"
	"sincerecloud.com/openstack_notifier/global"
	"sincerecloud.com/openstack_notifier/utils"
	"time"
)

var (
	DefaultRetries = 4
	DefaultSleep = 500 * time.Millisecond
)

type Func func(ctx context.Context, body []byte) error

func RetryDo(fn Func, ctx context.Context, body []byte, retries int, sleep time.Duration) error {
	if sleep == 0 {
		sleep = DefaultSleep
	}

	if err := fn(ctx, body); err != nil {
		retries--
		if retries <= 0 {
			return err
		}
        global.LOG.Info(fmt.Sprintf("[req-%s]Retry publish %d times", utils.GetRequestID(ctx), DefaultRetries - retries + 1))
		sleep += (time.Duration(rand.Int63n(int64(sleep)))) / 2
		time.Sleep(sleep)
		return RetryDo(fn, ctx, body, retries, 2*sleep)
	}
	return nil
}


