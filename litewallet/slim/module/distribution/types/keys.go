package types

import (
	"encoding/binary"
	"fmt"
	btypes "github.com/QOSGroup/litewallet/litewallet/slim/base/types"
)

const (
	AddrLen = 20

	MapperName = "distribution"

	Distribution        = "distribution"
	ValidatorPeriodInfo = "validatorPeriodInfo"
	DelegatorIncomeInfo = "delegatorIncomeInfo"
)

var (
	//存储社区的QOS收益
	//value: bigint
	communityFeePoolKey = []byte{0x01}
	//上一块proposer地址
	//value: address
	lastBlockProposerKey = []byte{0x02}

	//每块待分配的QOS数量 = mint数量 + tx fee
	//value: bigint
	blockDistributionKey = []byte{0x04}

	//delegator收益计算信息,key = prefix+validatorAddr+delegatorAddr
	//value: delegatorEarningsStartInfo
	delegatorEarningsStartInfoPrefixKey = []byte{0x12}
	//validator历史计费点汇总收益,key = prefix + validatorAddr + period
	//value: bigint
	validatorHistoryPeriodSummaryPrefixKey = []byte{0x13}
	//validator当前计费点收益信息,key = prefix + validatorAddr
	//value: bigint
	validatorCurrentPeriodSummaryPrefixKey = []byte{0x14}

	//validator获得收益信息
	validatorEcoFeePoolPrefixKey = []byte{0x15}

	//delegators某高度下是否发放收益信息: key = prefix + blockheight + validatorAddress+delegatorAddress
	//value: true
	delegatorPeriodIncomePrefixKey = []byte{0x31}
)

func BuildCommunityFeePoolKey() []byte {
	return communityFeePoolKey
}

func BuildLastProposerKey() []byte {
	return lastBlockProposerKey
}

func BuildBlockDistributionKey() []byte {
	return blockDistributionKey
}

func GetValidatorCurrentPeriodSummaryPrefixKey() []byte {
	return validatorCurrentPeriodSummaryPrefixKey
}

func GetValidatorHistoryPeriodSummaryPrefixKey() []byte {
	return validatorHistoryPeriodSummaryPrefixKey
}

func GetDelegatorEarningsStartInfoPrefixKey() []byte {
	return delegatorEarningsStartInfoPrefixKey
}

func GetDelegatorPeriodIncomePrefixKey() []byte {
	return delegatorPeriodIncomePrefixKey
}

func GetValidatorEcoFeePoolPrefixKey() []byte {
	return validatorEcoFeePoolPrefixKey
}

func BuildDelegatorEarningStartInfoKey(validatorAddr btypes.ValAddress, delegatorAddress btypes.AccAddress) []byte {
	return append(append(delegatorEarningsStartInfoPrefixKey, validatorAddr...), delegatorAddress...)
}

func GetDelegatorEarningStartInfoAddr(key []byte) (valAddr btypes.ValAddress, deleAddr btypes.AccAddress) {
	if len(key) != (1 + 2*AddrLen) {
		panic("invalid delegatorEarningStartInfoKey length")
	}

	return btypes.ValAddress(key[1 : 1+AddrLen]), btypes.AccAddress(key[1+AddrLen:])
}

func BuildValidatorHistoryPeriodSummaryPrefixKey(validatorAddr btypes.ValAddress) []byte {
	return append(validatorHistoryPeriodSummaryPrefixKey, validatorAddr...)
}

func BuildValidatorHistoryPeriodSummaryKey(validatorAddr btypes.ValAddress, period int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(period))
	return append(append(validatorHistoryPeriodSummaryPrefixKey, validatorAddr...), b...)
}

func GetValidatorHistoryPeriodSummaryAddrPeriod(key []byte) (valAddr btypes.ValAddress, period int64) {
	if len(key) != (1 + 8 + AddrLen) {
		panic("invalid ValidatorHistoryPeriodSummaryKey length")
	}

	valAddr = btypes.ValAddress(key[1 : 1+AddrLen])
	b := key[1+AddrLen:]
	period = int64(binary.LittleEndian.Uint64(b))
	return
}

func BuildValidatorCurrentPeriodSummaryKey(validatorAddr btypes.ValAddress) []byte {
	return append(validatorCurrentPeriodSummaryPrefixKey, validatorAddr...)
}

func GetValidatorCurrentPeriodSummaryAddr(key []byte) btypes.ValAddress {
	if len(key) != (1 + AddrLen) {
		panic("invalid ValidatorCurrentPeriodSummaryKey length")
	}
	return btypes.ValAddress(key[1:])
}

func BuildDelegatorPeriodIncomeKey(validatorAddr btypes.ValAddress, delegatorAddress btypes.AccAddress, height int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(height))
	return append(append(append(delegatorPeriodIncomePrefixKey, b...), validatorAddr...), delegatorAddress...)
}

func BuildDelegatorPeriodIncomePrefixKey(height int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(height))
	return append(delegatorPeriodIncomePrefixKey, b...)
}

func BuildValidatorEcoFeePoolKey(validatorAddr btypes.ValAddress) []byte {
	return append(validatorEcoFeePoolPrefixKey, validatorAddr...)
}

func GetValidatorEcoPoolAddress(key []byte) btypes.ValAddress {
	return btypes.ValAddress(key[1 : 1+AddrLen])
}

func GetDelegatorPeriodIncomeHeightAddr(key []byte) (valAddr btypes.ValAddress, deleAddr btypes.AccAddress, height int64) {
	if len(key) != (1 + 8 + 2*AddrLen) {
		panic("invalid DelegatorsPeriodIncomeKey length")
	}

	b := key[1:9]
	return btypes.ValAddress(key[9 : 9+AddrLen]), btypes.AccAddress(key[9+AddrLen:]), int64(binary.LittleEndian.Uint64(b))
}

func BuildQueryValidatorPeriodInfoCustomQueryPath(valAddr btypes.ValAddress) string {
	return fmt.Sprintf("custom/%s/%s/%s", Distribution, ValidatorPeriodInfo, valAddr.String())
}

func BuildQueryDelegatorIncomeInfoCustomQueryPath(delegator btypes.AccAddress, valAddr btypes.ValAddress) string {
	return fmt.Sprintf("custom/%s/%s/%s/%s", Distribution, DelegatorIncomeInfo, delegator.String(), valAddr.String())
}
