package types

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	btypes "github.com/QOSGroup/litewallet/litewallet/slim/base/types"
)

const (
	AddrLen = 20

	MapperName = "validator"

	//------query-------
	Stake         = "stake"
	Delegation    = "delegation"
	Delegations   = "delegations"
	ValidatorFlag = "validator"
	Delegator     = "delegator"
	Unbondings    = "Unbondings"
	Redelegations = "Redelegations"
)

var (
	//keys see docs/spec/staking.md
	validatorKey            = []byte{0x01} // 保存Validator信息. key: OperatorAddress
	validatorByConsensusKey = []byte{0x02} // 保存consensus address与Validator的映射关系. key: consensusAddress, value : OperatorAddress

	validatorByInactiveKey  = []byte{0x04} // 保存处于`inactive`状态的Validator. key: ValidatorInactiveTime + OperatorAddress
	validatorByVotePowerKey = []byte{0x05} // 按VotePower排序的Validator地址,不包含`pending`状态的Validator. key: VotePower + OperatorAddress

	//keys see docs/spec/staking.md
	validatorVoteInfoKey         = []byte{0x11} // 保存Validator在窗口的统计信息
	validatorVoteInfoInWindowKey = []byte{0x12} // 保存Validator在指定窗口签名信息

	DelegationByDelValKey = []byte{0x31} // key: delegator add + validator OperatorAddress add, value: delegationInfo
	DelegationByValDelKey = []byte{0x32} // key: OperatorAddress owner add + delegator add, value: nil

	UnbondingHeightDelegatorValidatorKey = []byte{0x41} // key: height + delegator + validator OperatorAddress addr, value: the amount of qos going to be unbonded on this height
	UnbondingDelegatorHeightValidatorKey = []byte{0x42} // key: delegator + height + validator OperatorAddress addr, value: nil
	UnbondingValidatorHeightDelegatorKey = []byte{0x43} // key: validator + height + delegator add, value: nil

	RedelegationHeightDelegatorFromValidatorKey = []byte{0x51} // key: height + delegator + fromValidator add, value: redelegations going to be complete on this height
	RedelegationDelegatorHeightFromValidatorKey = []byte{0x52} // key: delegator + height + fromValidator add, value: nil
	RedelegationFromValidatorHeightDelegatorKey = []byte{0x53} // key: fromValidator + height + delegator add, value: nil

	currentValidatorsAddressKey = []byte("currentValidatorsAddressKey")
)

func BuildStakeStoreQueryPath() []byte {
	return []byte(fmt.Sprintf("/store/%s/key", MapperName))
}

func BuildCurrentValidatorsAddressKey() []byte {
	return currentValidatorsAddressKey
}

func BuildValidatorKey(valAddress btypes.ValAddress) []byte {
	return append(validatorKey, valAddress...)
}

func BuildValidatorPrefixKey() []byte {
	return validatorKey
}

func BuildValidatorByConsensusKey(consensusAddress btypes.ConsAddress) []byte {

	lenz := 1 + len(consensusAddress)
	bz := make([]byte, lenz)

	copy(bz[0:1], validatorByConsensusKey)
	copy(bz[1:len(consensusAddress)+1], consensusAddress)

	return bz
}

func BuildInactiveValidatorKeyByTime(inactiveTime time.Time, valAddress btypes.ValAddress) []byte {
	return BuildInactiveValidatorKey(inactiveTime.UTC().Unix(), valAddress)
}

func BuildInactiveValidatorKey(sec int64, valAddress btypes.ValAddress) []byte {
	lenz := 1 + 8 + len(valAddress)
	bz := make([]byte, lenz)

	secBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(secBytes, uint64(sec))

	copy(bz[0:1], validatorByInactiveKey)
	copy(bz[1:9], secBytes)
	copy(bz[9:len(valAddress)+9], valAddress)

	return bz
}

func GetValidatorByInactiveKey() []byte {
	return validatorByInactiveKey
}

func GetValidatorByVotePowerKey() []byte {
	return validatorByVotePowerKey
}

func GetValidatorVoteInfoInWindowKey() []byte {
	return validatorVoteInfoInWindowKey
}

func GetValidatorVoteInfoKey() []byte {
	return validatorVoteInfoKey
}

func BuildValidatorByVotePower(votePower int64, valAddress btypes.ValAddress) []byte {
	lenz := 1 + 8 + len(valAddress)
	bz := make([]byte, lenz)

	votePowerBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(votePowerBytes, uint64(votePower))

	copy(bz[0:1], validatorByVotePowerKey)
	copy(bz[1:9], votePowerBytes)
	copy(bz[9:len(valAddress)+9], valAddress)

	return bz
}

