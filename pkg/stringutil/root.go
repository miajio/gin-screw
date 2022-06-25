package stringutil

const (
	NoBreakSp       = '\u00A0' // no-break space
	NarrowNoBreakSp = '\u202F' // narrow no-break space
	WideNoBreakSp   = '\uFEFF' // wide no-break space
	EnNoBreakSp     = '\u0020' // no-break space
	ZhWideNoBreakSp = '\u3000' // wide no-break space

	LeftToRightEmbedding = '\u202a' // left-to-right embedding
)

// DeSpace delete space in the val
func DeSpace(val string) string {
	result := make([]rune, 0)
	for _, v := range val {
		if v == NoBreakSp || v == NarrowNoBreakSp || v == WideNoBreakSp || v == EnNoBreakSp || v == ZhWideNoBreakSp {
			continue
		}
		result = append(result, v)
	}
	return string(result)
}
