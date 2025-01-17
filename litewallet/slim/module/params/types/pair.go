package types

import (
	btypes "github.com/QOSGroup/litewallet/litewallet/slim/base/types"
)

type KeyValuePair struct {
	Key   []byte
	Value interface{}
}

type KeyValuePairs []KeyValuePair

type ParamSet interface {
	KeyValuePairs() KeyValuePairs
	Validate(key string, value string) (interface{}, btypes.Error)
	GetParamSpace() string
}
