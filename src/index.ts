type PeerConnection = {
  displayName: string
  pc: RTCPeerConnection
  dataChannel?: RTCDataChannel
  videoDecoder?: VideoDecoder
  audioDecoder?: AudioDecoder
}

let audioEncoder: AudioEncoder | null = null
let videoEncoder: VideoEncoder | null = null
let serverConnection: WebSocket | null = null
let peerConnections: Record<string, PeerConnection> = {}
let localUuid: string, localDisplayName: string, localStream: MediaStream

const WS_PORT = 8443
const PEER_CONNECTION_CFG: RTCConfiguration = {
  iceServers: [{ urls: 'stun:stun.l.google.com:19302' }, { urls: 'stun:stun1.l.google.com:19302' }],
}

const VIDEO_STREAM_WIDTH: number = 1280
const VIDEO_STREAM_HEIGHT: number = 720
const VIDEO_STREAM_FRAME_RATE: number = 30

const CONSTRAINTS: MediaStreamConstraints = {
  video: {
    width: { ideal: VIDEO_STREAM_WIDTH, max: VIDEO_STREAM_WIDTH },
    height: { ideal: VIDEO_STREAM_HEIGHT, max: VIDEO_STREAM_HEIGHT },
    frameRate: { ideal: VIDEO_STREAM_FRAME_RATE, max: VIDEO_STREAM_FRAME_RATE },
  },
  audio: true,
}

const VIDEO_ENCODER_CFG: VideoEncoderConfig = {
  codec: 'vp8',
  width: VIDEO_STREAM_WIDTH,
  height: VIDEO_STREAM_HEIGHT,
  framerate: VIDEO_STREAM_FRAME_RATE,
  bitrate: 1_000_000,
  latencyMode: 'realtime',
}

const VIDEO_DECODER_CFG: VideoDecoderConfig = {
  codec: VIDEO_ENCODER_CFG.codec,
  codedWidth: VIDEO_ENCODER_CFG.width,
  codedHeight: VIDEO_ENCODER_CFG.height,
  hardwareAcceleration: 'no-preference',
}

const AUDIO_ENCODER_CFG: AudioEncoderConfig = {
  codec: 'opus',
  numberOfChannels: 1,
  sampleRate: 48000,
  bitrate: 128000,
}

const AUDIO_DECODER_CFG: AudioDecoderConfig = {
  codec: AUDIO_ENCODER_CFG.codec,
  numberOfChannels: AUDIO_ENCODER_CFG.numberOfChannels,
  sampleRate: AUDIO_ENCODER_CFG.sampleRate,
}

async function initializeCodecs(stream: MediaStream): Promise<void> {
  if (!('VideoEncoder' in window)) return

  const videoTrack = stream.getVideoTracks()[0]
  const audioTrack = stream.getAudioTracks()[0]

  const videoProcessor = new window.MediaStreamTrackProcessor({ track: videoTrack })
  const videoReader = videoProcessor.readable.getReader()

  videoEncoder = new VideoEncoder({
    output: (encodedChunk: EncodedVideoChunk) => {
      const data = new ArrayBuffer(encodedChunk.byteLength)
      const view = new Uint8Array(data)
      encodedChunk.copyTo(view)

      Object.values(peerConnections).forEach((peer) => {
        if (peer.dataChannel?.readyState === 'open') {
          peer.dataChannel.send(
            JSON.stringify({
              type: encodedChunk.type,
              timestamp: encodedChunk.timestamp,
              duration: encodedChunk.duration,
            }),
          )
          peer.dataChannel.send(data)
        }
      })
    },
    error: (e) => console.error(e),
  })

  videoEncoder.configure(VIDEO_ENCODER_CFG)

  const processVideoFrames = async () => {
    try {
      while (true) {
        const { value: frame, done } = await videoReader.read()
        if (done) break
        if (videoEncoder?.state === 'configured') {
          videoEncoder?.encode(frame!)
        }
        frame!.close()
      }
    } catch (e) {
      console.error('Error processing frames:', e)
    }
  }

  processVideoFrames()

  if (audioTrack && 'AudioEncoder' in window) {
    console.log(audioTrack.getSettings())
    const audioProcessor = new MediaStreamTrackProcessor({ track: audioTrack })
    const audioReader = audioProcessor.readable.getReader()

    audioEncoder = new AudioEncoder({
      output: (encodedChunk: EncodedAudioChunk) => {
        const data = new ArrayBuffer(encodedChunk.byteLength)
        const view = new Uint8Array(data)
        encodedChunk.copyTo(view)

        Object.values(peerConnections).forEach((peer) => {
          if (peer.dataChannel?.readyState === 'open') {
            peer.dataChannel.send(
              JSON.stringify({
                type: 'audio',
                timestamp: encodedChunk.timestamp,
                duration: encodedChunk.duration,
              }),
            )
            peer.dataChannel.send(data)
          }
        })
      },
      error: (e) => console.error('Audio Encoder Error:', e),
    })

    audioEncoder.configure(AUDIO_ENCODER_CFG)

    const processAudioFrames = async () => {
      try {
        while (true) {
          const { value: frame, done } = await audioReader.read()
          if (done) break
          if (audioEncoder?.state === 'configured' && frame) {
            try {
              audioEncoder.encode(frame)
            } catch (e) {
              console.error('Error encoding audio frame:', e)
            }
          }
          frame.close()
        }
      } catch (e) {
        console.error('Error processing audio frames:', e)
      }
    }

    processAudioFrames()
  }
}

