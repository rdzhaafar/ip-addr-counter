package ipstore

const (
	nBuckets = 0xff_ff_ff
)

type IPStore struct {
	buckets [][]uint32
	count   int
}

func NewIPStore() *IPStore {
	buckets := make([][]uint32, nBuckets)
	for i := 0; i < nBuckets; i++ {
		buckets = append(buckets, nil)
	}
	return &IPStore{
		buckets: buckets,
		count:   0,
	}
}

func (s *IPStore) Insert(ip uint32) {
	hash := ipStoreHash(ip)
	bucket := s.buckets[hash]
	for _, stored := range bucket {
		if stored == ip {
			return
		}
	}
	s.buckets[hash] = append(s.buckets[hash], ip)
	s.count++
}

func (s *IPStore) Count() int {
	return s.count
}

func ipStoreHash(ip uint32) int {
	return int((ip & 0xff_ff_ff_00) >> 8)
}
