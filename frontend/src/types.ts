export type PeerConnection = {
  displayName: string
  pc: RTCPeerConnection
  uuid: string
  ssVideoStream: MediaStream | null
  microphoneStream: MediaStream | null
  taskQueue: ((...args: any) => Promise<void>)[]
  trackMetadata?: 'microphone' | 'screen-share' | 'webcam'
  shareScreenStreamid?: string
  microphoneStreamId?: string
  ssVideoTrack?: MediaStreamTrack
  ssAudioSender?: RTCRtpSender
  ssVideoSender?: RTCRtpSender
  microphoneSender?: RTCRtpSender
  ssAudioTrack?: MediaStreamTrack | null
  userVolume: number
  isMuted: boolean
}

//ss stand for share screen

export {}

export type Server = {
  title: string
  voice_channels: VoiceChannel[]
  text_chats: TextChat[]
  accociated_users: ServerUser[]
}

type VoiceChannel = {
  id: number
  title: string
}

type TextChat = {
  id: number
  title: string
  messages: Message[]
}

type User = {}

type Message = {
  id: string
  content: string
  sender_id: string
  created_at: string
  updated_at: string
}

type ServerUser = {
  id: number
  username: string
  avatar_url: string
}

export type MessagesByDate = Record<string, { dateKey: string; messages: ChatMessage[] }>

export type ChatInfo = {
  title: string
  messages: ChatMessage[]
}
export type ChatMessage = {
  text: string
  timestamp: string
  sender: string
}
