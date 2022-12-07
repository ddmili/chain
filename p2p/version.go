package p2p

import (
	"bytes"
	"encoding/gob"
	"log"
)

//版本信息 默认0
const versionInfo = byte(0x00)

type version struct {
	Version  byte
	Height   int
	AddrFrom string
}

func (v version) serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(v)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func (v *version) deserialize(d []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(v)
	if err != nil {
		log.Panic(err)
	}
}