async function start(): Promise<void> {
  localUuid = createUUID()
  const urlParams = new URLSearchParams(window.location.search)
  localDisplayName = urlParams.get('displayName') || localUuid

  const labelContainer = document.getElementById('localVideoContainer')
  if (labelContainer) labelContainer.appendChild(makeLabel(localDisplayName))

  try {
    localStream = await getMediaStream()

    const localVideo = document.getElementById('localVideo') as HTMLVideoElement
    localVideo.srcObject = localStream

    await initializeCodecs(localStream)

    serverConnection = new WebSocket(`wss://${location.hostname}:${WS_PORT}/ws/`)
    serverConnection.onmessage = (msg) => gotMessageFromServer(msg)
    serverConnection.onopen = () => {
      serverConnection!.send(
        JSON.stringify({
          displayName: localDisplayName,
          uuid: localUuid,
          dest: 'all',
        }),
      )
    }
  } catch (error) {
    errorHandler(error)
  }
}

async function getMediaStream(): Promise<MediaStream> {
  if (navigator.mediaDevices.getDisplayMedia) {
    return await navigator.mediaDevices.getDisplayMedia(CONSTRAINTS)
  } else if (navigator.mediaDevices.getUserMedia) {
    return await navigator.mediaDevices.getUserMedia(CONSTRAINTS)
  } else {
    alert('Your browser does not support getDisplayMedia or getUserMedia API')
    throw new Error('Unsupported media API')
  }
}

