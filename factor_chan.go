package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
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
    process func(a int) map[int]int, exit func()) (chan int, chan map[int]int) {
    inStream := make(chan int, buffer)
    outStream := make(chan map[int]int, buffer)
    ack := make(chan bool)

    for t_worker := 0; t_worker < n_threads; t_worker++ {
        go func() {
            isClosed := false
            for !isClosed {
                num, ok := <-inStream
                if ok {
                    factoredMap := process(num)
                    outStream <- factoredMap
                } else {
                    isClosed = true
                    ack <- true
                }
            }
        }()
    }
    go func() {
        for i := 0; i < n_threads; i++ {
            <-ack
        }
        defer func() { close(outStream) }()
    }()
    return inStream, outStream
}

func main() {
    args := os.Args[1:]
    file, _ := os.Open("nums.txt")
    reader := bufio.NewReader(file)
    nums := make([]int, 0)
    factors := make([]map[int]int, 0)

    n_threads, _ := strconv.Atoi(args[0])

    factorInStream, factorOutStream := worker(n_threads, 100, factor, func() {

    })

    for true {
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        n, _ := strconv.ParseInt(strings.Trim(line, "\n"), 10, 32)
        nums = append(nums, int(n))
    }

    ack := make(chan bool)
    go func() {
        for _, num := range nums {
            factorInStream <- num
        }
        close(factorInStream)
        defer func() { ack <- true }()
    }()

    go func() {
        for true {
            m, ok := <-factorOutStream
            if ok {
                factors = append(factors, m)
            } else {
                break
            }
        }
        // defer func() { writeToFile(nums, factors); ack <- true }()
        defer func() { ack <- true }()
    }()
    <-ack
    <-ack
    close(ack)
}

/*
Generate the file nums.txt
$ python -c "import random; print '\n'.join(str(random.randint(1, 10**7)) for _ in range(1000000))" > nums.txt
$ go run threadripper.go
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
