package session

import (
	"../../utils"
	"github.com/phayes/freeport"
	log "github.com/cihub/seelog"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const ffmpeg_cmd_linux = "ffmpeg"
const ffmpeg_cmd_windows = "ffmpeg.exe"
const ffmpeg_arg1 = "-listen -1 -rtmp_live live -i "
const ffmpeg_arg2 = " -c:v libx264 -crf 21 -preset veryfast -g 25 -sc_threshold 0 -c:a aac -b:a 128k -ac 2 -f hls -hls_flags delete_segments -hls_segment_size 500000 -hls_time 4 -hls_playlist_type event "


type Session struct {
	ip string
	uuid string
	port int
	path string
	exited bool
	cmdr *utils.Commander
}

func NewSession(ip string, exitAction utils.OnExitFunc) *Session {

	const rndStrLen = 20;

	session := new(Session)
	session.ip = ip

	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Error("Cannot create uuid")
		return nil
	}

	/* remove new line */
    strUiid := strings.TrimSuffix(string(uuid), "\n")
	log.Info(strUiid )
	session.uuid = strUiid

	port, err := freeport.GetFreePort()
	if err != nil {
		log.Error("Cannot find free port")
		return nil
	}
	session.port = port

	session.path = utils.RandomString(rndStrLen)

	session.exited = true;
	session.cmdr = session.prepareFFmpeg(exitAction)
	session.cmdr.Execute(session.path)

	log.Debug("New session, path: " + session.path)

	return session
}

/* 	stops external command execution */
func (session *Session) Stop() {

	if nil != session {

		session.cmdr.Stop()
	}
}

/* restarts external command */
func (session *Session) Restart(exitAction utils.OnExitFunc) {

	if nil != session {

		log.Debug("Restarting session, path: " + session.path)
		os.RemoveAll(session.path)
		os.Mkdir(session.path, 0755)
		session.cmdr = session.prepareFFmpeg(exitAction)
		session.cmdr.Execute(session.path)
		session.exited = false
	}
}

func (session *Session) prepareFFmpeg(exitAction utils.OnExitFunc) *utils.Commander {

	if nil != session {

		os.Mkdir(session.path, 0755)

		address := session.GetRTMPurl()
		args := ffmpeg_arg1 + address + ffmpeg_arg2 + session.Getm3u8Path()

		ffmpeg_cmd := ""

		if runtime.GOOS == "windows" {
			ffmpeg_cmd += ffmpeg_cmd_windows
		} else {
			ffmpeg_cmd += ffmpeg_cmd_linux
		}

		return utils.NewCommander(ffmpeg_cmd, args, exitAction)
	}

	return nil
}

func (session *Session) GetRTMPurl() string {

	url := ""

	if nil != session {

		url += "rtmp://" + session.ip + ":" + strconv.Itoa(session.port) + "/test/" + session.path
	}

	return url
}

func (session *Session) GetHTTPPath() string {

	url := ""

	if nil != session {

		url += "http://" + session.ip + "/play/" + session.path + ".m3u8"
	}

	return url
}

func (session *Session) Getm3u8Path() string {

	path := ""

	if nil != session {

		path += session.path + "/" + session.path + ".m3u8"
	}

	return path
}

func (session *Session) GetUUID() string {

	uuid := ""

	if nil != session {

		uuid += session.uuid
	}

	return uuid
}

func (session *Session) GetPath() string {

	path := ""

	if nil != session {

		path += session.path
	}

	return path
}

func (session *Session) OnExit() {

	session.exited = true;
}

func (session *Session) WasExited() bool {

	return session.exited
}