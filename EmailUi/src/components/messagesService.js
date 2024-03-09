import { useFetch } from "./fetch"
import { ref } from 'vue'
export function messagesService() {
  const currentMessages = ref()
  const { data, error, post, get, loading, postFile } = useFetch()

  const getQuery = (query, from, to) => {
    const must = []

    if (!query) {
      must.push({ match_all: {} })
    } else {
      must.push({ query_string: { query } })
    }

    const payload = {
      query: { bool: { must } }, from, to, sort: ["-Date"]
    }
    return payload
  }


  const getMessages = (query, callback, from = 0, to = 2, ) => {

    const searchParams = new URLSearchParams({
      from, to, search: query,
    })
    get(
        `http://localhost:8000/messages?${searchParams}`,
        callback
    )
  }

  const sumMessages = (func) => {
    return function(){
      const newMessages = func()
      if (newMessages) {
        currentMessages.value = currentMessages.value.concat(callback(newMessages))
        return newMessages
      }
    }

  }

  const addMessages = (query, from = 0, to = 0, callback) => {
    const payload = getQuery(query, from, to)
    post('http://localhost:4080/es/Messages/_search', payload, sumMessages(callback))

  }

  const uploadFile = (file, callback) => {
    const form = new FormData()
    form.append("file", file)
    postFile('http://localhost:8000/import_file', form, callback)

  }

  return {
    getMessages,
    addMessages,
    uploadFile,
    currentMessages,
    data,
    loading,
    error
  }
}