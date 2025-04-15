package value

import (
	"vm-go/util"
)

// TODO: add ToList() when lists are in the language
type Iterator interface {
    HasNext() bool
    Advance()
    GetNext() Value
}

// ---

type RangeIterator struct {
    Range ValueRange
    Count float64
}

func NewRangeIterator(rg ValueRange) RangeIterator {
    return RangeIterator{
        Range: CopyValue(rg).(ValueRange), // to avoid race conditions while iterating
        Count: *rg.Start,
    }
}

// impl Iterator for *RangeIterator
func (r *RangeIterator) HasNext() bool {
	// Stop the iteration if the range is unreachable
    if !util.IsRangeReachable(*r.Range.Start, *r.Range.End, *r.Range.Step, *r.Range.Inclusive) {
        return false
    }

    if *r.Range.Inclusive {
        if *r.Range.Step > 0 {
            return r.Count <= *r.Range.End
        } else {
            return r.Count >= *r.Range.End
        }
    } else {
        if *r.Range.Step > 0 {
            return r.Count < *r.Range.End
        } else {
            return r.Count > *r.Range.End
        }
    }
}

func (r *RangeIterator) Advance() {
    r.Count += *r.Range.Step
}

func (r *RangeIterator) GetNext() Value {
    return ValueNumber{ Value: r.Count }
}

// impl Value for RangeIterator
func (x RangeIterator) String() string {
    return "<range iterator>"
}

func (x RangeIterator) Type() string {
    return "iterator"
}

// ---

type StrBytesIterator struct {
	Str string
	Pos int
}

func NewStrBytesIterator(str string) StrBytesIterator {
	return StrBytesIterator{
		Str: str,
		Pos: 0,
	}
}

// impl Iterator for *StrBytesIterator
func (r *StrBytesIterator) HasNext() bool {
	return r.Pos < len(r.Str)
}

func (r *StrBytesIterator) Advance() {
	r.Pos++
}

// TODO: add type char and make this return it
func (r *StrBytesIterator) GetNext() Value {
	return ValueString{ Value: string(r.Str[r.Pos]) }
}

// impl Value for StrBytesIterator
func (x StrBytesIterator) String() string {
    return "<string bytes iterator>"
}

func (x StrBytesIterator) Type() string {
    return "iterator"
}
