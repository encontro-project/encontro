const express = require('express')
const cors = require('cors')
const [
  serverSchema,
  userSchema,
  joinServerSchema,
  chatSchema,
  messageSchema,
] = require('./serverSchema')
const { Pool } = require('pg')
const bodyParser = require('body-parser')

const app = express()
app.use(cors())
app.use(bodyParser.json())
const PORT = 3001

require('dotenv').config()

console.log(process.env.DB_PORT)

const pool = new Pool({
  user: process.env.DB_USERNAME,
  host: 'localhost',
  database: 'crud_db',
  password: process.env.DB_PASSWORD,
  port: process.env.DB_PORT,
  max: 20,
  idleTimeoutMillis: 30000,
  connectionTimeoutMillis: 2000,
})

const Ajv = require('ajv')

const ajv = new Ajv()
const validateServer = ajv.compile(serverSchema)
const validateUser = ajv.compile(userSchema)
const validateJoinServer = ajv.compile(joinServerSchema)
const validateChat = ajv.compile(chatSchema)
const validateMessage = ajv.compile(messageSchema)

app.get('/get-server/:serverId', async (req, res) => {
  try {
    const serverId = req.params.serverId
    const server = await pool.query(`SELECT * FROM servers WHERE id = ${serverId}`)
    res.status(200).send(server.rows[0])
  } catch (error) {
    console.log(error)
  } finally {
  }
})

app.post('/add-server', async (req, res) => {
  try {
    const body = req.body
    const query = 'INSERT INTO servers (title, avatar_url) VALUES ($1, $2) RETURNING *'
    if (await validateServer(body)) {
      const values = [body.title, body.avatar_url]
      await pool.query(query, values)
      res.send('added server')
    } else {
      res.status(400).send('request data doesnt match schema')
    }
  } catch (error) {
    res.status(400).send('invalid request')
    console.log(error)
  } finally {
  }
})

app.get('/server/:serverId/chats', async (req, res) => {
  try {
    const serverId = req.params.serverId
    const query = `SELECT * FROM chats WHERE server_id = ${serverId}`
    const result = await pool.query(query)
    res.status(200).send(result.rows)
  } catch (error) {
    console.log(error)
    res.status(400).send('sdffsdg')
  }
})

app.post('/chat/add-message', async (req, res) => {
  try {
    const body = req.body
    const chatId = body.chat_id
    const content = body.content
    const senderId = body.sender_id
    const chatExists = (await pool.query(`SELECT * FROM chats WHERE id = ${chatId}`)).rowCount > 0
    const senderExists =
      (await pool.query(`SELECT * FROM users WHERE id = ${senderId}`)).rowCount > 0
    if (chatExists && senderExists && (await validateMessage(body))) {
      const values = [content, senderId, chatId]
      const query =
        'INSERT INTO messages (content, sender_id, chat_id) VALUES ($1, $2, $3) RETURNING *'
      await pool.query(query, values)
      res.status(200).send('added message')
    } else {
      res.status(400).send('chat/user doesnt exist or data doesnt match schema')
    }
  } catch (error) {
    console.log(error)
    res.status(400).send('invalid request')
  }
})

app.get('/chat/:chatId/get-messages', async (req, res) => {
  const chatId = req.params.chatId
  const result = await pool.query(
    `SELECT * FROM messages WHERE chat_id = ${chatId} ORDER BY sent_at`,
  )

  res.status(200).send(result.rows)
})

app.post('/server/add-chat', async (req, res) => {
  try {
    const body = req.body
    const serverId = body.server_id
    const title = body.title
    const serverExists =
      (await pool.query(`SELECT * FROM servers WHERE id = ${serverId}`)).rowCount > 0
    if (serverExists && (await validateChat(body))) {
      const values = [title, serverId]
      const query = 'INSERT INTO chats (title, server_id) VALUES ($1, $2) RETURNING *'
      await pool.query(query, values)
      res.status(200).send('added chat')
    } else {
      res.status(400).send('server doesnt exist or data doesnt match schema')
    }
  } catch (error) {
    res.status(400).send('invalid request')
  }
})

app.get('/user/:userId/servers', async (req, res) => {
  try {
    const userId = req.params.userId

    const userExists = (await pool.query(`SELECT * FROM users WHERE id = ${userId}`)).rowCount > 0

    if (userExists) {
      const result = await pool.query(
        `SELECT (server_id) FROM server_users WHERE user_id = ${userId} ORDER BY joined_at DESC`,
      )
      const serverIds = []
      for (i = 0; i < result.rowCount; i++) [serverIds.push(result.rows[i].server_id)]

      const servers = await pool.query(
        'SELECT * FROM servers WHERE id IN (SELECT unnest($1::int[]))',
        [serverIds],
      )

      res.status(200).send(servers.rows)
    } else {
      res.status(400).send('user doesnt exist')
    }
  } catch (error) {
    console.log(error)
    res.status(400).send('an error occured on request')
  }
})

app.post('/join-server/:serverId', async (req, res) => {
  try {
    const body = req.body
    const serverId = req.params.serverId
    if (validateJoinServer(body)) {
      const serverExists =
        (await pool.query(`SELECT * FROM servers WHERE id = ${serverId}`)).rowCount > 0
      const userExists =
        (await pool.query(`SELECT * FROM users WHERE id = ${body.user_id}`)).rowCount > 0
      if (userExists && serverExists) {
        const values = [body.user_id, serverId]
        const query =
          'INSERT INTO server_users (user_id, server_id, joined_at) VALUES ($1, $2, NOW()) ON CONFLICT DO NOTHING'
        const result = await pool.query(query, values)
        if (result.rowCount > 0) {
          res.status(200).send('server joined')
        } else {
          res.status(400).send('user already joined this server')
        }
      } else {
        res.status(400).send('server or user doesnt exist')
      }
    } else {
      res.status(400).send('request data doesnt match schema')
    }
  } catch (error) {
    console.log(error)
  }
})

app.get('/get-users', async (req, res) => {
  try {
    const users = await pool.query(`SELECT * FROM users`)
    res.send(users.rows)
  } catch (error) {}
})

app.post('/register-user', async (req, res) => {
  try {
    const body = req.body
    console.log()

    if (validateUser(body)) {
      const values = [body.username, body.avatar_url]
      const query = 'INSERT INTO users (username, avatar_url) VALUES ($1, $2) RETURNING *'
      await pool.query(query, values)
      res.status(200).end()
    } else {
      res.status(400).send('request data doesnt match schema')
    }
  } catch (error) {
    res.status(400).end()
    console.log(error)
  }
})

app.post('/leave-server/:serverId', async (req, res) => {})

app.listen(PORT, () => {
  console.log('эщкере')
})