func ParseValidatorVotePowerKey(key []byte) (votePower int64, valAddress btypes.ValAddress, err error) {
	if len(key) < 10 {
		return 0, btypes.ValAddress{}, errors.New("incorrect key length")
	}

	if key[0] != validatorByVotePowerKey[0] {
		return 0, btypes.ValAddress{}, errors.New("incorrect key type, not validatorByVotePowerKey key")
	}

	votePower = int64(binary.BigEndian.Uint64(key[1:9]))
	valAddress = btypes.ValAddress(key[9:])
	err = nil
	return
}

func BuildDelegationByDelValKey(delAdd btypes.AccAddress, valAdd btypes.ValAddress) []byte {
	bz := append(DelegationByDelValKey, delAdd...)
	return append(bz, valAdd...)
}

func BuildDelegationByValDelKey(valAdd btypes.ValAddress, delAdd btypes.AccAddress) []byte {
	bz := append(DelegationByValDelKey, valAdd...)
	return append(bz, delAdd...)
}

func BuildDelegationByValidatorPrefix(valAdd btypes.ValAddress) []byte {
	return append(DelegationByValDelKey, valAdd...)
}

func GetDelegationValDelKeyAddress(key []byte) (valAddr btypes.ValAddress, deleAddr btypes.AccAddress) {
	if len(key) != 1+2*AddrLen {
		panic("invalid DelegationValDelKey length")
	}

	valAddr = key[1 : 1+AddrLen]
	deleAddr = key[1+AddrLen:]
	return
}

func BuildValidatorVoteInfoKey(valAddress btypes.ValAddress) []byte {
	return append(validatorVoteInfoKey, valAddress...)
}

func BuildValidatorVoteInfoInWindowPrefixKey(valAddress btypes.ValAddress) []byte {
	return append(validatorVoteInfoInWindowKey, valAddress...)
}

func GetValidatorVoteInfoAddr(key []byte) btypes.ValAddress {
	return btypes.ValAddress(key[1:])
}

func BuildValidatorVoteInfoInWindowKey(index int64, valAddress btypes.ValAddress) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(index))

	bz := append(validatorVoteInfoInWindowKey, valAddress...)
	bz = append(bz, b...)

	return bz
}

func GetValidatorVoteInfoInWindowIndexAddr(key []byte) (int64, btypes.ValAddress) {
	addr := btypes.ValAddress(key[1 : AddrLen+1])
	index := int64(binary.LittleEndian.Uint64(key[AddrLen+1:]))
	return index, addr
}

func BuildUnbondingHeightDelegatorValidatorKey(height int64, deleAddr btypes.AccAddress, valAddr btypes.ValAddress) []byte {
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(height))

	bz := append(UnbondingHeightDelegatorValidatorKey, heightBytes...)
	bz = append(bz, deleAddr...)
	return append(bz, valAddr...)
}

func BuildUnbondingDelegationByHeightPrefix(height int64) []byte {
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(height))

	return append(UnbondingHeightDelegatorValidatorKey, heightBytes...)
}

func GetUnbondingDelegationHeightDelegatorValidator(key []byte) (height int64, deleAddr btypes.AccAddress, valAddr btypes.ValAddress) {

	if len(key) != (1 + 8 + 2*AddrLen) {
		panic("invalid UnbondingHeightDelegatorKey length")
	}

	height = int64(binary.BigEndian.Uint64(key[1:9]))
	deleAddr = key[9 : AddrLen+9]
	valAddr = key[AddrLen+9:]
	return
}

func GetUnbondingDelegationDelegatorHeightValidator(key []byte) (deleAddr btypes.AccAddress, height int64, valAddr btypes.ValAddress) {

	if len(key) != (1 + 8 + 2*AddrLen) {
		panic("invalid UnbondingDelegatorHeightKey length")
	}

	deleAddr = key[1 : AddrLen+1]
	height = int64(binary.BigEndian.Uint64(key[AddrLen+1 : AddrLen+9]))
	valAddr = key[AddrLen+9:]
	return
}

func BuildUnbondingDelegatorHeightValidatorKey(delAddr btypes.AccAddress, height int64, valAddr btypes.ValAddress) []byte {
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(height))

	bz := append(UnbondingDelegatorHeightValidatorKey, delAddr...)
	bz = append(bz, heightBytes...)
	return append(bz, valAddr...)
}

func BuildUnbondingByDelegatorPrefix(delAddr btypes.AccAddress) []byte {

	return append(UnbondingDelegatorHeightValidatorKey, delAddr...)
}

func BuildUnbondingValidatorHeightDelegatorKey(valAddr btypes.ValAddress, height int64, delAddr btypes.AccAddress) []byte {
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(height))

	bz := append(UnbondingDelegatorHeightValidatorKey, valAddr...)
	bz = append(bz, heightBytes...)
	return append(bz, delAddr...)
}

func BuildUnbondingByValidatorPrefix(valAddr btypes.ValAddress) []byte {

	return append(UnbondingValidatorHeightDelegatorKey, valAddr...)
}

