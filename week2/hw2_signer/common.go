package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)


type job func(in, out chan interface{})

const (
	MaxInputDataLen = 100
)

var (
	dataSignerOverheat uint32 = 0
	DataSignerSalt            = ""
)

var OverheatLock = func() {
	for {
		if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 0, 1); !swapped {
			fmt.Println("OverheatLock happend")
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

var OverheatUnlock = func() {
	for {
		if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 1, 0); !swapped {
			fmt.Println("OverheatUnlock happend")
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

var DataSignerMd5 = func(data string) string {
	OverheatLock()
	defer OverheatUnlock()
	data += DataSignerSalt
	dataHash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	time.Sleep(10 * time.Millisecond)
	fmt.Println("DataSignerMd5 has done, result:", dataHash)
	return dataHash
}

var DataSignerCrc32 = func(data string) string {
	data += DataSignerSalt
	crcH := crc32.ChecksumIEEE([]byte(data))
	dataHash := strconv.FormatUint(uint64(crcH), 10)
	time.Sleep(time.Second)

	fmt.Println("DataSignerCrc32 has done, result:", dataHash)

	return dataHash
}



// все что ниже есть всратая моя работа
// и я хз как вообще его начать
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func ExecutePipeline(works ...job) {
	for _, iWork := range works {
		in := make(chan interface{}, 1)
		out := make(chan interface{}, 1)
		go iWork(in, out)
	}
}

func main() {
	runtime.GOMAXPROCS(0)
	//inChan := make(chan int)
	//outChan := make(chan interface{})

	//SingleHash()
	MultiHash()
	//time.Sleep(199 * time.Second)
}

// crc32(md5(data)) + "~" + cdc32(data)
func SingleHash() {
	defer timeTrack(time.Now(), "SingleHash:")
	//fmt.Println(";ll")
	strFromIn := fmt.Sprintf("%v", 1)

	Md5Crc32Chan := make(chan string)
	Crc32Chan := make(chan string)

	go func(in chan string) {
		in <- DataSignerMd5(<-in)
	}(Md5Crc32Chan)

	go func(in chan string) {
		in <- DataSignerCrc32(<-in)
	}(Crc32Chan)

	Crc32Chan <- strFromIn
	Md5Crc32Chan <- strFromIn

	// DataSignerMd5 take 10 msec to execute
	// DataSignerCrc32 take 1 sec to execute
	// There was no reason to go one more goroutine and it didn't have to do that, but why not?
	result := <-Md5Crc32Chan
	fmt.Println(result)
	go func(in chan string) {
		in <- DataSignerCrc32(<-in)
	}(Md5Crc32Chan)
	Md5Crc32Chan <- result

	result = <-Crc32Chan
	fmt.Println(result)
	result += "~"
	result += <-Md5Crc32Chan

	fmt.Println("SingleHash result:", result)

}

// * MultiHash считает значение crc32(th+data)) (конкатенация цифры, приведённой к строке и строки),
//где th=0..5 ( т.е. 6 хешей на каждое входящее значение ), потом берёт конкатенацию результатов в
//порядке расчета (0..5), где data - то что пришло на вход (и ушло на выход из SingleHash)

func MultiHash() {
	data := "4108050209~502633748"

	defer timeTrack(time.Now(), "MultiHash")
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	type threadVal struct {
		th  int
		val string
	}

	crc32Chan := make(chan threadVal)
	var intermResult = make(map[int]string, 6)

	for th := 0; th < 6; th++ {
		wg.Add(1)
		go func(inIn chan threadVal, hashMap map[int]string, mu * sync.Mutex, waiter *sync.WaitGroup) {
			defer waiter.Done()
			thrV := <- inIn
			thrV.val = DataSignerCrc32(thrV.val)
			mu.Lock()
			hashMap[thrV.th] = thrV.val
			mu.Unlock()
		}(crc32Chan, intermResult, mu, wg)

		crc32Chan <- threadVal{ th, strconv.Itoa(th) + data}
	}
	wg.Wait()

	result := ""

	for i := 0; i < 6; i++ {
		mu.Lock()
		result += intermResult[i]
		mu.Unlock()
	}
	fmt.Println(result)

}
