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

func contains(slice []error, element error) bool {
	for _, value := range slice {
		if value.Error() == element.Error() {
			return true
		}
	}
	return false
}

func removeDuplication(list []error) []error {
	var newList []error

	for _, element := range list {
		if !contains(newList, element) {
			newList = append(newList, element)
		}
	}
	return newList
}