function setUpPeer(peerUuid: string, displayName: string, initCall = false): void {
  const peerConnection: PeerConnection = {
    displayName,
    pc: new RTCPeerConnection(PEER_CONNECTION_CFG),
  }

  const dataChannel = peerConnection.pc.createDataChannel('media-channel')
  peerConnection.dataChannel = dataChannel

  let metadata: any

  let videoWritable: WritableStreamDefaultWriter<VideoFrame> | undefined
  let videoStreamGenerator: MediaStreamTrackGenerator<VideoFrame> | undefined

  let audioStreamGenerator: MediaStreamTrackGenerator<AudioData> | undefined
  let audioWritable: WritableStreamDefaultWriter<AudioData> | undefined

  dataChannel.onopen = async () => {
    if (!('VideoDecoder' in window)) return

    peerConnection.videoDecoder = new VideoDecoder({
      output: (frame: VideoFrame) => {
        if (!videoStreamGenerator) {
          videoStreamGenerator = new MediaStreamTrackGenerator({
            kind: 'video',
          })
          videoWritable = videoStreamGenerator.writable.getWriter()

          const stream = new MediaStream([videoStreamGenerator])
          const vidElement = document.querySelector(
            `#remoteVideo_${peerUuid} video`,
          ) as HTMLVideoElement
          if (vidElement) {
            vidElement.srcObject = stream
            vidElement.play().catch((e) => console.error('Error playing video:', e))
          }
        }
        videoWritable!.write(frame)
      },
      error: (e) => console.error(e),
    })

    peerConnection.videoDecoder.configure(VIDEO_DECODER_CFG)

    if ('AudioDecoder' in window) {
      peerConnection.audioDecoder = new AudioDecoder({
        output: async (audioData: AudioData) => {
          if (!audioStreamGenerator) {
            audioStreamGenerator = new MediaStreamTrackGenerator({
              kind: 'audio',
            })
            audioWritable = audioStreamGenerator.writable.getWriter()

            const stream = new MediaStream([audioStreamGenerator])
            const audioElement = document.querySelector(
              `#remoteAudio_${peerUuid}`,
            ) as HTMLAudioElement
            if (audioElement) {
              audioElement.srcObject = stream
              audioElement.play().catch((e) => console.error('Error playing audio:', e))
            }
          }
          audioWritable!.write(audioData)
        },
        error: (e) => console.error('Audio decoder error:', e),
      })

      peerConnection.audioDecoder.configure(AUDIO_DECODER_CFG)
    }
  }

  dataChannel.onmessage = async (event: MessageEvent) => {
    if (typeof event.data === 'string') {
      metadata = JSON.parse(event.data)
      console.log(metadata)
    } else if (metadata.type === 'video' || metadata.type == 'delta') {
      const chunk = new EncodedVideoChunk({
        type: 'key',
        timestamp: metadata.timestamp,
        duration: metadata.duration,
        data: event.data,
      })

      if (peerConnection.videoDecoder?.state === 'configured') {
        peerConnection.videoDecoder.decode(chunk)
      }
    } else if (metadata.type === 'audio') {
      const chunk = new EncodedAudioChunk({
        type: 'key',
        timestamp: metadata.timestamp,
        duration: metadata.duration,
        data: event.data,
      })

      if (peerConnection.audioDecoder?.state === 'configured') {
        peerConnection.audioDecoder.decode(chunk)
      }
    }
  }

  peerConnection.pc.ontrack = (event) => {
    if (event.track.kind === 'video') {
      gotRemoteStream(event, peerUuid)
    }
  }
  peerConnection.pc.onicecandidate = (event) => gotIceCandidate(event, peerUuid)
  peerConnection.pc.oniceconnectionstatechange = () => checkPeerDisconnect(peerUuid)
  localStream.getTracks().forEach((track) => {
    peerConnection.pc.addTrack(track, localStream)
  })
  peerConnections[peerUuid] = peerConnection

  if (initCall) {
    peerConnection.pc
      .createOffer()
      .then((description) => createdDescription(description, peerUuid))
      .catch(errorHandler)
  }
}

function gotMessageFromServer(message: MessageEvent): void {
  const signal = JSON.parse(message.data)
  const peerUuid = signal.uuid

  if (peerUuid === localUuid || (signal.dest !== localUuid && signal.dest !== 'all')) return

  if (signal.type === 'peer-disconnect') {
    handlePeerDisconnect(peerUuid)
    return
  }

  if (signal.displayName && signal.dest === 'all') {
    setUpPeer(peerUuid, signal.displayName)
    serverConnection!.send(
      JSON.stringify({
        displayName: localDisplayName,
        uuid: localUuid,
        dest: peerUuid,
      }),
    )
  } else if (signal.displayName && signal.dest === localUuid) {
    setUpPeer(peerUuid, signal.displayName, true)
  } else if (signal.sdp) {
    handleSdpSignal(signal, peerUuid)
  } else if (signal.ice) {
    handleIceCandidate(signal, peerUuid)
  }
}

function handleSdpSignal(signal: any, peerUuid: string): void {
  peerConnections[peerUuid].pc
    .setRemoteDescription(new RTCSessionDescription(signal.sdp))
    .then(() => {
      if (signal.sdp.type === 'offer') {
        peerConnections[peerUuid].pc
          .createAnswer()
          .then((description) => createdDescription(description, peerUuid))
          .catch(errorHandler)
      }
    })
    .catch(errorHandler)
}

