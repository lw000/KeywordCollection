package models

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type WSCMDReader interface {
	// Encode ...
	Encode() ([]byte, error)
	// Decode ...
	Decode(data []byte) error
}

// WSCMD ...
type WSCMD struct {
	MainID int    `json:"main_id,omitempty"`
	SubID  int    `json:"sub_id,omitempty"`
	Data   string `json:"data,omitempty"`
}

type WSAckQuery struct {
	Code int    `json:"code,omitempty"`
	Id   int    `json:"id,omitempty"`
	Rank int    `json:"rank,omitempty"`
	Data string `json:"data,omitempty"`
}

type WSAckOpen struct {
	Code int    `json:"code,omitempty"`
	Data string `json:"data,omitempty"`
}

func (w WSCMD) String() string {
	return fmt.Sprintf("MainID:%d SubID:%d Data:%s", w.MainID, w.SubID, w.Data)
}

// Encode ...
func (w *WSCMD) Encode(data string) ([]byte, error) {
	w.Data = data
	buf, err := json.Marshal(w)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return buf, nil
}

// Encode ...
func (w *WSCMD) EncodeCmd(reader WSCMDReader) ([]byte, error) {
	if reader != nil {
		dbuf, err := reader.Encode()
		if err != nil {
			log.Error(err)
			return nil, err
		}
		w.Data = string(dbuf)
	}
	buf, err := json.Marshal(w)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return buf, nil
}

// Decode ...
func (w *WSCMD) Decode(data []byte) error {
	if err := json.Unmarshal(data, w); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// Encode ...
func (w *WSAckQuery) Encode() ([]byte, error) {
	buf, err := json.Marshal(w)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return buf, nil
}

// Decode ...
func (w *WSAckQuery) Decode(data []byte) error {
	if err := json.Unmarshal(data, w); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
