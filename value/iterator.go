package value

type Iterator interface {
    hasNext() bool
    next() Value
}

type RangeIterator struct {
    Range ValueRange
    Count float64
}

// impl Iterator for RangeIterator
func (r RangeIterator) hasNext() bool {
    return r.Count < *r.Range.End
}

func (r *RangeIterator) next() Value {
    r.Count += *r.Range.Step
}

