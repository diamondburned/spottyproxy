package api

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"

	"github.com/librespot-org/librespot-golang/librespot"
	librespotcore "github.com/librespot-org/librespot-golang/librespot/core"
	"libdb.so/spottyproxy/server/internal/crudkv"
	"libdb.so/spottyproxy/server/internal/syncg"
)

type SessionStore struct {
	store    crudkv.Store[sessionData]
	sessions syncg.Map[string, *Session]
	initing  syncg.SingleflightGroup[*Session]
}

// NewSessionStore returns a new SessionStore.
func NewSessionStore(store crudkv.BasicStore) *SessionStore {
	return &SessionStore{
		store: crudkv.Wrap[sessionData](store, "session_store"),
	}
}

// Acquire acquires a session for the given user. If the session does not exist,
// it is created. If the user has never logged in, a NotFound error is returned.
func (s *SessionStore) Acquire(token string) (*Session, error) {
	session, ok := s.sessions.Load(token)
	if ok {
		return session, nil
	}

	session, err := s.initing.Do(token, func() (*Session, error) {
		data, err := s.store.Get(token)
		if err != nil {
			return nil, err
		}

		coreSession, err := librespot.LoginSaved(data.Username, data.AuthData, data.DeviceName)
		if err != nil {
			return nil, err
		}

		session := &Session{
			Session: coreSession,
			token:   token,
		}

		s.sessions.Store(token, session)
		return session, nil
	})
	if err != nil {
		return nil, err
	}

	session.mutex.Lock()
	return session, nil
}

// LoginAndAcquire logs in the given user and acquires a new session for them.
func (s *SessionStore) LoginAndAcquire(username, password, device string) (*Session, error) {
	coreSession, err := librespot.Login(username, password, device)
	if err != nil {
		return nil, err
	}

	var token string
	err = s.store.Tx(true, func(txn crudkv.Transaction[sessionData]) error {
		for {
			token = generateToken()

			_, err := txn.Get(token)
			if err == nil {
				continue
			}

			if !errors.Is(err, crudkv.ErrNotFound) {
				return err
			}

			return txn.Set(token, sessionData{
				Username:   username,
				AuthData:   coreSession.ReusableAuthBlob(),
				DeviceName: device,
			})
		}
	})
	if err != nil {
		return nil, err
	}

	return s.Acquire(token)
}

func generateToken() string {
	randBytes := make([]byte, 16)

	_, err := rand.Read(randBytes)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(randBytes)
}

type sessionData struct {
	Username   string
	DeviceName string
	AuthData   []byte
}

// Session is a shared librespot session.
type Session struct {
	*librespotcore.Session
	mutex sync.Mutex
	token string
}

// Release releases the session. The session is no longer usable after this
// method is called.
func (s *Session) Release() {
	s.mutex.Unlock()
}
