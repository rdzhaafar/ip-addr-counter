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

func TestIPStoreBucketDistribution(t *testing.T) {
	store := NewIPStore()
	for i := uint32(0); i < 0x00_00_ff_00; i += 0x00_00_01_00 {
		store.Insert(i)
	}
	for i := 0; i < 0xff; i++ {
		bucket := store.buckets[i]
		if !bucket.has(bitmap{1, 0, 0, 0}) {
			t.Fatalf("store.buckets[%d]: %s", i, bucket)
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

func TestBitmapApply(t *testing.T) {
	bm1 := bitmap{0, 0, 0, 1}
	mask1 := bitmap{1, 0, 0, 0}
	res1 := bm1.apply(mask1)
	if res1[0] != 1 || res1[1] != 0 || res1[2] != 0 || res1[3] != 1 {
		t.Errorf("Bitmap apply mask failed: res1[0]: %b res1[1]: %b res1[2]: %b res1[3]: %b", res1[0], res1[1], res1[2], res1[3])
	}

	bm2 := bitmap{1, 1, 1, 1}
	mask2 := bitmap{1 << 63, 1 << 63, 1 << 63, 1 << 63}
	res2 := bm2.apply(mask2)
	for i, b := range res2 {
		if b != 0x8000000000000001 {
			t.Errorf("bitmap.apply() failed: res2[%d]: %b", i, b)
		}
	}

	bm3 := bitmap{1, 0, 0, 0}
	mask3 := bitmap{1, 0, 0, 0}
	res3 := bm3.apply(mask3)
	if res3[0] != 1 || res3[1] != 0 || res3[2] != 0 || res3[3] != 0 {
		t.Errorf("bitmap.apply() failed: res3[0]: %b res3[1]: %b res3[2]: %b res3[3]: %b", res3[0], res3[1], res3[2], res3[3])
	}
}

func TestBitmaskHas(t *testing.T) {
	bm1 := bitmap{0, 0, 0, 0}
	mask1 := bitmap{1, 0, 0, 0}
	if bm1.has(mask1) {
		t.Errorf("bitmap.has() failed: bitmap: %v mask: %s", bm1, mask1)
	}

	bm2 := bitmap{1, 0, 0, 0}
	mask2 := bitmap{1, 0, 0, 0}
	if !bm2.has(mask2) {
		t.Errorf("bitmap.has() failed: bitmap: %s mask: %s", bm2, mask2)
	}

	bm3 := bitmap{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64}
	mask3 := bitmap{0, 0, 1, 0}
	if !bm3.has(mask3) {
		t.Errorf("bitmap.has() failed: bitmap: %s mask %s", bm3, mask3)
	}
}

func TestIPToBitmap(t *testing.T) {
	var cases = []struct {
		ip uint32
		bm bitmap
	}{
		{0, bitmap{1, 0, 0, 0}},
		{0xff_ff_ff_00, bitmap{1, 0, 0, 0}},
		{math.MaxUint32, bitmap{0, 0, 0, 0x8000000000000000}},
		{1, bitmap{2, 0, 0, 0}},
		{0x80, bitmap{0, 0, 1, 0}},
		{0xff, bitmap{0, 0, 0, 1 << 63}},
		{0xc0, bitmap{0, 0, 0, 1}},
	}

	for _, tc := range cases {
		bm := ipToBitmap(tc.ip)
		for i, b := range bm {
			if b != tc.bm[i] {
				t.Errorf("ipToBitmap(%x) failed: want: %s got: %s", tc.ip, tc.bm, bm)
				break
			}
		}
	}
}
