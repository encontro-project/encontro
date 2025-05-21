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
}

//ss stand for share screen

export {}

export type Server = {
  title: string
  voice_channels: VoiceChannel[]
  text_chats: TextChat[]
  accociated_users: MessageUser[]
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

type MessageUser = {
  id: number
  username: string
  avatar_url: string
}

export type ChatInfo = {
  title: string
  messages: ChatMessage[]
}
type ChatMessage = {
  text: string
  timestamp: number
  sender: number
}
