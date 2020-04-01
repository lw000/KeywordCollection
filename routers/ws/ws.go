package ws

import (
	"KeywordCollection/global"
	"KeywordCollection/models"
	"github.com/lw000/gocommon/utils"
	"net/url"

	"github.com/olahol/melody"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// RegisterService ...
func RegisterService(engine *gin.Engine) {

	m := melody.New()

	engine.GET("ws", func(c *gin.Context) {
		err := m.HandleRequest(c.Writer, c.Request)
		if err != nil {
			log.Error(err)
		}
	})

	m.HandleConnect(func(s *melody.Session) {
		values, err := url.ParseQuery(s.Request.URL.RawQuery)
		if err != nil {
			log.Error(err)
			err = s.CloseWithMsg([]byte("Illegal access"))
			if err != nil {
				log.Error(err)
				return
			}
			return
		}

		clientID := values.Get("clientID")
		if clientID == "" {
			clientID = tyutils.UUID()
		}

		s.Set("clientID", clientID)
		global.QueryClients.Store(clientID, s)

		log.WithFields(log.Fields{"clientID": clientID, "RemoteAddr": s.Request.RemoteAddr}).Info("join")

		cmd := &models.WSCMD{MainID: 1, SubID: 1}
		data, err := cmd.Encode(clientID)
		if err != nil {
			log.Error(err)
			return
		}
		if err = s.Write(data); err != nil {
			log.Error(err)
		}
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		if err := m.Broadcast(msg); err != nil {
			log.Error(err)
			return
		}
	})

	m.HandleDisconnect(func(s *melody.Session) {
		v, exists := s.Get("clientID")
		if exists {
			log.WithFields(log.Fields{"RemoteAddr": s.Request.RemoteAddr}).Info("leave")
			global.QueryClients.Delete(v.(string))
		}
	})
}
