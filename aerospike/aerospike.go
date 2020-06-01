// Package aerospike client
package aerospike

import (
	"strconv"
	"strings"
	"time"

	as "github.com/aerospike/aerospike-client-go"
	"github.com/kdpujie/log4go"
)

// Client aerospike client
type Client struct {
	C *as.Client
}

type Option interface {
	apply(do *as.ClientPolicy)
}

type optionFunc func(do *as.ClientPolicy)

func (fn optionFunc) apply(do *as.ClientPolicy) {
	fn(do)
}

func AuthMode(authMode int) Option {
	return optionFunc(func(do *as.ClientPolicy) {
		do.AuthMode = as.AuthMode(authMode)
	})
}

func User(user string) Option {
	return optionFunc(func(do *as.ClientPolicy) {
		do.User = user
	})
}
func Password(password string) Option {
	return optionFunc(func(do *as.ClientPolicy) {
		do.Password = password
	})
}

func Timeout(timeout time.Duration) Option {
	return optionFunc(func(do *as.ClientPolicy) {
		do.Timeout = timeout // 30s
	})
}
func IdleTimeout(idleTimeout time.Duration) Option {
	return optionFunc(func(do *as.ClientPolicy) {
		do.IdleTimeout = idleTimeout // 14s
	})
}
func LoginTimeout(loginTimeout time.Duration) Option {
	return optionFunc(func(do *as.ClientPolicy) {
		do.LoginTimeout = loginTimeout // 10s
	})
}

func ConnectionQueueSize(connectionQueueSize int) Option {
	return optionFunc(func(do *as.ClientPolicy) {
		do.ConnectionQueueSize = connectionQueueSize // 256
	})
}

func MinConnectionsPerNode(minConnectionsPerNode int) Option {
	return optionFunc(func(do *as.ClientPolicy) {
		do.MinConnectionsPerNode = minConnectionsPerNode // 0
	})
}

func LimitConnectionsToQueueSize(limitConnectionsToQueueSize bool) Option {
	return optionFunc(func(do *as.ClientPolicy) {
		do.LimitConnectionsToQueueSize = limitConnectionsToQueueSize // true
	})
}

func OpeningConnectionThreshold(openingConnectionThreshold int) Option {
	return optionFunc(func(do *as.ClientPolicy) {
		do.OpeningConnectionThreshold = openingConnectionThreshold // 0
	})
}

func TendInterval(tendInterval time.Duration) Option {
	return optionFunc(func(do *as.ClientPolicy) {
		do.TendInterval = tendInterval // 1s
	})
}

// New
func New(hosts []string, options ...Option) (*Client, error) {
	var asHosts []*as.Host
	for _, host := range hosts {
		hostSplit := strings.Split(host, ":")
		host := strings.TrimSpace(hostSplit[0])
		port, _ := strconv.ParseInt(hostSplit[1], 10, 64)
		asHosts = append(asHosts, &as.Host{Name: host, Port: int(port)})
	}
	// set the value you set
	policy := &as.ClientPolicy{}
	for _, option := range options {
		option.apply(policy)
	}

	c, err := as.NewClientWithPolicyAndHost(policy, asHosts...)
	if err != nil {
		return nil, err
	}
	isConnected := c.IsConnected()
	if !isConnected {
		log4go.Error("[aerospike] connected failed")
	}
	return &Client{C: c}, err
}

// Close the client
func (c *Client) Close() (err error) {
	c.C.Close()
	return err
}
