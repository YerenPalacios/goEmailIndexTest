import { useFetch } from "./fetch"
import moment from 'moment'

export function messagesService() {
  const { data, error, post } = useFetch()

  const get_from_name = xfrom => {
    const splitedText = xfrom.split('<')
    if (splitedText.length > 1) return splitedText[0]
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
    return messages
  }

  const getMessages = (query) => {
    const must = []

    if (!query) {
      must.push({ match_all: {} })
    } else {
      must.push({ query_string: { query } })
    }

    const payload = {
      query: {
        bool: {
          must
        }
      }
    }
    post('http://localhost:4080/es/Messages2/_search', payload, messageSerializer)
  }

  return { getMessages, data, error }
}