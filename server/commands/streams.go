package commands

import (
	"strconv"
	"strings"
	"time"
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

	waiters := StreamWaiters[key]

	if len(waiters) > 0 {

		for _, waiter := range waiters {

			waiter <- entry
		}

		delete(StreamWaiters, key)
	}

	v.Stream = append(v.Stream, entry)

	Store[key] = v
	IncrementKeyVersion(key)
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

func XRead(parts []string) string {

	if len(parts) < 4 {
		return RespError(
			"wrong number of arguments for 'XREAD'",
		)
	}

	blocking := false
	timeout := 0
	offset := 1

	if strings.ToUpper(parts[1]) == "BLOCK" {

		blocking = true

		t, err := strconv.Atoi(parts[2])

		if err != nil {
			return RespError("invalid timeout")
		}

		timeout = t

		offset = 3
	}

	if strings.ToUpper(parts[offset]) != "STREAMS" {
		return RespError("syntax error")
	}

	remaining := parts[offset+1:]

	if len(remaining)%2 != 0 {
		return RespError(
			"unbalanced stream keys and IDs",
		)
	}

	streamCount := len(remaining) / 2

	streamNames := remaining[:streamCount]

	streamIDs := remaining[streamCount:]

	finalResp := ""

	matchedStreams := 0

	var waitKey string

	for i := 0; i < streamCount; i++ {

		key := streamNames[i]

		lastID := streamIDs[i]

		waitKey = key

		Mu.RLock()

		value, exists := Store[key]

		Mu.RUnlock()

		if !exists {
			continue
		}

		if value.Type != "stream" {
			return RespError(
				"WRONGTYPE Operation against wrong kind of value",
			)
		}

		if lastID == "$" {

			if len(value.Stream) == 0 {
				lastID = "0-0"
			} else {
				lastID = value.Stream[len(value.Stream)-1].ID
			}
		}

		var entries []StreamEntry

		for _, entry := range value.Stream {

			if utils.CompareIDs(
				entry.ID,
				lastID,
			) <= 0 {

				continue
			}

			entries =
				append(entries, entry)
		}

		if len(entries) == 0 {
			continue
		}

		finalResp += BuildStreamResponse(
			key,
			entries,
		)

		matchedStreams++
	}

	if matchedStreams > 0 {

		return "*" +
			strconv.Itoa(
				matchedStreams,
			) +
			"\r\n" +
			finalResp
	}

	if !blocking {
		return RespNull()
	}

	ch := make(chan StreamEntry)

	Mu.Lock()

	StreamWaiters[waitKey] =
		append(
			StreamWaiters[waitKey],
			ch,
		)

	Mu.Unlock()

	if timeout == 0 {

		entry := <-ch

		resp :=
			BuildStreamResponse(
				waitKey,
				[]StreamEntry{
					entry,
				},
			)

		return "*1\r\n" + resp
	}

	select {

	case entry := <-ch:

		resp :=
			BuildStreamResponse(
				waitKey,
				[]StreamEntry{
					entry,
				},
			)

		return "*1\r\n" + resp

	case <-time.After(
		time.Duration(timeout) *
			time.Millisecond,
	):

		RemoveStreamWaiter(
			waitKey,
			ch,
		)

		return RespNull()
	}
}


func BuildEntryResponse(entry StreamEntry) string {

	fieldResp := ""
	fieldCount := 0

	for field, value := range entry.Fields {

		fieldResp += RespBulkString(field)
		fieldResp += RespBulkString(value)

		fieldCount += 2
	}

	resp := "*2\r\n"

	resp += RespBulkString(entry.ID)

	resp += "*" + strconv.Itoa(fieldCount) + "\r\n"

	resp += fieldResp

	return resp
}

func BuildStreamResponse(
	key string,
	entries []StreamEntry,
) string {

	resp := "*2\r\n"

	resp += RespBulkString(key)

	resp += "*" +
		strconv.Itoa(len(entries)) +
		"\r\n"

	for _, entry := range entries {
		resp += BuildEntryResponse(entry)
	}

	return resp
}