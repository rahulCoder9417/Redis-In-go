package server

import "strconv"

func EncodeRESP(
	parts []string,
) string {

	resp :=
		"*" +
		strconv.Itoa(
			len(parts),
		) +
		"\r\n"

	for _, part := range parts {

		resp +=
			"$" +
			strconv.Itoa(
				len(part),
			) +
			"\r\n"

		resp +=
			part +
			"\r\n"
	}

	return resp
}