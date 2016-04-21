package helm

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/deis/tiller/pkg/proto/hapi/services"
)

type client struct {
	cfg  *config
	conn *grpc.ClientConn
	impl services.ReleaseServiceClient
}

func (c *client) dial() (err error) {
	c.conn, err = grpc.Dial(c.cfg.ServAddr, c.cfg.DialOpts()...)
	c.impl = services.NewReleaseServiceClient(c.conn)
	return
}

func (c *client) install(req *services.InstallReleaseRequest) (res *services.InstallReleaseResponse, err error) {
	if err = c.dial(); err != nil {
		return
	}

	defer c.Close()

	return c.impl.InstallRelease(context.TODO(), req, c.cfg.CallOpts()...)
}

func (c *client) Close() error {
	return c.conn.Close()
}
