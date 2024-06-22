package ipstore

import (
	"math"
	"testing"
)

func TestIPStoreHash(t *testing.T) {
	var cases = []struct {
		ip   uint32
		hash int
	}{
		{0, 0},
		{0xff_00_00_00, 0xff_00_00},
		{0x00_ff_ff_00, 0x00_ff_ff},
		{math.MaxUint32, nBuckets},
	}

	for _, c := range cases {
		h := ipStoreHash(c.ip)
		if h != c.hash {
			t.Errorf("ipStoreHash() failed. IP: %d, want: %d, got %d", c.ip, c.hash, h)
		}
	}
}

func TestIPStoreInsert(t *testing.T) {
	store := NewIPStore()

	store.Insert(uint32(0))
	if len(store.buckets[0]) == 0 || store.buckets[0][0] != 0 {
		t.Errorf("IPStore.Insert(0) did not work correctly")
	}

	store.Insert(math.MaxUint32)
	if len(store.buckets[nBuckets]) == 0 || store.buckets[nBuckets][0] != math.MaxUint32 {
		t.Errorf("IPStore.Insert(math.MaxUint32) did not work correctly")
	}
}

func TestIPStoreBucketDistribution(t *testing.T) {
	store := NewIPStore()
	for i := uint32(0); i < 0xff; i++ {
		ip := i << 8
		store.Insert(ip)
	}
	for i := 0; i < 0xff; i++ {
		bucket := store.buckets[i]
		if len(bucket) != 1 {
			t.Fatalf("len(IPStore.buckets[%d]) == %d", i, len(bucket))
		}
	}
}

func TestIPStoreCount(t *testing.T) {
	const elements = 1_000_000

	store1 := NewIPStore()
	for i := 0; i < elements; i++ {
		store1.Insert(uint32(i))
	}
	if store1.Count() != elements {
		t.Errorf("IPStore.Count() == %d after inserting %d unique elements", store1.Count(), elements)
	}

	store2 := NewIPStore()
	for i := 0; i < elements; i++ {
		store2.Insert(uint32(0))
	}
	if store2.Count() != 1 {
		t.Errorf("IPStore.Count() == %d after inserting 1 unique element", store2.Count())
	}
}
