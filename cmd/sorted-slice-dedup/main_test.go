package main

import "testing"

func TestCountIPs(t *testing.T) {
	var cases = []struct {
		filename string
		ips      int
	}{
		{"../../test/static/0.txt", 0},
		{"../../test/static/1.txt", 1},
		{"../../test/static/64.txt", 64},
		{"../../test/static/65536.txt", 65536},
	}

	const dedup = 65536

	for _, tc := range cases {
		got, err := sortedSliceCountIPs(tc.filename, dedup)
		if err != nil {
			t.Fatalf("countIPs(%s) returned error: %v", tc.filename, err)
		}
		if got != tc.ips {
			t.Fatalf("countIPs(%s) found %d unique IPs instead of %d", tc.filename, got, tc.ips)
		}
	}
}
