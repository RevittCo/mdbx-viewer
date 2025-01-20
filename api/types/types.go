package types

import (
	"fmt"

	"github.com/RevittConsulting/mdbx-viewer/pkg/utils"
)

type Value []byte

type KeyValuePair struct {
	Key   Value `json:"key"`
	Value Value `json:"value"`
}

type KeyValuePairString struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (v Value) Hex() string {
	return fmt.Sprintf("0x%x", v)
}

func (v Value) Uint64() uint64 {
	return utils.BytesToUint64(v)
}

func (kv KeyValuePair) IntKeyHexValue() KeyValuePairString {
	return KeyValuePairString{
		Key:   fmt.Sprintf("%d", kv.Key.Uint64()),
		Value: kv.Value.Hex(),
	}
}

func (kv KeyValuePair) HexKeyHexValue() KeyValuePairString {
	return KeyValuePairString{
		Key:   kv.Key.Hex(),
		Value: kv.Value.Hex(),
	}
}

func (kv KeyValuePair) IntKeyIntValue() KeyValuePairString {
	return KeyValuePairString{
		Key:   fmt.Sprintf("%d", kv.Key.Uint64()),
		Value: fmt.Sprintf("%d", kv.Value.Uint64()),
	}
}

func (kv KeyValuePair) HexKeyIntValue() KeyValuePairString {
	return KeyValuePairString{
		Key:   kv.Key.Hex(),
		Value: fmt.Sprintf("%d", kv.Value.Uint64()),
	}
}
