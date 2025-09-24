package api

import (
	"github.com/go-kit/log"
	"github.com/labstack/echo/v4"
	"net/http"
	"github.com/772005himanshu/Mingo-Blockchain/core"
	"fmt"
	"strconv" // string Converter
	"time"
)

type Block struct {
	Version uint32
	DataHash string
	PrevBlockHash string
	Height uint32
	Timestamp time.Time
}


type ServerConfig struct {
	Logger log.Logger
	ListenAddr string 
}


type Server struct {
	ServerConfig
	bc *core.Blockchain
}


func NewServer(cfg ServerConfig, bc *core.Blockchain) *Server{
	return &Server{
		ServerConfig: cfg,
		bc : bc,
	}

}

func (s *Server) Start() error {
	e := echo.New()

	e.GET("/block/:hashorid", s.handleGetBlock)

	return e.Start(s.ListenAddr)
}

func (s *Server) handleGetBlock(c echo.Context) error {
	hashOrID := c.Param("hashorid")
	height, err := strconv.Atoi(hashOrID) // if this has then return the error , because the hash never converted back to the ID again 
	if err == nil {
		block, err := s.bc.GetBlock(uint32(height))
		if err != nil {
			return err
		}  

		// Taking Our Block Get Converted to the Json Block For Postman -> Getting the String From It Converted Easily 
		jsonBlock := Block {
			Version: block.Header.Version,
			Height: block.Header.Height,
			DataHash: block.Header.DataHash.String(), //  Here we convert it to the String Then Sending to the JSON API Postman that all 
			PrevBlockHash: block.Header.PrevBlockHash.String(),
			// Timestamp: uint64(time.Now().UnixNano()),
		}

		
		return c.JSON(http.StatusOK, jsonBlock) // what is the use of this functionality here , this go the backed (Postman) update the srever Request JSON like that 
	}
	fmt.Println(hashOrID)
	return c.JSON(http.StatusOK, map[string]any{"msg": "it works!"})
}