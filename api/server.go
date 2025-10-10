package api

import (
	"github.com/go-kit/log"
	"github.com/labstack/echo/v4"
	"net/http"
	"github.com/772005himanshu/Mingo-Blockchain/core"
	"github.com/772005himanshu/Mingo-Blockchain/types"
	"strconv" // string Converter
	"encoding/hex"
	"encoding/gob"
	"fmt"
	
)

type TxResponse struct {
	TxCount uint
	Hashes []string
}

type APIError struct {
	Error string
} // Why we do this because In simple Error Handling for the API/ JSON we have to handle the process of the marshal and unmarshal on top of it 

type Block struct {
	Hash string
	Version uint32
	DataHash string
	PrevBlockHash string
	Height uint32
	Timestamp int32
	Validator string
	Signature string

	TxReponse TxResponse
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
	e.GET("/tx/:hash", s.handleGetTx)
	e.POST("/tx", s.handlePostTx)


	return e.Start(s.ListenAddr)
}

func (s *Server) handlePostTx(c echo.Context) error {
	tx := &core.Transaction{}
	if err := gob.NewDecoder(c.Request().Body).Decode(tx); err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	fmt.Printf("%+v\n", tx)

	return nil
}

func (s *Server) handleGetBlock(c echo.Context) error {
	hashOrID := c.Param("hashorid")
	height, err := strconv.Atoi(hashOrID) 
	// if the error is nil we can assume the height of the Block is given here 

	if err == nil {  
		block, err := s.bc.GetBlock(uint32(height))
		if err != nil {
			return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()}) // .Error() interface used to convert to string
		}  

		// Taking Our Block Get Converted to the Json Block For Postman -> Getting the String From It Converted Easily 
		jsonBlock := intoJSONBlock(block)

		return c.JSON(http.StatusOK, jsonBlock) // what is the use of this functionality here , this go the backed (Postman) update the srever Request JSON like that 
	}

	// Other wise assume its  the hash

	b, err := hex.DecodeString(hashOrID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	block, err := s.bc.GetBlockByHash(types.HashFromBytes(b))
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}
	
	return c.JSON(http.StatusOK, block)
}

func (s *Server) handleGetTx(c echo.Context) error {
	hash := c.Param("hash")
	b, err := hex.DecodeString(hash)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}
	tx, err := s.bc.GetTxByHash(types.HashFromBytes(b))
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, tx)
}


func intoJSONBlock(block *core.Block) Block {
	txResponse := TxResponse {
		TxCount : uint(len(block.Transactions)),
		Hashes : make([]string , len(block.Transactions)),
	}

	for i := 0 ; i < int(txResponse.TxCount); i++ {
		txResponse.Hashes[i] = block.Transactions[i].Hash(core.TxHasher{}).String()
	}

	return Block {
		Hash : block.Hash(core.BlockHasher{}).String(),
		Version: block.Header.Version,
		Height: block.Header.Height,
		DataHash: block.Header.DataHash.String(), //  Here we convert it to the String Then Sending to the JSON API Postman that all 
		PrevBlockHash: block.Header.PrevBlockHash.String(),
		Timestamp: int32(block.Header.Timestamp),
		Validator: block.Validator.Address().String(),
		Signature: block.Signature.String(),
		TxReponse: txResponse,
	}
}

