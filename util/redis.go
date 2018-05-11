package util

import (
	"errors"
	"time"

	"github.com/FZambia/go-sentinel"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
)

func NewSentinelPool() *redis.Pool {
	sntnl := &sentinel.Sentinel{
		Addrs:      []string{"redis-sentinel:26379"},
		MasterName: "mymaster",
		Dial: func(addr string) (redis.Conn, error) {
			timeout := 500 * time.Millisecond
			c, err := redis.DialTimeout("tcp", addr, timeout, timeout, timeout)
			if err != nil {
				glog.Infof("newSentinelPool redis.DialTimeout error addr:%v timeout:%v err:%v", addr, timeout, err)
				return nil, err
			}
			glog.Infof("newSentinelPool redis.DialTimeout SUCCESS addr:%v timeout:%v", addr, timeout)
			return c, nil
		},
	}
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   64,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			masterAddr, err := sntnl.MasterAddr()
			if err != nil {
				glog.Infof("newSentinelPool sntnl.MasterAddr sntnl:%#v err:%v", sntnl, err)
				return nil, err
			}
			glog.Infof("newSentinelPool sntnl.MasterAddr SUCCESS masterAddr:%#v, sntnl:%#v", masterAddr, sntnl)
			c, err := redis.Dial("tcp", masterAddr)
			if err != nil {
				glog.Infof("newSentinelPool masterAddr:%v redis.Dial sntnl:%#v err:%v", masterAddr, sntnl, err)
				return nil, err
			}
			glog.Infof("newSentinelPool redis.Dial SUCCESS masterAddr:%#v, sntnl:%#v", masterAddr, sntnl)
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if !sentinel.TestRole(c, "master") {
				glog.Infof("newSentinelPool sentinel.TestRole failed sntnl:%#v redis.Conn:%#v", sntnl, c)
				return errors.New("Role check failed")
			} else {
				return nil
			}
		},
	}
}
