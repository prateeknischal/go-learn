package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	N     = int(1e7 + 5)
	prime = make([]int, N)
)

func init() {
	prime[1] = 1
	prime[2] = 0
	for i := 3; i*i < N; i++ {
		if prime[i] == 0 {
			for j := i * i; j < N; j += (i * 2) {
				prime[j] = i
			}
		}
	}
	for i := 1; i < N; i++ {
		if i%2 == 0 {
			prime[i] = 2
		} else if prime[i] == 0 {
			prime[i] = i
		}
	}
}

func factor(n int) map[int]int {
	mp := make(map[int]int)
	for n > 1 {
		mp[prime[n]]++
		n = n / prime[n]
	}
	return mp
}

func writeToFile(nums []int, v []map[int]int) {
	writer, _ := os.Create("factored.txt")
	for i := 0; i < len(nums); i++ {
		fmt.Fprintf(writer, "%d -> ", nums[i])
		for k := range v[i] {
			fmt.Fprintf(writer, "(%d:%d),", k, v[i][k])
		}
		fmt.Fprintln(writer, "")
	}
}

func worker(n_threads int, buffer int,
	process func(a int) map[int]int, wg *sync.WaitGroup) (chan int, chan map[int]int) {
	inStream := make(chan int, buffer)
	outStream := make(chan map[int]int, buffer)

	for t_worker := 0; t_worker < n_threads; t_worker++ {
		wg.Add(1)
		go func() {
			isClosed := false
			for !isClosed {
				num, ok := <-inStream
				if ok {
					factoredMap := process(num)
					outStream <- factoredMap
				} else {
					isClosed = true
				}
			}
			defer wg.Done()
		}()
	}
	return inStream, outStream
}

func main() {
	args := os.Args[1:]
	file, _ := os.Open("nums.txt")
	reader := bufio.NewReader(file)
	nums := make([]int, 0)
	factors := make([]map[int]int, 0)

	n_threads, _ := strconv.Atoi(args[0])
	var wg sync.WaitGroup
	var m_wg sync.WaitGroup

	factorInStream, factorOutStream := worker(n_threads, 100, factor, &wg)

	for true {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		n, _ := strconv.ParseInt(strings.Trim(line, "\n"), 10, 32)
		nums = append(nums, int(n))
	}

	go func() {
		m_wg.Add(1)
		for _, num := range nums {
			factorInStream <- num
		}
		close(factorInStream)
		defer m_wg.Done()
	}()

	go func() {
		m_wg.Add(1)
		for true {
			m, ok := <-factorOutStream
			if ok {
				factors = append(factors, m)
			} else {
				break
			}
		}
		// defer func() { writeToFile(nums, factors); m_wg.Done() }()
		defer m_wg.Done()
	}()
	wg.Wait()
	close(factorOutStream)
	m_wg.Wait()
}

/*
Generate the file nums.txt
$ python -c "import random; print '\n'.join(str(random.randint(1, 10**7)) for _ in range(1000000))" > nums.txt
$ time go run factor_wg.go 8
performance:
Go version : go version go1.8.3 darwin/amd64
Specs : Intel(R) Core(TM) i7-4770HQ CPU @ 2.20GHz, 2x8G mem
Sieve = 10000000
No of integers to be factored  = 100000
n_threads time(s)
1         1.60
2         1.37
4         1.24
8         1.06
10        1.1
*/
