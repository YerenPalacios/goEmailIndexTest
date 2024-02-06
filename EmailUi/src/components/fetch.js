import { ref } from 'vue'

export function useFetch() {
	const data = ref(null)
	const error = ref(null)

	const fetchData = (url, options, serializer) => {
		data.value = null
		error.value = null

		fetch(url, options)
			.then((res) => {
				if (res.status >= 400) throw new Error('Error was found')
				return res.json()
			})
			.then((json) => {
				data.value = serializer(json)
			})
			.catch((err) => { error.value = err; console.log(err) })
	}

	const post = (url, body, serializer) => {
		fetchData(url, {
			method: 'POST',
			body: JSON.stringify(body),
			headers: {
				'Content-type': 'application/json',
				'Authorization': 'Basic YWRtaW46Q29tcGxleHBhc3MjMTIz'
			}
		}, serializer)
	}

	return { data, error, post }
}