import { ref } from 'vue'

export function useFetch() {
	const data = ref(null)
	const error = ref(null)
	const loading = ref(false)

	const fetchData = (url, options, callback) => {
		data.value = null
		error.value = null
		loading.value = true

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
			.finally(()=>loading.value=false)
	}

	const post = (url, body, callback = () => { }) => {
		fetchData(url, {
			method: 'POST',
			body: JSON.stringify(body)
		}, callback)
	}

	const get = (url, callback) => {
		fetchData(url, {}, callback)
	}

	const postFile = (url, body, callback  = () => { }) => {
		fetchData(url, {
			method: 'POST',
			body,
		}, callback)
	}

	return { data, error, loading, post, get, postFile }
}