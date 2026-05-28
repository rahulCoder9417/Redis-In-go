package utils

import (
	"strconv"
	"strings"
	"time"
)

func ParseStreamID(id string) (int64, int64, bool) {
	parts := strings.SplitN(id, "-", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}

	ms, err1 := strconv.ParseInt(parts[0], 10, 64)
	seq, err2 := strconv.ParseInt(parts[1], 10, 64)

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
func GenerateStreamID(lastID string) string {

	currentMs := time.Now().UnixMilli()

	if lastID == "" {
		return strconv.FormatInt(currentMs, 10) + "-0"
	}

	lastMs, lastSeq, ok := ParseStreamID(lastID)

	if !ok {
		return strconv.FormatInt(currentMs, 10) + "-0"
	}

	// same millisecond
	if currentMs == lastMs {
		return strconv.FormatInt(currentMs, 10) + "-" + strconv.FormatInt(lastSeq+1, 10)
	}

	// newer millisecond
	if currentMs > lastMs {
		return strconv.FormatInt(currentMs, 10) + "-0"
	}

	// clock moved backwards
	return strconv.FormatInt(lastMs, 10) + "-" + strconv.FormatInt(lastSeq+1, 10)
}

func GeneratePartialId(msPart ,lastId string)(string,bool){
	ms,err := strconv.ParseInt(msPart, 10, 64)
	if err != nil {
		return "",false
	}

	if lastId ==""{
		return strconv.FormatInt(ms, 10) + "-0",true
	}

	lastMs, lastSeq, ok := ParseStreamID(lastId)
	if !ok {
		return strconv.FormatInt(ms, 10) + "-0",true
	}

	if ms < lastMs {
		return "",false
	}

	if ms == lastMs {
		return strconv.FormatInt(ms, 10) + "-" + strconv.FormatInt(lastSeq+1, 10),true
	}

	return strconv.FormatInt(ms, 10) + "-0",true
}
