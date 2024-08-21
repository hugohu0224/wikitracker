package internal

//
//import (
//	"context"
//	"github.com/gin-gonic/gin"
//	"github.com/gorilla/websocket"
//	"go.uber.org/zap"
//	"net/http"
//	"sync"
//	"wikitracker/pkg/models"
//)
//
//var upgrader = websocket.Upgrader{
//	ReadBufferSize:  64,
//	WriteBufferSize: 64,
//	CheckOrigin: func(r *http.Request) bool {
//		return true
//	},
//}
//
//type Client struct {
//	ID           string
//	Conn         *websocket.Conn
//	Send         chan *models.GameMsg
//	mu            sync.Mutex
//}
//
//func WebsocketEndpoint(c *gin.Context) {
//	// http
//	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
//	if err != nil {
//		zap.S().Error("websocket upgrade failed", zap.Error(err))
//		zap.S().Debug("request details", zap.Any("headers", c.Request.Header))
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade connection"})
//		return
//	}
//	zap.S().Debugf("websocket connected to server %v\n", conn.RemoteAddr())
//	serveWs()
//
//func serveWs() {
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	var wg sync.WaitGroup
//	wg.Add(2)
//
//	errChan := make(chan error, 2)
//
//	go func() {
//		defer wg.Done()
//		if err := client.ReadPump(ctx); err != nil {
//			errChan <- err
//		}
//	}()
//
//	go func() {
//		defer wg.Done()
//		if err := client.WritePump(ctx); err != nil {
//			errChan <- err
//		}
//	}()
//
//	go func() {
//		wg.Wait()
//		close(errChan)
//	}()
//
//	for err := range errChan {
//		if err != nil {
//			zap.S().Errorf("Error in client pump: %v", err)
//			break
//		}
//	}
//	zap.S().Infof("finished serving client %v", client.ID)
//}
