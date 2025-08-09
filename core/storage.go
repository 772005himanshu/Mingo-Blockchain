package core

// Storage is only for the rpc and json rpc , if you want to retrive a block and other block askes the nodes to sync 

type Storage interface {
	Put(*Block) error
}


type MemoryStore struct {

}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}


func (s *MemoryStore) Put(b *Block) error {
	return nil
}