package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/rdzhaafar/ip-addr-counter/pkg/ipstore/v1"
)

var (
	file        = flag.String("file", "", "Input file")
	cpuprofile  = flag.String("cpuprofile", "", "Write CPU profile to file")
	heapprofile = flag.String("heapprofile", "", "Write heap profile to file")
	chunks      = flag.Int("chunks", 32, "Number of file chunks to process in parallel")
	buffer      = flag.Int("buffer", 4096, "Buffer size for task channels")
)

func die(m string, a ...interface{}) {
	message := "error: " + fmt.Sprintf(m, a...)
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

func main() {
	flag.Parse()
	file := *file
	cpuprofile := *cpuprofile
	heapprofile := *heapprofile
	chunks := *chunks
	buffer := *buffer

	if file == "" {
		die("expected an input file")
	}
	if cpuprofile != "" {
		file, err := os.Create(cpuprofile)
		if err != nil {
			die("%v", err)
		}
		defer file.Close()
		if err := pprof.StartCPUProfile(file); err != nil {
			die("%v", err)
		}
		defer pprof.StopCPUProfile()
	}

	start := time.Now()
	unique, err := countIPs(file, chunks, buffer)
	if err != nil {
		die("%v", err)
	}
	fmt.Printf("found %d unique IP addresses in %s\n", unique, file)
	fmt.Printf("took %v\n", time.Since(start))

	if heapprofile != "" {
		file, err := os.Create(heapprofile)
		if err != nil {
			die("%v", err)
		}
		defer file.Close()
		if err := pprof.WriteHeapProfile(file); err != nil {
			die("%v", err)
		}
	}
}

func countIPs(filename string, chunks int, buffer int) (int, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	chunkSize := stat.Size() / int64(chunks)

	ips := make(chan uint32, buffer)
	errors := make(chan error)
	count := make(chan int)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go collectIPs(ips, count, ctx)
	offset := int64(0)
	for i := 0; i < chunks; i++ {
		go readFileChunk(filename, offset, chunkSize, ips, errors, ctx, &wg)
		wg.Add(1)
		offset += chunkSize
	}
	wg.Wait()
	close(ips)

	select {
	case err := <-errors:
		return 0, err
	case c := <-count:
		return c, nil
	}
}

func readFileChunk(filename string, offset int64, chunk int64, ips chan<- uint32, errors chan<- error, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(filename)
	if err != nil {
		errors <- err
		return
	}
	defer file.Close()

	// If we're not at the beginning of the file, we need to scroll back to
	// find the end of the last line
	offset, chunk, err = findLastLineOffset(file, offset, chunk)
	if err != nil {
		errors <- err
		return
	}

	// Read assigned chunk and populate the ips channel
	_, err = file.Seek(offset, 0)
	if err != nil {
		errors <- err
		return
	}
	scanner := bufio.NewScanner(file)
	read := int64(0)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			s := scanner.Text()
			if s == "" {
				read += 1
			} else {
				ip, err := parseIP(s)
				if err != nil {
					errors <- err
					return
				}
				ips <- ip
				read += int64(len(s)) + 1
			}
			if read >= chunk {
				return
			}
		}
	}
}

func findLastLineOffset(file *os.File, offset int64, chunk int64) (int64, int64, error) {
	if offset == 0 {
		return offset, chunk, nil
	}

	const bufsize = 16
	buf := make([]byte, bufsize)
	i := bufsize - 1
	_, err := file.ReadAt(buf, offset-bufsize)
	atIPStart := false
	if err != nil {
		return 0, 0, err
	}
	for offset != 0 && !atIPStart {
		b := buf[i]
		if b == '\n' {
			atIPStart = true
			break
		}
		i--
		offset--
	}
	chunk += int64(bufsize + i - 1)
	return offset, chunk, nil
}

func collectIPs(ips <-chan uint32, count chan<- int, ctx context.Context) {
	store := ipstore.NewIPStore()

	for {
		select {
		case <-ctx.Done():
			return
		case ip, ok := <-ips:
			if !ok {
				count <- store.Count()
				return
			}
			store.Insert(ip)
		}
	}
}

func parseIP(s string) (uint32, error) {
	netIP := net.ParseIP(s)
	if netIP == nil {
		return 0, fmt.Errorf("invalid IPv4: %s", s)
	}
	return binary.BigEndian.Uint32(netIP[12:16]), nil
}
