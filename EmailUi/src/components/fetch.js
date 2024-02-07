import { ref } from 'vue'

export function useFetch() {
	const data = ref(null)
	const error = ref(null)

	const fetchData = (url, options, callback) => {
		data.value = null
		error.value = null

		fetch(url, options)
			.then((res) => {
				if (res.status >= 500) throw new Error('Error was found')
				if (res.status >= 400) return res.text()
				return res.json()
			})
			.then((json) => {
				data.value = callback(json)
			})
			.catch((err) => { error.value = err; console.log(err) })
	}

	const post = (url, body, callback = () => { }) => {
		fetchData(url, {
			method: 'POST',
			body: JSON.stringify(body),
			headers: {
				'Content-type': 'application/json',
				'Authorization': 'Basic YWRtaW46Q29tcGxleHBhc3MjMTIz'
			}
		}, callback)
	}

	const postFile = (url, body, callback  = () => { }) => {
		fetchData(url, {
			method: 'POST',
			body,
		}, callback)
	}

	return { data, error, post, postFile }
}