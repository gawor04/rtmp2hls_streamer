package session

import (
	"sync"
)

type SessionContainer struct {

	sesMap map[string] *Session
	mutex              sync.Mutex
	wg                 sync.WaitGroup
}

type SessionProperties struct {

	Session_id     string `json:"session_id"`
	Ingest_address string `json:"ingest_address"`
	Playback_url   string `json:"playback_url"`
}

/* global variable */
var Sessions = NewSessionContainer()

/* SessionContainer constructor */
func NewSessionContainer() *SessionContainer {

	sesCont := new(SessionContainer)
	sesCont.sesMap = make(map[string]*Session)
	sesCont.wg.Add(1)

	return sesCont
}

/* action if external command exits without error */
func ExitAction(identifier string) {

	Sessions.mutex.Lock()
	Sessions.sesMap[identifier].OnExit()
	Sessions.mutex.Unlock()

	defer Sessions.wg.Done()
}

/* need to restart ffmpeg after connection lost (external command exits without error) */
func SessionsObserver() {
	
	go func() {
		for {

			/* waits for ExitAction */
			Sessions.wg.Wait()
			Sessions.mutex.Lock()

			/* search which session was exited */
			for _, value := range Sessions.sesMap {

				if true == value.WasExited() {

					value.Restart(ExitAction)
				}
			}

			Sessions.mutex.Unlock()
			Sessions.wg.Add(1)
		}
	}()
}

/* clean up on program exit */
func SessionsClear() {

	Sessions.mutex.Lock()

	/* search which session is exited */
	for _, value := range Sessions.sesMap {

		value.Stop();
	}

	Sessions.mutex.Unlock()
}

/* creates new rtmp to ffpeg stream session */
func (sesCont *SessionContainer) NewSession(ip string, port string) SessionProperties {

	prop := SessionProperties{"", "", ""}

	if nil != sesCont {
		sesCont.mutex.Lock()

		ses := NewSession(ip, port, ExitAction)
		sesCont.sesMap[ses.GetPath()] = ses

		prop.Ingest_address = ses.GetRTMPurl()
		prop.Session_id = ses.GetUUID()
		prop.Playback_url = ses.GetHTTPPath()

		sesCont.mutex.Unlock()
	}

	return prop
}

/* stops external command and deletes session from container */
func (sesCont *SessionContainer) DeleteSession(uuid string) {

	if nil != sesCont {
		sesCont.mutex.Lock()

		for key, value := range sesCont.sesMap {

			if uuid == value.uuid {

				sesCont.sesMap[key].Stop()
				delete(sesCont.sesMap, key)
			}
		}

		sesCont.mutex.Unlock()
	}
}

/* returns possible video directories */
func (sesCont *SessionContainer) GetAvailableKeys() []string {

	sesCont.mutex.Lock()

	keys := make([]string, 0, len(sesCont.sesMap))

	for k := range sesCont.sesMap {

		keys = append(keys, k)
	}

	sesCont.mutex.Unlock()

	return keys
}
