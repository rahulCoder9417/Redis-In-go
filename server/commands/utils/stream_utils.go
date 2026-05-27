package utils

import (
	"strconv"
	"strings"
)

func ParseStreamID(id string) (int, int, bool) {
	parts := strings.SplitN(id, "-", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}
	
	ms, err1 := strconv.Atoi(parts[0])
	seq, err2 := strconv.Atoi(parts[1])
	
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	

	if ms == 0 && seq == 0 {
		return 0, 0, false
	}
	return ms, seq, true
}

func CompareIDs(a,b string)int{
	ms1, seq1, ok1 := ParseStreamID(a)
	ms2, seq2, ok2 := ParseStreamID(b)
	
	if !ok1 || !ok2 {
		return 0
	}
	
	if ms1 > ms2 {
		return 1
	}
	if ms1 < ms2 {
		return -1
	}
	
	if seq1 > seq2 {
		return 1
	}
	if seq1 < seq2 {
		return -1
	}
	
	return 0
}