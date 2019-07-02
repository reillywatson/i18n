// Copyright (c) 2011 Jad Dittmar
// See: https://github.com/Confunctionist/finance
//
// Some changes by Oliver Eilhard
package i18n

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
)

type Money struct {
	M int64
	C CurrencyCode
}

var (
	ErrMoneyOverflow              = errors.New("i18n: money overflow")
	ErrMoneyDivideByZero          = errors.New("i18n: money division by zero")
	ErrMoneyDecimalPlacesTooLarge = errors.New("i18n: money decimal places too large")
	ErrMoneyCannotSplit           = errors.New("i18n: cannot split money into more chunks than cents in total")
	ErrMoneyZeroOrLessChunks      = errors.New("i18n: cannot split money into zero or less chunks")
)

const (
	MAXDEC = 18
	Round  = .5
	Roundn = Round * -1
)

type moneyMarshalContainer struct {
	M int64        `json:"M"`
	C CurrencyCode `json:"C"`
	F float64      `json:"F"`
}

func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(moneyMarshalContainer{M: m.M, C: m.C, F: m.Get()})
}

func (m *Money) UnmarshalJSON(b []byte) error {
	var container moneyMarshalContainer
	err := json.Unmarshal(b, &container)
	if err != nil {
		return err
	}
	if container.M != 0 && container.F != 0 && (Money{M: container.M, C: container.C}).Get() != container.F {
		return errors.New("M and F were both specified, but they aren't equivalent")
	}
	m.C = container.C
	m.M = container.M
	if m.M == 0 && container.F != 0 {
		*m = MakeMoney(m.C, container.F)
	}
	return nil
}

func MakeMoney(currency CurrencyCode, amount float64) Money {
	dpf := Money{C: currency}.dpf()
	fDPf := amount * dpf
	r := int64(amount * dpf)
	return Money{C: currency, M: rnd(r, fDPf-float64(r))}
}

// Returns the absolute value of Money.
func (m Money) Abs() Money {
	if m.M < 0 {
		return m.Neg()
	}
	return m
}

// Adds two money types.
func (m Money) Add(n Money) Money {
	r := m.M + n.M
	if (r^m.M)&(r^n.M) < 0 {
		panic(ErrMoneyOverflow)
	}
	return Money{C: m.C, M: r}
}

// Divides one Money type from another.
func (m Money) Div(n Money) Money {
	if n.M == 0 {
		panic(ErrMoneyDivideByZero)
	}
	guardf := 100.0
	f := guardf * m.dpf() * float64(m.M) / float64(n.M) / guardf
	i := int64(f)
	return Money{C: m.C, M: rnd(i, f-float64(i))}
}

// Splits the money amount into equal parts, up to 1 cent variation to not lose cents
// The larger cent amounts go in the beginning of the array
// chunks - how many parts to split money into, cannot exceed total cents in money amount
func (m Money) Split(chunks int64) []Money {
	if chunks <= 0 {
		panic(ErrMoneyZeroOrLessChunks)
	}
	if chunks > m.M {
		panic(ErrMoneyCannotSplit)
	}

	chunkAmount := m.M / chunks
	remainder := m.M % chunks //remainder cannot be > chunks

	var result []Money
	for i := int64(0); i < chunks; i++ {
		result = append(result, Money{M: chunkAmount, C: m.C})
	}

	//give out missing cents starting from the top
	for i := int64(0); i < remainder; i++ {
		result[i] = result[i].Add(Money{M: 1, C: m.C})
	}

	return result
}

// rounding value, expressed as 10^N where N is the number of decimal places
// (ie 2 decimals places == 10^2 == 100)
func (m Money) dp() int64 {
	for _, loc := range Locales {
		if loc.CurrencyCode == m.C {
			return int64(math.Pow10(loc.CurrencyDecimalDigits))
		}
	}
	return 100
}

// number of decimal places to use for rounding (float verison)
func (m Money) dpf() float64 {
	return float64(m.dp())
}

// Gets value of money truncating after DP (see Value() for no truncation).
func (m Money) Gett() int64 {
	return m.M / m.dp()
}

// Gets the float64 value of money (see Value() for int64).
func (m Money) Get() float64 {
	return float64(m.M) / m.dpf()
}

// Multiplies two Money types.
func (m Money) Mul(n Money) Money {
	return m.Mulf(n.Get())
}

// Multiplies a Money with a float to return a money-stored type.
func (m Money) Mulf(f float64) Money {
	mf := m.Get()
	return MakeMoney(m.C, f*mf)
}

// Returns the negative value of Money.
func (m Money) Neg() Money {
	if m.M != 0 {
		return Money{C: m.C, M: m.M * -1}
	}
	return m
}

