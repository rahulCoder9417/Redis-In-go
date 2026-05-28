package commands

import (
	"strconv"
	"strings"
	"github.com/rahulCoder9417/Redis-in-go/server/commands/utils"
)

func XAdd(parts []string) string {

	if len(parts) < 5 {
		return RespError("wrong number of arguments for 'XADD'")
	}

	if (len(parts)-3)%2 != 0 {
		return RespError("field/value pairs must be specified")
	}

	key := parts[1]
	id := parts[2]

	fields := make(map[string]string)

	for i := 3; i < len(parts); i += 2 {
		fields[parts[i]] = parts[i+1]
	}

	Mu.Lock()
	defer Mu.Unlock()

	v, exists := Store[key]

	if !exists {
		v = Value{
			Type:   "stream",
			Stream: []StreamEntry{},
		}
	}

	if v.Type != "stream" {
		return RespError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	lastID := ""

	if len(v.Stream) > 0 {
		lastID = v.Stream[len(v.Stream)-1].ID
	}

	// FULL AUTO ID
	if id == "*" {

		id = utils.GenerateStreamID(lastID)

	// PARTIAL AUTO ID
	} else if strings.HasSuffix(id, "-*") {

		msPart := strings.TrimSuffix(id, "-*")

		generatedID, ok := utils.GeneratePartialId(msPart, lastID)

		if !ok {
			return RespError("ID must be greater than previous ID")
		}

		id = generatedID

	// MANUAL ID
	} else {

		_, _, ok := utils.ParseStreamID(id)

		if !ok {
			return RespError("invalid stream entry ID")
		}

		if lastID != "" && utils.CompareIDs(id, lastID) <= 0 {
			return RespError("ID must be greater than previous ID")
		}
	}

	entry := StreamEntry{
		ID:     id,
		Fields: fields,
	}

	v.Stream = append(v.Stream, entry)

	Store[key] = v

	return RespBulkString(id)
}

func XRange(parts []string) string {
	if len(parts) != 4 {
		return RespError("wrong number of arguments for 'XRANGE'")
	}
	
	key := parts[1]
	start := parts[2]
	end := parts[3]
	Mu.RLock()
	value,exists:=Store[key]
	Mu.RUnlock()

	if !exists{
		return RespArray([]string{})
	}

	if value.Type != "stream" {
		return RespError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	result :=""
	count :=0

	if(start == "-"){
		start = "0-0"
	}
	if(end == "+"){
		end = "9999999999999-999999999"
	}

	for _,entry :=range value.Stream{
		if utils.CompareIDs(entry.ID, start) >= 0 && 
			utils.CompareIDs(entry.ID, end) <= 0 {
			
			fieldResp := ""
			fieldCount := 0

			for field,value := range entry.Fields{
				fieldResp += RespBulkString(field) + RespBulkString(value)
				fieldCount += 2
			}

			entryResp := "*2\r\n"

			entryResp+=RespBulkString(entry.ID) 

			entryResp += "*" + strconv.Itoa(fieldCount) + "\r\n"
			entryResp += fieldResp

			result += entryResp
			count++
		}
	}
	
	return "*" + strconv.Itoa(count) + "\r\n" + result
}
