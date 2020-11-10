package main

import (
	"runtime"
	"strconv"
	"sync"
	"time"
	"fmt"
)

// все что ниже есть всратая моя работа
// и я хз как вообще его начать
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func ExecutePipeline(works ...job) {
	jobNum := len(works)
	chanSlice := make([]chan string, jobNum + 1)

	for i := 0; i < jobNum + 1; i++ {
		chanSlice[i] = make(chan string)
	}


}

func main() {
	runtime.GOMAXPROCS(0)
	inChan := make(chan interface{})
	outChan := make(chan interface{})

	//go SingleHash(inChan, outChan)
	//inChan <- 0
	//fmt.Println("4108050209~502633748" == <- outChan)
	go MultiHash(inChan, outChan)
	inChan <- "4108050209~502633748"
	fmt.Println("29568666068035183841425683795340791879727309630931025356555" == <-outChan)
	//MultiHash()
	//time.Sleep(199 * time.Second)
}

// crc32(md5(data)) + "~" + cdc32(data)
func SingleHash(in, out chan interface{}) {
	defer timeTrack(time.Now(), "SingleHash:")

	strFromIn := fmt.Sprintf("%v", <-in)

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

	go func(in chan string) {
		in <- DataSignerCrc32(<-in)
	}(Md5Crc32Chan)

	result := <-Crc32Chan
	out <- result + "~" + <-Md5Crc32Chan
}

// * MultiHash считает значение crc32(th+data)) (конкатенация цифры, приведённой к строке и строки),
//где th=0..5 ( т.е. 6 хешей на каждое входящее значение ), потом берёт конкатенацию результатов в
//порядке расчета (0..5), где data - то что пришло на вход (и ушло на выход из SingleHash)

func MultiHash(in, out chan interface{}) {
	data := fmt.Sprintf("%v", <-in)

	defer timeTrack(time.Now(), "MultiHash")

	wg := &sync.WaitGroup{}

	type threadVal struct {
		th  int
		val string
	}

	crc32ChanIn := make(chan threadVal, 5)
	crc32ChanOut := make(chan threadVal, 5)
	var intermResult = make(map[int]string, 6)

	go func(in, out chan threadVal){
		wg.Add(6)

		for th := 0; th < 6; th++ {
			go func(in, out chan threadVal, waiter *sync.WaitGroup) {
				defer waiter.Done()

				thrV := <-in
				thrV.val = DataSignerCrc32(thrV.val)
				out <- thrV

			}(crc32ChanIn, crc32ChanOut, wg)

			crc32ChanIn <- threadVal{th, strconv.Itoa(th) + data}
		}
		wg.Wait()
		close(crc32ChanOut)

	}(crc32ChanIn, crc32ChanOut)

	result := ""

	for outThreadValue := range crc32ChanOut {
		intermResult[outThreadValue.th] = outThreadValue.val
	}
	for i := 0; i < 6; i++ {
		result += intermResult[i]
	}

	out <- result
}

func MultiHash2(in, out chan interface{}) {
	data := fmt.Sprintf("%v", <-in)

	defer timeTrack(time.Now(), "MultiHash")
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	type threadVal struct {
		th  int
		val string
	}

	crc32Chan := make(chan threadVal, 5)
	var intermResult = make(map[int]string, 6)

	wg.Add(6)
	for th := 0; th < 6; th++ {
		go func(in chan threadVal, hashMap map[int]string, mu *sync.Mutex, waiter *sync.WaitGroup) {
			defer waiter.Done()

			thrV := <-in
			thrV.val = DataSignerCrc32(thrV.val)
			mu.Lock()
			hashMap[thrV.th] = thrV.val
			mu.Unlock()
		}(crc32Chan, intermResult, mu, wg)

		crc32Chan <- threadVal{th, strconv.Itoa(th) + data}
	}
	wg.Wait()

	result := ""

	for i := 0; i < 6; i++ {
		result += intermResult[i]
	}
	//fmt.Println(result)
	out <- result
}

func CombineResults(in, out chan interface{}){

}
