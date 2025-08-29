package network

type GetBlocksMessage struct {
	From uint32
	// if To is 0 , then we have to return the maximum blocks
	To uint32
}


type GetStatusMessage struct {

}

type StatusMessage struct {
	// the id of the server 
	ID string
	CurrentHeight uint32
	Version  uint32
}

