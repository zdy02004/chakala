package main

import "strings"

func sql_low_inject(in string) string {
	low_level := [...]string{"select", "and", "or", "like", "regxp", "from", "where", "update", "exec", "order by", "having", "drop", "delete"}

	var index int = -1

	for _, k := range low_level {
		index = strings.Index(strings.ToLower(in), k)
		if index >= 0 {
			return k
		}
	}

	return "ok"

}

func sql_high_inject(in string) string {
	hight_level := [...]string{"select", "and", "or", "like", "regxp", "from", "where", "update", "exec", "order by", "having", "drop", "delete", "(", ")", "[", "]", "<", ">", ",", ".", ";", ":", "'", "\"", "#", "%", "+", "-", "_", "=", "/", "*", "@"}

	var index int = -1

	for _, k := range hight_level {
		index = strings.Index(strings.ToLower(in), k)
		if index >= 0 {
			return k
		}
	}

	return "ok"

}
