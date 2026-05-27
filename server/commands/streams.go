package commands

import "redis-from-scratch/server/commands/utils"

func XAdd(parts []string) string {
	if len(parts) < 5 {
		return RespError("wrong number of arguments for 'XADD'")
	}
	
	if (len(parts)-3)%2 != 0{
		return RespError("field/value pairs must be specified")
	}
	id:=parts[2]

	if(id != "*"){
		_,_,ok := utils.ParseStreamID(id)

		if !ok {
			return RespError("invalid stream entry ID")
		}
	}

	key := parts[1]

	fields := make(map[string]string)
	
	for i := 3; i < len(parts); i += 2 {
		fields[parts[i]] = parts[i+1]
	}
	
	Mu.Lock()
	defer Mu.Unlock()
	
	v,exists := Store[key]

	if !exists {
		v = Value{
			Type:"stream",
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

	if id == "*" {
		id = utils.GenerateStreamID(lastID)
	} else {
		if lastID != "" && utils.CompareIDs(id, lastID) <= 0 {
			return RespError("ID must be greater than previous ID")
		}
	}

	entry := StreamEntry{
		ID: id,
		Fields: fields,
	}
	v.Stream = append(v.Stream, entry)
	Store[key] = v
	
	return RespBulkString(id)
}
