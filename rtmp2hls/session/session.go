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

const ffmpegCmdLinux = "ffmpeg"
const ffmpegCmdWindows = "ffmpeg.exe"
const ffmpegArg1 = "-listen -1 -rtmp_live live -i "
const ffmpegArg2 = " -c:v libx264 -crf 21 -preset veryfast -g 25 -sc_threshold 0 -c:a aac -b:a 128k -ac 2 -f hls -hls_flags delete_segments -hls_segment_size 500000 -hls_time 4 -hls_playlist_type event "
const rndStrLen = 20;

type Session struct {
	ip       string
	httpPort string
	uuid     string
	port     int
	path     string
	exited   bool
	cmdr     *utils.Commander
}

/* Session constructor */
func NewSession(ip string, http_port string, exitAction utils.OnExitFunc) *Session {

	session         := new(Session)
	session.ip       = ip
	session.httpPort = http_port

	/* generate uuid */
	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {

		log.Error("Cannot create uuid")
		return nil
	}

	/* remove new line from uuid string */
	strUiid := strings.TrimSuffix(string(uuid), "\n")
	log.Info(strUiid)
	session.uuid = strUiid

	/* find free port */
	port, err := freeport.GetFreePort()
	if err != nil {

		log.Error("Cannot find free port")
		return nil
	}
	session.port = port

	/* generate random path name */
	session.path = utils.RandomString(rndStrLen)

	session.exited = false;

	/* start ffmpeg command execution - concurrently */
	session.cmdr = session.prepareFFmpeg(exitAction)
	session.cmdr.Execute(session.path)

	log.Debug("New session, path: " + session.path)

	return session
}

/* stops external command execution */
func (session *Session) Stop() {

	if nil != session {

		session.cmdr.Stop()
		os.RemoveAll(session.path)
		os.Remove(session.path)
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

/* prepares directory and command for stream session */
func (session *Session) prepareFFmpeg(exitAction utils.OnExitFunc) *utils.Commander {

	if nil != session {

		os.Mkdir(session.path, 0755)

		address := session.GetRTMPurl()
		args := ffmpegArg1 + address + ffmpegArg2 + session.Getm3u8Path()

		ffmpegCmd := ""

		if runtime.GOOS == "windows" {

			ffmpegCmd += ffmpegCmdWindows
		} else {

			ffmpegCmd += ffmpegCmdLinux
		}

		return utils.NewCommander(ffmpegCmd, args, exitAction)
	}

	return nil
}

/*  returns rtmp server adress for this concrete stream */
func (session *Session) GetRTMPurl() string {

	url := ""

	if nil != session {

		url += "rtmp://" + session.ip + ":" + strconv.Itoa(session.port) + "/test/" + session.path
	}

	return url
}

/* returns hls output stream address  */
func (session *Session) GetHTTPPath() string {

	url := ""

	if nil != session {

		host := session.ip

		if len(session.httpPort) > 0 {

			host += ":" + session.httpPort
		}

		url += "http://" + host + "/play/" + session.path + ".m3u8"
	}

	return url
}

/* returns m3u8 file path (on filesystem) */
func (session *Session) Getm3u8Path() string {

	path := ""

	if nil != session {

		path += session.path + "/" + session.path + ".m3u8"
	}

	return path
}

/* returns session uuid */
func (session *Session) GetUUID() string {

	uuid := ""

	if nil != session {

		uuid += session.uuid
	}

	return uuid
}

/* returns path where output sream files are located */
func (session *Session) GetPath() string {

	path := ""

	if nil != session {
		path += session.path
	}

	return path
}

/* should be called if fmmpeg exits because of connection lost */
func (session *Session) OnExit() {

	session.exited = true;
}

/* checks if stream was excited */
func (session *Session) WasExited() bool {

	return session.exited
}
