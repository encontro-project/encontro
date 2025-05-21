export const getFormatedDate = (timestamp: string) => {
  const monthsDict = {
    0: 'Янв',
    1: 'Фев',
    2: 'Мар',
    3: 'Апр',
    4: 'Май',
    5: 'Июн',
    6: 'Июл',
    7: 'Авг',
    8: 'Сен',
    9: 'Окт',
    10: 'Ноя',
    11: 'Дек',
  }
  const date = new Date(timestamp)
  const day = date.getDate()
  const month = date.getMonth()
  const year = date.getFullYear()

  return `${day} ${(monthsDict as any)[month]}, ${year}`
}
