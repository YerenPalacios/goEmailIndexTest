import { useFetch } from "./fetch"
import { ref } from 'vue'
export function messagesService() {
  const currentMessages = ref()
  const { data, error, post, loading, postFile } = useFetch()

  const get_from_name = xfrom => {
    const splitedText = xfrom?.split('<')
    if (splitedText && splitedText.length > 1) return splitedText[0]
    else return xfrom
  }

  const messageSerializer = (data) => {
    const hits = data?.hits?.hits || []

    const messages = hits.map(message => {
      const item = {
        id: message._id,
        subject: message._source.Subject || "(No subject provided)",
        content: message._source.content,
        from_email: message._source.From,
        from_name: get_from_name(message._source['X-From']),
        to_email: message._source.To,
        to_name: get_from_name(message._source['X-To']),
        date: message._source.Date
      }
      return item
    })
    currentMessages.value = messages
    return messages
  }

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


  const getMessages = (query, from = 0, to = 10) => {
    const payload = getQuery(query, from, to)
    post('http://localhost:4080/es/Messages/_search', payload, messageSerializer)
  }

  const sumMessages = (data) => {
    currentMessages.value = currentMessages.value.concat(messageSerializer(data))
    return data
  }

  const addMessages = (query, from = 0, to = 0) => {
    const payload = getQuery(query, from, to)
    post('http://localhost:4080/es/Messages/_search', payload, sumMessages)

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