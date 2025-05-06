export type PeerConnection = {
  displayName: string
  pc: RTCPeerConnection
  uuid: string
  ssVideoStream: MediaStream | null
  microphoneStream: MediaStream | null
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
