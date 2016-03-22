package benchmarks

import (
	"fmt"
	"log"
	"runtime"
	"testing"
	"time"
)

import proto "github.com/golang/protobuf/proto"
import (
	"github.com/fakewechat/message"
)

var sum int32

func init() {
	fmt.Println("Init() called")

}

func Benchmark_pack_10_1000000(b *testing.B) {
	benchmark_run_pool_go(b, 10, 1000000)
}

func Benchmark_pack_100_100000(b *testing.B) {
	benchmark_run_pool_go(b, 100, 100000)
}

func Benchmark_parse_1_10000000(b *testing.B) {
	benchmark_run_pool_go(b, 1, 200)
}

func benchmark_run_pool_go(b *testing.B, threads int, loop int) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	defer runtime.GOMAXPROCS(1)

	channel := make(chan int, 10224*16)

	datachannel := make(chan *[]byte, 2*1000*10*1000)
	started := time.Now()

	for i := 0; i < threads; i++ {
		go pool_2(channel, loop, datachannel)
	}
	for i := 0; i < threads*loop; i++ {
		msg := <-channel
		if msg != 1 {
			panic("Out of sequence")
		}
	}
	b.StopTimer()

	finished := time.Now()
	fmt.Println(finished.Sub(started))
}

func pool_2(p chan int, loops int, datachannel chan *[]byte) {
	for i := 0; i < loops; i++ {
		pool(datachannel)
		p <- 1
	}
}

func pool(datachannel chan *[]byte) {
	msg := &message.MonitorStatus{
		Timestamp:     proto.Int64(time.Now().UnixNano()),
		Name:          proto.String("hello"),
		Servername:    proto.String("td01"),
		Servertype:    proto.String("tdserver"),
		Absolutevalue: proto.Int64(10000),
		Changesvalue:  proto.Int64(-200),
	} //msg init
	buf, err := proto.Marshal(msg) //SerializeToOstream
	if err != nil {
		fmt.Println("发送数据失败!", err)
	}
	datachannel <- &buf
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		panic(err)
	}
}
