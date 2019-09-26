package client

import (
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/appfirewall/appfirewall/eventInfo"
	"github.com/appfirewall/appfirewall/protocol"
	"github.com/appfirewall/appfirewall/rule"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// GrpcClient connects to the AppFirewall UI that acts as server
type GrpcClient struct {
	sync.Mutex

	unixSocketPath string
	clientConn     *grpc.ClientConn
	afiClient      protocol.AFIClient
}

// NewGrpcClient creates a new GrpcClient with a connection to a unix Socket
func NewGrpcClient(unixSocketPath string) *GrpcClient {
	grpcClient := &GrpcClient{
		unixSocketPath: removeUnixPrefix(unixSocketPath),
	}

	go grpcClient.connecter()
	return grpcClient
}

func removeUnixPrefix(unixSocketPath string) string {
	if strings.HasPrefix(unixSocketPath, "unix://") == true {
		return unixSocketPath[7:]
	}
	return unixSocketPath
}

func (c *GrpcClient) connecter() {
	log.Printf("connecter started for socket %s", c.unixSocketPath)
	wasConnected := false
	for true {
		isConnected := c.isConnected()
		if wasConnected != isConnected {
			c.logConnectionChange(isConnected)
			wasConnected = isConnected
		}

		if err := c.connect(); err != nil {
			log.Printf("error while connecting to: %s", err)
		}

		time.Sleep(1 * time.Second)
	}
}

func (c *GrpcClient) isConnected() bool {
	c.Lock()
	defer c.Unlock()
	if c.clientConn == nil || c.clientConn.GetState() != connectivity.Ready {
		return false
	}
	return true
}

func (c *GrpcClient) logConnectionChange(connected bool) {
	if connected {
		log.Printf("connected to %s", c.unixSocketPath)
	} else {
		log.Fatal("connection lost.")
	}
}

func (c *GrpcClient) connect() (err error) {
	if c.isConnected() {
		return
	}

	c.clientConn, err = grpc.Dial(c.unixSocketPath, grpc.WithInsecure(),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}))

	if err != nil {
		c.clientConn = nil
		return err
	}

	c.afiClient = protocol.NewAFIClient(c.clientConn)
	return nil
}

// Prompt shows the current connection to the user and asks for an action
func (c *GrpcClient) Prompt(event *eventInfo.EventPayload) (*rule.Rule, bool) {

	// TODO block connection by default

	c.Lock()
	defer c.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1800) // Setup Timeout
	defer cancel()
	reply, err := c.afiClient.Prompt(ctx, event.ToAFConnectionInfo())
	if err != nil {
		log.Fatal("error while asking for rule:", err)
		// TODO block connection on error
	}

	return rule.FromAFRule(reply), true
}