function handleIceCandidate(signal: any, peerUuid: string): void {
  peerConnections[peerUuid].pc.addIceCandidate(new RTCIceCandidate(signal.ice)).catch(errorHandler)
}

function handlePeerDisconnect(peerUuid: string): void {
  console.log(`Peer ${peerUuid} disconnected`)
  delete peerConnections[peerUuid]
  const vidElement = document.getElementById(`remoteVideo_${peerUuid}`)
  if (vidElement) document.getElementById('videos')?.removeChild(vidElement)
  updateLayout()
}

function gotRemoteStream(event: RTCTrackEvent, peerUuid: string): void {
  const vidElement = document.createElement('video')
  vidElement.setAttribute('autoplay', '')
  vidElement.setAttribute('muted', '')
  vidElement.srcObject = event.streams[0]

  const vidContainer = document.createElement('div')
  vidContainer.id = `remoteVideo_${peerUuid}`
  vidContainer.className = 'videoContainer'
  vidContainer.appendChild(vidElement)
  vidContainer.appendChild(makeLabel(peerConnections[peerUuid].displayName))

  document.getElementById('videos')?.appendChild(vidContainer)
  updateLayout()
}

function checkPeerDisconnect(peerUuid: string): void {
  const state = peerConnections[peerUuid].pc.iceConnectionState
  if (['failed', 'closed', 'disconnected'].includes(state)) {
    handlePeerDisconnect(peerUuid)
  }
}

function gotIceCandidate(event: RTCPeerConnectionIceEvent, peerUuid: string): void {
  if (event.candidate) {
    serverConnection!.send(
      JSON.stringify({
        ice: event.candidate,
        uuid: localUuid,
        dest: peerUuid,
      }),
    )
  }
}

function createdDescription(description: RTCSessionDescriptionInit, peerUuid: string): void {
  peerConnections[peerUuid].pc
    .setLocalDescription(description)
    .then(() => {
      serverConnection!.send(
        JSON.stringify({
          sdp: peerConnections[peerUuid].pc.localDescription,
          uuid: localUuid,
          dest: peerUuid,
        }),
      )
    })
    .catch(errorHandler)
}

function updateLayout(): void {
  const numVideos = Object.keys(peerConnections).length + 1
  const rowHeight = numVideos > 1 && numVideos <= 4 ? '48vh' : numVideos > 4 ? '32vh' : '98vh'
  const colWidth = numVideos > 1 && numVideos <= 4 ? '48vw' : numVideos > 4 ? '32vw' : '98vw'
  document.documentElement.style.setProperty('--rowHeight', rowHeight)
  document.documentElement.style.setProperty('--colWidth', colWidth)
}

function makeLabel(label: string): HTMLElement {
  const vidLabel = document.createElement('div')
  vidLabel.textContent = label
  vidLabel.className = 'videoLabel'
  return vidLabel
}

function createUUID(): string {
  return [1e7, -1e3, -4e3, -8e3, -1e11]
    .join('')
    .replace(/[018]/g, (c: string) =>
      (
        parseInt(c, 16) ^
        (crypto.getRandomValues(new Uint8Array(1))[0] & (15 >> (parseInt(c, 16) / 4)))
      ).toString(16),
    )
}

function errorHandler(error: any): void {
  console.error(error)
}

window.addEventListener('beforeunload', () => {
  if (videoEncoder?.state !== 'closed') {
    videoEncoder?.close()
  }
  if (audioEncoder?.state !== 'closed') {
    audioEncoder?.close()
  }

  Object.values(peerConnections).forEach((peer) => {
    if (peer.videoDecoder?.state !== 'closed') {
      peer.videoDecoder?.close()
    }
    if (peer.audioDecoder?.state !== 'closed') {
      peer.audioDecoder?.close()
    }
    if (peer.dataChannel) {
      peer.dataChannel.close()
    }
  })

  serverConnection?.send(
    JSON.stringify({
      type: 'peer-disconnect',
      uuid: localUuid,
    }),
  )
})
