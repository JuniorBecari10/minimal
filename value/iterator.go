package value

type Iterator interface {
    HasNext() bool
    Advance()
    GetNext() Value
}

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