// Rounds int64 remainder rounded half towards plus infinity
// trunc = the remainder of the float64 calc
// r     = the result of the int64 cal
func rnd(r int64, trunc float64) int64 {
	if trunc > 0 {
		if trunc >= Round {
			r++
		}
	} else {
		if trunc < Roundn {
			r--
		}
	}
	return r
}

// Returns the Sign of Money 1 if positive, -1 if negative.
func (m Money) Sign() int {
	if m.M < 0 {
		return -1
	}
	return 1
}

// String for money type representation in basic monetary unit (DOLLARS CENTS).
func (m Money) String() string {
	if m.Sign() > 0 {
		return fmt.Sprintf("%d.%02d %s", m.Value()/m.dp(), m.Value()%m.dp(), m.C)
	}
	// Negative value
	return fmt.Sprintf("-%d.%02d %s", m.Abs().Value()/m.dp(), m.Abs().Value()%m.dp(), m.C)
}

func (m Money) Format(locale string) string {
	l, found := Locales[locale]
	if !found {
		// If we don't have any information about the currency format,
		// we'll try our best to display something useful.
		return m.String()
	}

	// DP is a measure for decimals: 2 decimal digits => dp = 10^2
	currencySymbol := string(m.C)
	curr, found := Currencies[m.C]
	if found {
		currencySymbol = curr.Symbol
	}

	// DP is a measure for decimals: 2 decimal digits => dp = 10^2
	dp := int64(math.Pow10(l.CurrencyDecimalDigits))

	// Group DP is a measure for grouping: 3 decimal digits => groupDp = 10^3
	var groupDp int64
	var groupSize int
	if len(l.CurrencyGroupSizes) == 0 {
		// BUG(oe): Handle currency group size
		groupDp = int64(math.Pow10(3))
		groupSize = 3
	} else if len(l.CurrencyGroupSizes) >= 1 {
		// BUG(oe): Handle currency group size
		groupDp = int64(math.Pow10(l.CurrencyGroupSizes[0]))
		groupSize = l.CurrencyGroupSizes[0]
	}

	// We use absolute values (as int64) from here on, because the
	// negative sign is part of the currency format pattern.
	absVal := m.Value()
	if m.Sign() < 0 {
		absVal = -absVal
	}
	wholeVal := absVal / dp
	decVal := absVal % dp

	// The unformatted string (without grouping and with a decimal sep of ".")
	var unformatted string
	if l.CurrencyDecimalDigits > 0 {
		unformatted = fmt.Sprintf("%d.%0"+fmt.Sprintf("%d", l.CurrencyDecimalDigits)+"d", wholeVal, decVal)
	} else {
		unformatted = fmt.Sprintf("%d", wholeVal)
	}

	// Perform grouping operation of the whole number
	// For 1234, this returns this array: [234 1]
	groups := make([]string, 0)
	for {
		if groupDp > wholeVal {
			// do not prepend zeros
			groups = append(groups, fmt.Sprintf("%d", wholeVal%groupDp))
		} else {
			// prepend zeros
			groups = append(groups, fmt.Sprintf("%0"+fmt.Sprintf("%d", groupSize)+"d", wholeVal%groupDp))
		}
		wholeVal /= groupDp
		if wholeVal == 0 {
			break
		}
	}
	var wholeBuf bytes.Buffer
	for i := range groups {
		if i > 0 {
			wholeBuf.WriteString(l.CurrencyGroupSeparator)
		}
		wholeBuf.WriteString(groups[len(groups)-i-1])
	}

	// Which pattern do we need?
	// Notice that the minus sign is part of the pattern
	var pattern string
	if m.Sign() > 0 {
		pattern = l.CurrencyPositivePattern
	} else {
		pattern = l.CurrencyNegativePattern
	}

	// Split into whole and decimal and build formatted number
	var formatted string
	parts := strings.SplitN(unformatted, ".", 2)
	if len(parts) > 1 {
		formatted = fmt.Sprintf("%s%s%s", wholeBuf.String(), l.CurrencyDecimalSeparator, parts[1])
	} else {
		formatted = wholeBuf.String()
	}

	output := strings.Replace(pattern, "$", currencySymbol, -1)
	output = strings.Replace(output, "n", formatted, -1)

	return output
}

// Subtracts one Money type from another.
func (m Money) Sub(n Money) Money {
	r := m.M - n.M
	if (r^m.M)&^(r^n.M) < 0 {
		panic(ErrMoneyOverflow)
	}
	return Money{C: m.C, M: r}
}

// Returns in int64 the value of Money (also see Gett(), See Get() for float64).
func (m Money) Value() int64 {
	return m.M
}
