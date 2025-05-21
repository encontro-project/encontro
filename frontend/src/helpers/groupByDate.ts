import type { ChatMessage, MessagesByDate } from '@/types'
import { getFormatedDate } from './getFormatedDate'

//Вынести в композабл
export const groupByDate = (messages: ChatMessage[]) => {
  const dateGroups = messages.reduce((acc: MessagesByDate, item) => {
    const date = new Date(item.timestamp)

    const dateKey = `${date.getDate()}-${date.getMonth()}-${date.getFullYear()}`

    if (!acc[dateKey]) {
      acc[dateKey] = {
        dateKey: getFormatedDate(item.timestamp),
        messages: [],
      }
    }
    acc[dateKey].messages.push(item)

    return acc
  }, {})
  console.log(dateGroups)
  return dateGroups
}
