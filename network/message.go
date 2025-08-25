package network


type GetStatusMessage struct {

}

type StatusMessage struct {
	// the id of the server 
	ID string
	CurrentHeight uint32
	Version  uint32
}

