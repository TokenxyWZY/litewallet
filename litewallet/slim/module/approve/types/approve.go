package types

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	btypes "github.com/QOSGroup/litewallet/litewallet/slim/base/types"
	"github.com/QOSGroup/litewallet/litewallet/slim/types"
)

// 授权 Common 结构
type Approve struct {
	From btypes.AccAddress `json:"from"` // 授权账号
	To   btypes.AccAddress `json:"to"`   // 被授权账号
	QOS  btypes.BigInt     `json:"qos"`  // QOS
	QSCs types.QSCs        `json:"qscs"` // QSCs
}

func NewApprove(from btypes.AccAddress, to btypes.AccAddress, qos btypes.BigInt, qscs types.QSCs) Approve {
	return Approve{
		From: from,
		To:   to,
		QOS:  qos.NilToZero(),
		QSCs: qscs,
	}
}

// 校验
func (approve Approve) Valid() error {
	// 地址不能为空
	if len(approve.From) == 0 {
		return ErrInvalidInput("from address is empty")
	}
	if len(approve.To) == 0 {
		return ErrInvalidInput("to address is empty")
	}
	// 地址不能相同
	if approve.From.Equals(approve.To) {
		return ErrInvalidInput("addresses of from and to are the same")
	}
	// 币值为正
	if !approve.IsPositive() {
		return ErrInvalidInput("coins must be positive")
	}
	// QSCs中不能包含QOS，不能包含重复币种
	m := make(map[string]bool)
	for _, val := range approve.QSCs {
		val.Name = strings.ToUpper(strings.TrimSpace(val.Name))
		// 不能包含QOS
		if strings.ToUpper(val.Name) == types.QOSCoinName {
			return ErrInvalidInput("QSCs can not contain qos, not case sensitive")
		}

		// 不能重复
		if _, ok := m[val.Name]; !ok {
			m[val.Name] = true
		} else {
			return ErrInvalidInput(fmt.Sprintf("repeat qsc:%s", val.Name))
		}
	}

	return nil
}

// 是否为正值
func (approve Approve) IsPositive() bool {
	if approve.QOS.IsNil() || approve.QOS.IsZero() {
		return approve.QSCs.IsPositive()
	} else if approve.QOS.GT(btypes.NewInt(0)) {
		return approve.IsNotNegative()
	} else {
		return false
	}
}

// 是否为非负值
func (approve Approve) IsNotNegative() bool {
	approve.QOS = approve.QOS.NilToZero()

	if approve.QOS.LT(btypes.NewInt(0)) {
		return false
	}

	return approve.QSCs.IsNotNegative()
}

// 返回相反值
func (approve Approve) Negative() Approve {

	return NewApprove(approve.From, approve.To, approve.QOS.Neg(), approve.QSCs.Negative())
}

// Plus
func (approve Approve) Plus(qos btypes.BigInt, qscs types.QSCs) Approve {

	return NewApprove(approve.From, approve.To, approve.QOS.Add(qos.NilToZero()), approve.QSCs.Plus(qscs))
}

// Minus
func (approve Approve) Minus(qos btypes.BigInt, qscs types.QSCs) Approve {
	approve.QOS = approve.QOS.NilToZero()
	qos = qos.NilToZero()

	return NewApprove(approve.From, approve.To, approve.QOS.Add(qos.Neg()), approve.QSCs.Minus(qscs))
}

// 是否大于等于
func (approve Approve) IsGTE(qos btypes.BigInt, qscs types.QSCs) bool {
	approve.QOS = approve.QOS.NilToZero()
	qos = qos.NilToZero()

	if approve.QOS.LT(qos) {
		return false
	}

	return approve.QSCs.IsGTE(qscs)
}

// 是否大于
func (approve Approve) IsGT(qos btypes.BigInt, qscs types.QSCs) bool {
	approve.QOS = approve.QOS.NilToZero()
	qos = qos.NilToZero()

	if approve.QOS.LT(qos) {
		return false
	} else if approve.QOS.Equal(qos) {
		return !approve.QSCs.IsLT(qscs) && !approve.QSCs.IsEqual(qscs)
	} else {
		return qscs.IsNotNegative()
	}
}

// 重写Equals
func (approve Approve) Equals(approveB Approve) bool {
	return approve.String() == approveB.String()
}

// 输出字符串
func (approve Approve) String() string {
	approve.QOS = approve.QOS.NilToZero()

	var buf bytes.Buffer
	buf.WriteString("from:" + approve.From.String() + " ")
	buf.WriteString("to:" + approve.To.String() + " ")
	buf.WriteString("qos:" + approve.QOS.String() + " ")
	names := make([]string, 0, len(approve.QSCs))
	m1 := make(map[string]btypes.BigInt)
	for _, val := range approve.QSCs {
		names = append(names, val.Name)
		m1[val.Name] = val.Amount
	}
	sort.Strings(names)
	for _, name := range names {
		buf.WriteString(name + ":")
		buf.WriteString(m1[name].String() + " ")
	}
	return buf.String()
}
