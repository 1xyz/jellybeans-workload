package workload

import (
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"time"
)

type Consumer struct {
	tube    string
	addr    string
	timeout time.Duration
}

func NewConsumer(addr, tube string, timeout time.Duration) *Consumer {
	return &Consumer{
		tube:    tube,
		addr:    addr,
		timeout: timeout,
	}
}

func (c *Consumer) reserveOne(t *beanstalk.TubeSet) error {
	id, bytes, err := t.Reserve(c.timeout)
	if err != nil {
		return fmt.Errorf("reserve error = %w", err)
	}

	log.Infof("Reserved job id = %v, payload-bytes-len=%v", id, len(bytes))
	db, err := NewDataBlobFromBytes(bytes)
	if err != nil {
		return fmt.Errorf("newDataBlobFromBytes %w", err)
	}
	log.Infof("Datablob de-serialized & verified hash slen = %v hashlen: %v",
		len(db.Data), len(db.Hash))

	if err := t.Conn.Delete(id); err != nil {
		return fmt.Errorf("delete err %w", err)
	}
	log.Infof("deleted job w/ id = %v", id)

	return nil
}

func (c *Consumer) Run() error {
	conn, err := beanstalk.Dial("tcp", c.addr)
	if err != nil {
		return fmt.Errorf("error dial beanstald %w", err)
	}

	ts := beanstalk.NewTubeSet(conn, c.tube)
	for {
		if err := c.reserveOne(ts); err != nil {
			log.Errorf("reserveOne: err = %v", err)
		}
	}
}

func CmdConsumer(argv []string, version string) error {
	usage := `usage: consumer [options]
Options:
  -h --help               Show this screen.
  --addr=<addr>           Beanstalk addr [default: localhost:11300].
  --tube=<name>           Name of tube to consume [default: workload]
  --timeout-secs=<secs>   Reservation timeout in seconds [default: 10]..
`
	opts, err := docopt.ParseArgs(usage, argv[1:], version)
	if err != nil {
		log.Fatalf("error parsing arguments. err=%v", err)
	}

	addr, err := opts.String("--addr")
	if err != nil {
		return err
	}
	tube, err := opts.String("--tube")
	if err != nil {
		return err
	}
	timeoutSecs, err := opts.Int("--timeout-secs")
	if err != nil {
		return err
	}

	c := NewConsumer(addr, tube, time.Duration(timeoutSecs)*time.Second)
	return c.Run()
}
