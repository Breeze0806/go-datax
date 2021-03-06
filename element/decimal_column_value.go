package element

import (
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

//NilDecimalColumnValue 空值高精度实数型列值
type NilDecimalColumnValue struct {
	*nilColumnValue
}

//NewNilDecimalColumnValue 生成空值高精度实数型列值
func NewNilDecimalColumnValue() ColumnValue {
	return &NilDecimalColumnValue{
		nilColumnValue: &nilColumnValue{},
	}
}

//Type 列类型
func (n *NilDecimalColumnValue) Type() ColumnType {
	return TypeDecimal
}

//Clone 克隆
func (n *NilDecimalColumnValue) Clone() ColumnValue {
	return NewNilDecimalColumnValue()
}

//DecimalColumnValue 高精度实数列值
type DecimalColumnValue struct {
	notNilColumnValue

	val decimal.Decimal //高精度实数
}

//NewDecimalColumnValueFromFloat 根据float64 f生成高精度实数列值
func NewDecimalColumnValueFromFloat(f float64) ColumnValue {
	return &DecimalColumnValue{
		val: decimal.NewFromFloat(f),
	}
}

//NewDecimalColumnValue 根据高精度实数 d生成高精度实数列值
func NewDecimalColumnValue(d decimal.Decimal) ColumnValue {
	return &DecimalColumnValue{
		val: d,
	}
}

//NewDecimalColumnValueFromString 根据字符串s生成高精度实数列值
//不是数值型或者科学计数法的字符串，就会报错
func NewDecimalColumnValueFromString(s string) (ColumnValue, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil, err
	}
	return &DecimalColumnValue{
		val: d,
	}, nil
}

//Type 列类型
func (d *DecimalColumnValue) Type() ColumnType {
	return TypeDecimal
}

//AsBool 非0转化为true, 0转化为false
func (d *DecimalColumnValue) AsBool() (bool, error) {
	return d.val.Cmp(decimal.Zero) != 0, nil
}

//AsBigInt 对高精度实数取整，如123.67转化为123 123.12转化为123
func (d *DecimalColumnValue) AsBigInt() (*big.Int, error) {
	exp := d.val.Exponent()
	value := d.val.Coefficient()
	if exp == 0 {
		return value, nil
	}
	diff := math.Abs(-float64(exp))

	expScale := new(big.Int).Exp(_IntTen, big.NewInt(int64(diff)), nil)
	if 0 > exp {
		value = value.Quo(value, expScale)
	} else if 0 < exp {
		value = value.Mul(value, expScale)
	}

	return value, nil
}

//AsDecimal 转化为高精度实数
func (d *DecimalColumnValue) AsDecimal() (decimal.Decimal, error) {
	return d.val, nil
}

//AsString 转化为字符串， 如10.123 转化为10.123
func (d *DecimalColumnValue) AsString() (string, error) {
	return d.val.String(), nil
}

//AsBytes 转化为字节流， 如10.123 转化为10.123
func (d *DecimalColumnValue) AsBytes() ([]byte, error) {
	return []byte(d.val.String()), nil
}

//AsTime 目前无法转化为时间
func (d *DecimalColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, NewTransformErrorFormColumnTypes(d.Type(), TypeTime, fmt.Errorf(" val: %v", d.String()))
}

func (d *DecimalColumnValue) String() string {
	return d.val.String()
}

//Clone 克隆高精度实数列值
func (d *DecimalColumnValue) Clone() ColumnValue {
	return &DecimalColumnValue{
		val: d.val,
	}
}

//Cmp  返回1代表大于， 0代表相等， -1代表小于
func (d *DecimalColumnValue) Cmp(right ColumnValue) (int, error) {
	rightValue, err := right.AsDecimal()
	if err != nil {
		return 0, err
	}

	return d.val.Cmp(rightValue), nil
}
