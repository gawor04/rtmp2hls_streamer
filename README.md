# rtmp2hls_streamer

rtmp2hls_streamer is GO powered, ffmpeg based rtmp to hls live-streaming application. Stream sessions can be controlled (add/delete) via REST API.

### Installation

rtmp2hls_streamer requires ffmpeg installation on your machine.

External go the dependencies:
```sh
$ go get -u github.com/cihub/seelog
$ go get github.com/phayes/freeport
```

Download rtmp2hls_streamer
```sh
$ git clone https://github.com/gawor04/rtmp2hls_streamer
```

run rtmp2hls_streamer
```sh
$ go run main.go
```

### Create new streaming session
```sh
$ curl -o out.json http://127.0.0.1:8080/sessions
$cat out.json 
{"session_id":"07b74323-c1c3-4cba-85ea-2ad412401eec","ingest_address":"rtmp://127.0.0.1:33005/test/5a84jjJkwzDkh9h2fhfU","playback_url":"http://127.0.0.1/play/5a84jjJkwzDkh9h2fhfU.m3u8"}
```

rtmp://127.0.0.1:33005/test/5a84jjJkwzDkh9h2fhfU - rtmp stream input
http://127.0.0.1/play/5a84jjJkwzDkh9h2fhfU.m3u8  - HLS stream output
07b74323-c1c3-4cba-85ea-2ad412401eec             - session id

### Testing
1. In OBS set rtmp://127.0.0.1:33005/test as server and 5a84jjJkwzDkh9h2fhfU as a key
2. Start streaming in OBS
3. Play video using ffplay
```sh
$ ffplay http://127.0.0.1:8080/play/5a84jjJkwzDkh9h2fhfU.m3u8
```

### Delete session (using session id)
```sh
$ curl -X DELETE http://127.0.0.1:8080/sessions/07b74323-c1c3-4cba-85ea-2ad412401eec
