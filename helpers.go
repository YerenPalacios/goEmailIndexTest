package main

func GetBatch(list []map[string]string, size int) [][]map[string]string {
	var grupos [][]map[string]string

	for i := 0; i < len(list); i += size {
		fin := i + size

		if fin > len(list) {
			fin = len(list)
		}

		grupos = append(grupos, list[i:fin])
	}

	return grupos
}

func Filter(numbers []string, condition func(string) bool) []string {
	var result []string

	for _, num := range numbers {
		if condition(num) {
			result = append(result, num)
		}
	}

	return result
}
