export type PeerConnection = {
  displayName: string
  pc: RTCPeerConnection
  uuid: string
  ssVideoStream: MediaStream | null
  ssVideoTrack?: MediaStreamTrack
  ssAudioSender?: RTCRtpSender
  ssVideoSender?: RTCRtpSender
  ssAudioTrack?: MediaStreamTrack | null
}

//ss stand for share screen
