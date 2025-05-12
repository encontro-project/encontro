const PORT = 3000
const express = require('express')
const cors = require('cors')

const app = express()

app.use(cors())

app.get('/channel-info/:channelId', (req, res) => {
  const channelId = req.params.channelId

  const mockData = [
    {
      textChannels: [{ channelTitle: 'Гомики пишут', url: 'textRoom1' }],
      voiceChannels: [
        { channelTitle: 'Первая комната', url: 'room1' },
        { channelTitle: 'Вторая комната', url: 'room2' },
      ],
    },
    {
      textChannels: [{ channelTitle: 'Гомики пишут', url: 'textRoom1' }],
      voiceChannels: [
        { channelTitle: 'Большой хуй', url: 'room3' },
        { channelTitle: 'Вторая хуесссссссссс', url: 'room4' },
      ],
    },
  ]
  res.send(mockData[channelId - 1])
})

app.post('/new-participant', (req, res) => {
  try {
    const body = req.body

    channelParticipants[body.channelId][body.channelViewId][body.userUuid] = body

    res.status(200).end()
  } catch (error) {
    res.status(400).end()
  }
})

app.listen(PORT, () => {
  console.log(PORT)
})
