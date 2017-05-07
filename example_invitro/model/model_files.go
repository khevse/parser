package model

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"net/http"
	"sync"

	"github.com/ncw/swift"
)

const FILE_NAME_PREFIX = "invitro"

type ModelFiles struct {
	conn      *swift.Connection
	container string
	once      sync.Once
}

func NewModelFiles(conn *swift.Connection, container string) *ModelFiles {
	return &ModelFiles{
		conn:      conn,
		container: container,
	}
}

func (m *ModelFiles) Add(data []byte) (string, error) {
	if err := m.internalInit(); err != nil {
		return "", err
	}

	name := FILE_NAME_PREFIX + md5Sum(data)
	contentType := http.DetectContentType(data)

	if err := m.conn.ObjectPutBytes(m.container, name, data, contentType); err != nil {
		log.Println(err)
		return "", err
	}

	return name, nil
}

func (m *ModelFiles) Get(name string) ([]byte, error) {
	if err := m.internalInit(); err != nil {
		return nil, err
	}

	if body, err := m.conn.ObjectGetBytes(m.container, name); err != nil {
		log.Println(err)
		return nil, err
	} else {
		return body, nil
	}
}

func (m *ModelFiles) internalInit() (retval error) {
	m.once.Do(func() {
		retval = m.removeAll()
		if retval != nil {
			log.Println(retval)
		}
	})

	return
}

func (m *ModelFiles) removeAll() error {

	opt := swift.ObjectsOpts{
		Prefix: FILE_NAME_PREFIX,
	}

	if list, err := m.conn.ObjectNames(m.container, &opt); err != nil {
		log.Println(err)
		return err
	} else if len(list) > 0 {
		if _, err := m.conn.BulkDelete(m.container, list); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func md5Sum(data []byte) string {

	shasum := make([]byte, 32)
	tmpbuf := md5.Sum(data)
	hex.Encode(shasum, tmpbuf[:])

	return string(shasum)
}
