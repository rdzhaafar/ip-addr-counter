package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

var (
	file        = flag.String("file", "", "Input file")
	cpuprofile  = flag.String("cpuprofile", "", "Write CPU profile to file")
	heapprofile = flag.String("heapprofile", "", "Write heap profile to file")
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

	if file == "" {
		die("expected an input file")
	}
	if cpuprofile != "" {
		file, err := os.Create(cpuprofile)
		if err != nil {
			die("%v", err)
		}
		defer file.Close()
		pprof.StartCPUProfile(file)
		defer pprof.StopCPUProfile()
	}

	start := time.Now()
	unique, err := countIPs(file)
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

func countIPs(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	ips := make(map[string]struct{})

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := scanner.Text()
		if ip == "" {
			continue
		}
		ips[ip] = struct{}{}
	}
	return len(ips), nil
}
