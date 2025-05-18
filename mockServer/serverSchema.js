const serverSchema = {
  type: 'object',
  properties: {
    title: { type: 'string' },
    avatar_url: { type: 'string' },
    /* server_settings: { type: 'object' }, //TODO Придумать структуру настроек для сервера */
  },
  required: ['title'],
}

const userSchema = {
  type: 'object',
  properties: {
    username: { type: 'string' },
    avatar_url: { type: 'string' },
  },
  required: ['username', 'avatar_url'],
}

const joinServerSchema = {
  type: 'object',
  properties: {
    user_id: { type: 'integer' },
  },
  required: ['user_id'],
}

const chatSchema = {
  type: 'object',
  properties: {
    title: { type: 'string' },
    server_id: { type: 'integer' },
  },
  required: ['title', 'server_id'],
}
const messageSchema = {
  type: 'object',
  properties: {
    content: { type: 'string' },
    chat_id: { type: 'integer' },
    sender_id: { type: 'integer' },
  },
  required: ['content', 'chat_id', 'sender_id'],
}

module.exports = [serverSchema, userSchema, joinServerSchema, chatSchema, messageSchema]
