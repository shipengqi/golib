package golib

import "time"

// Retry retries the given callback
func Retry(count int, interval time.Duration, f func(times int) error) error {
	var err error
	total := count
	for ; count > 0; count-- {
		if err = f(total - count); err != nil {
			if interval > 0 {
				<-time.After(interval)
			}
		} else {
			break
		}
	}
	return err
}
