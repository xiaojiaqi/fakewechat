package benchmarks

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"runtime"
	"strconv"
	"testing"
	"time"
)

var db *sql.DB
var sum int32

func init() {
	db, _ = sql.Open("mysql", "mygolang:1234@tcp(10.29.101.3:3306)/mydb?charset=utf8")
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(1000)
	db.Ping()
}

/*
func Benchmark_run_pool_100(b *testing.B) {
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < 100; i++ {
		pool()
	}

	b.StopTimer()
}

func Benchmark_grun_pool_10000(b *testing.B) {
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < 10000; i++ {
		pool()
	}

	b.StopTimer()
}


func Benchmark_grun_pool_100000(b *testing.B) {
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < 100000; i++ {
		pool()
	}

	b.StopTimer()
}

*/

func Benchmark_grun_pool_1_10000(b *testing.B) {
	benchmark_run_pool_go(b, 1, 10000)

}

func Benchmark_grun_pool_2_5000(b *testing.B) {
	benchmark_run_pool_go(b, 2, 5000)

}

func Benchmark_grun_pool_10_1000(b *testing.B) {
	benchmark_run_pool_go(b, 10, 1000)

}

func Benchmark_grun_pool_20_500(b *testing.B) {
	benchmark_run_pool_go(b, 20, 500)

}

func Benchmark_grun_pool_30_3340(b *testing.B) {
	benchmark_run_pool_go(b, 30, 334)

}

func Benchmark_grun_pool_50_200(b *testing.B) {
	benchmark_run_pool_go(b, 50, 200)

}

func Benchmark_grun_pool_100_100(b *testing.B) {
	benchmark_run_pool_go(b, 100, 100)

}

func benchmark_run_pool_go(b *testing.B, threads int, loop int) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	defer runtime.GOMAXPROCS(1)
	id = 1
	channel := make(chan int, 10224*16)

	started := time.Now()

	db, _ = sql.Open("mysql", "mygolang:1234@tcp(10.29.101.3:3306)/mydb?charset=utf8")

	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(1000)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < threads; i++ {
		go pool_2(channel, loop)
	}
	for i := 0; i < threads*loop; i++ {
		msg := <-channel
		if msg != 1 {
			panic("Out of sequence")
		}
	}
	b.StopTimer()
	db.Close()

	finished := time.Now()

	fmt.Println(finished.Sub(started))
}

func pool_2(p chan int, loops int) {
	for i := 0; i < loops; i++ {
		pool()
		p <- 1
	}
}

var id int

func pool() {
	sql := "SELECT * FROM rp where id1 = " + strconv.Itoa(id) + " order by id1 limit 1000"
	id += 1
	//fmt.Println(sql)
	rows, err := db.Query(sql)
	defer rows.Close()
	checkErr(err)

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]string)
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
	}

	//	fmt.Println(record)
	// fmt.Fprintln(w, "finish")
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		panic(err)
	}
}
