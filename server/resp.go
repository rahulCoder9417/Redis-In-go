package server

import(
	"bufio"
	"strconv"
)

type Resp struct{
	reader *bufio.Reader
}

func NewResp(reader *bufio.Reader)*Resp{
	return &Resp{
		reader: reader,
	}
}

func (r *Resp)readLine()(string,error){
	line,err :=r.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return line[:len(line)-2],nil
}


func (r *Resp) Read() ([]string, error) {

	line, err := r.readLine()

	if err != nil {
		return nil, err
	}

	if line[0] != '*' {
		return nil, nil
	}

	count, err := strconv.Atoi(line[1:])

	if err != nil {
		return nil, err
	}

	parts := make([]string, 0, count)

	for i := 0; i < count; i++ {

		_, err := r.readLine()

		if err != nil {
			return nil, err
		}

		arg, err := r.readLine()

		if err != nil {
			return nil, err
		}

		parts = append(parts, arg)
	}

	return parts, nil
}