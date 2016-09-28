package provingwork

import (
  "fmt"
	"bytes"
	"time"

	"math/big"

	"crypto/rand"
	"crypto/sha1"

	"encoding/base64"
	"encoding/binary"
)

var (
	DefaultBitStrength = 20
	DefaultSaltSize    = 16
)

// End goal format:
// 1:20:20160927155710:somedatatovalidate::aW5ZdXJQcm90b2NvbHMh:VvJC
// version, zero bits, date, resource, extension (ignored), rand, counter

type HashCashOptions struct {
	BitStrength int
	Extension   []byte
	Salt        []byte
	Timestamp   *time.Time
}

type HashCash struct {
	Resource []byte

	Counter int64

	*HashCashOptions
}

func NewHashCash(resource []byte, opts ...*HashCashOptions) *HashCash {
	hc := HashCash{
		Counter: 1,
		Resource: resource,
	}

	if len(opts) != 0 {
		hc.HashCashOptions = opts[0]
	} else {
		hc.HashCashOptions = &HashCashOptions{}
	}

	if hc.Timestamp == nil {
		t := time.Now()
		hc.Timestamp = &t
	}

	if hc.BitStrength == 0 {
		hc.BitStrength = DefaultBitStrength
	}

	if (len(hc.Salt) == 0) {
		hc.Salt = make([]byte, DefaultSaltSize)
		rand.Read(hc.Salt)
	}

	return &hc
}

func (hc *HashCash) Check() bool {
	if (hc.ZeroCount() >= hc.BitStrength) {
		return true
	}
	return false
}

func (hc *HashCash) CounterBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, &hc.Counter)
	return buf.Bytes()
}

func (hc *HashCash) FindProof() bool {
	for {
		if hc.Check() {
			return true
		}
		hc.Counter++
	}

	return false
}

func (hc *HashCash) String() string {
	return fmt.Sprintf(
		"1:%v:%v:%v:%v:%v:%v",
		hc.BitStrength,
		hc.Timestamp.Format("20060102150405"),
		string(hc.Resource),
		string(hc.Extension),
		base64.StdEncoding.EncodeToString(hc.Salt),
		base64.StdEncoding.EncodeToString(hc.CounterBytes()),
	)
}

func (hc *HashCash) ZeroCount() int {
	digest := sha1.Sum([]byte(hc.String()))
	digestHex := new(big.Int).SetBytes(digest[:])
	return (160 - digestHex.BitLen())
}
