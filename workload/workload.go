package workload

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type DataBlob struct {
	Data string   `json:"data"`
	Hash [32]byte `json:"hash"`
}

func NewDataBlobFromString(s string) *DataBlob {
	d := &DataBlob{
		Data: s,
	}
	d.ComputeHash()
	return d
}

func NewDataBlobFromBytes(b []byte) (*DataBlob, error) {
	d := &DataBlob{
		Data: "",
		Hash: [32]byte{},
	}
	if err := d.UnMarshal(b); err != nil {
		return nil, err
	}
	if err := d.VerifyHash(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DataBlob) ComputeHash() {
	d.Hash = sha256.Sum256([]byte(d.Data))
}

func (d *DataBlob) VerifyHash() error {
	h := sha256.Sum256([]byte(d.Data))
	for i, _ := range h {
		if h[i] != d.Hash[i] {
			return fmt.Errorf("hash mismatch for entry at index %d", i)
		}
	}
	return nil
}

func (d *DataBlob) Marshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d *DataBlob) UnMarshal(b []byte) error {
	return json.Unmarshal(b, d)
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func randStr(length int) string {
	log.Infof("length = %d", length)
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
