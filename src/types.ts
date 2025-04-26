export type PeerConnection = {
  displayName: string
  pc: RTCPeerConnection
  dataChannel?: RTCDataChannel
  videoStream?: MediaStream
  audioStream?: MediaStream
  videoDecoder?: VideoDecoder
  audioDecoder?: AudioDecoder
}