func GetUnbondingDelegationValidatorHeightDelegator(key []byte) (valAddr btypes.ValAddress, height int64, deleAddr btypes.AccAddress) {

	if len(key) != (1 + 8 + 2*AddrLen) {
		panic("invalid UnbondingDelegatorHeightKey length")
	}

	valAddr = key[1 : AddrLen+1]
	height = int64(binary.BigEndian.Uint64(key[AddrLen+1 : AddrLen+9]))
	deleAddr = key[AddrLen+9:]
	return
}

func BuildRedelegationHeightDelegatorFromValidatorKey(height int64, delAdd btypes.AccAddress, valAddr btypes.ValAddress) []byte {
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(height))

	bz := append(RedelegationHeightDelegatorFromValidatorKey, heightBytes...)
	bz = append(bz, delAdd...)
	return append(bz, valAddr...)
}

func BuildRedelegationByHeightPrefix(height int64) []byte {
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(height))

	return append(RedelegationHeightDelegatorFromValidatorKey, heightBytes...)
}

func GetRedelegationHeightDelegatorFromValidator(key []byte) (height int64, deleAddr btypes.AccAddress, valAddr btypes.ValAddress) {

	if len(key) != (1 + 8 + 2*AddrLen) {
		panic("invalid RedelegationHeightDelegatorKey length")
	}

	height = int64(binary.BigEndian.Uint64(key[1:9]))
	deleAddr = key[9 : AddrLen+9]
	valAddr = key[AddrLen+9:]
	return
}

func GetRedelegationDelegatorHeightFromValidator(key []byte) (deleAddr btypes.AccAddress, height int64, valAddr btypes.ValAddress) {

	if len(key) != (1 + 8 + 2*AddrLen) {
		panic("invalid RedelegationDelegatorHeightKey length")
	}

	deleAddr = key[1 : AddrLen+1]
	height = int64(binary.BigEndian.Uint64(key[AddrLen+1 : AddrLen+9]))
	valAddr = key[AddrLen+9:]
	return
}

func GetRedelegationFromValidatorHeightDelegator(key []byte) (valAddr btypes.ValAddress, height int64, deleAddr btypes.AccAddress) {

	if len(key) != (1 + 8 + 2*AddrLen) {
		panic("invalid RedelegationDelegatorHeightKey length")
	}

	valAddr = key[1 : AddrLen+1]
	height = int64(binary.BigEndian.Uint64(key[AddrLen+1 : AddrLen+9]))
	deleAddr = key[AddrLen+9:]
	return
}

func BuildRedelegationDelegatorHeightFromValidatorKey(delAddr btypes.AccAddress, height int64, valAddr btypes.ValAddress) []byte {
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(height))

	bz := append(RedelegationDelegatorHeightFromValidatorKey, delAddr...)
	bz = append(bz, heightBytes...)
	return append(bz, valAddr...)
}

func BuildRedelegationByDelegatorPrefix(delAddr btypes.AccAddress) []byte {

	return append(RedelegationDelegatorHeightFromValidatorKey, delAddr...)
}

func BuildRedelegationFromValidatorHeightDelegatorKey(valAddr btypes.ValAddress, height int64, delAddr btypes.AccAddress) []byte {
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(height))

	bz := append(RedelegationDelegatorHeightFromValidatorKey, valAddr...)
	bz = append(bz, heightBytes...)
	return append(bz, delAddr...)
}

func BuildRedelegationByFromValidatorPrefix(valAddr btypes.ValAddress) []byte {

	return append(RedelegationFromValidatorHeightDelegatorKey, valAddr...)
}

//-------------------------query path

func BuildGetDelegationCustomQueryPath(deleAddr btypes.AccAddress, valAddr btypes.ValAddress) string {
	return fmt.Sprintf("custom/%s/%s/%s/%s", Stake, Delegation, deleAddr.String(), valAddr.String())
}

func BuildQueryDelegationsByOwnerCustomQueryPath(valAddr btypes.ValAddress) string {
	return fmt.Sprintf("custom/%s/%s/%s/%s", Stake, Delegations, ValidatorFlag, valAddr.String())
}

func BuildQueryDelegationsByDelegatorCustomQueryPath(deleAddr btypes.AccAddress) string {
	return fmt.Sprintf("custom/%s/%s/%s/%s", Stake, Delegations, Delegator, deleAddr.String())
}

func BuildQueryUnbondingsByDelegatorCustomQueryPath(deleAddr btypes.AccAddress) string {
	return fmt.Sprintf("custom/%s/%s/%s", Stake, Unbondings, deleAddr.String())
}

func BuildQueryRedelegationsByDelegatorCustomQueryPath(deleAddr btypes.AccAddress) string {
	return fmt.Sprintf("custom/%s/%s/%s", Stake, Redelegations, deleAddr.String())
}
