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

export type ChannelDescription = {
  channelTitle: string
  url: string
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
