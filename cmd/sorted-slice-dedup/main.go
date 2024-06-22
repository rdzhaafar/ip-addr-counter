package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"runtime/pprof"
	"slices"
	"time"
)

var (
	file        = flag.String("file", "", "Input file")
	cpuprofile  = flag.String("cpuprofile", "", "Write CPU profile to file")
	heapprofile = flag.String("heapprofile", "", "Write heap profile to file")
	dedup       = flag.Int("dedup", 33554432, "Slice deduplication threshold")
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
	dedup := *dedup

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
	unique, err := sortedSliceCountIPs(file, dedup)
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

func sortedSliceCountIPs(filename string, dedup int) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	allIPs := make([]uint32, 0)
	scanner := bufio.NewScanner(file)
	sinceDedup := 0
	for scanner.Scan() {
		s := scanner.Text()
		if s == "" {
			continue
		}
		ip, err := parseIP(s)
		if err != nil {
			return 0, err
		}
		allIPs = append(allIPs, ip)
		sinceDedup++
		if sinceDedup >= dedup {
			allIPs = dedupIPs(allIPs)
			sinceDedup = 0
		}
	}

	slices.Sort(allIPs)
	unique := 0
	last := uint32(0)
	for i := 0; i < len(allIPs); i++ {
		ip := allIPs[i]
		if i == 0 || last != ip {
			unique++
		}
		last = ip
	}

	return unique, nil
}

func parseIP(s string) (uint32, error) {
	netIP := net.ParseIP(s)
	if netIP == nil {
		return 0, fmt.Errorf("invalid IPv4: %s", s)
	}
	return binary.BigEndian.Uint32(netIP[12:16]), nil
}

func dedupIPs(ips []uint32) []uint32 {
	slices.Sort(ips)
	deduped := make([]uint32, 0)
	last := uint32(0)
	for i := 0; i < len(ips); i++ {
		ip := ips[i]
		if i == 0 || ip != last {
			deduped = append(deduped, ip)
		}
		last = ip
	}
	return deduped
}
