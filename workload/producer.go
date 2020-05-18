package workload

import (
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type Producer struct {
	count      uint64
	sizeAvg    float64
	sizeStddev float64
	addr       string
	tube       string
}

var ttrSecs = 60 * time.Second

func NewProducer(addr, tube string, count uint64, sizeAvg, SizeStdDev float64) *Producer {
	return &Producer{
		count:      count,
		sizeAvg:    sizeAvg,
		sizeStddev: SizeStdDev,
		addr:       addr,
		tube:       tube,
	}
}

func putJob(t *beanstalk.Tube, b []byte) error {
	id, err := t.Put(b, 1, 0, ttrSecs)
	if err != nil {
		return fmt.Errorf("putJob err = %w", err)
	}
	log.Infof("putJob id = %v size = %v bytes", id, len(b))
	return nil
}

func (p *Producer) Run() error {
	c, err := beanstalk.Dial("tcp", p.addr)
	if err != nil {
		return fmt.Errorf("beanstalkDial err  = %w", err)
	}

	t := &beanstalk.Tube{Conn: c, Name: p.tube}
	i := uint64(0)
	for {
		if p.count > 0 && i >= p.count {
			log.Infof("i = %d iterations complete", i)
			return nil
		}
		i++
		log.Infof("iteration %d", i)

		size := rand.NormFloat64()*p.sizeStddev + p.sizeAvg
		if size < 0 {
			size = -size
		}
		d := NewDataBlobFromString(randStr(int(size)))
		b, err := d.Marshal()
		if err != nil {
			return fmt.Errorf("dataMarshal %w", err)
		}

		if err := putJob(t, b); err != nil {
			log.Errorf("putJob err %w", err)
		}
	}
}

func CmdProducer(argv []string, version string) error {
	usage := `usage: producer [options]
Options:
  -h --help                Show this screen.
  --addr=<addr>            Beanstalk addr [default: localhost:11300].
  --tube=<name>            Name of tube to consume [default: workload]
  --count=<n>              Number of items to producer (0 for no limit) [default: 0].
  --avg-data-size=<bytes>  Average data size in bytes of the blob to produce [default: 16384].
  --sd-data-size=<bytes>   Standard deviation in bytes of the data size to produce [default: 4096].
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
	count, err := opts.Int("--count")
	if err != nil {
		return err
	}
	avgDataSize, err := opts.Int("--avg-data-size")
	if err != nil {
		return err
	}
	sdDataSize, err := opts.Int(("--sd-data-size"))
	if err != nil {
		return err
	}

	p := NewProducer(addr, tube, uint64(count), float64(avgDataSize), float64(sdDataSize))
	return p.Run()
}
