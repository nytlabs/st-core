package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/nytlabs/st-core/core"
	"github.com/thejerf/suture"
)

// BlockLedger is the information the API keeps about the block
// TODO Token and ID should probably be the same thing
type BlockLedger struct {
	Name      string
	BlockType string
	Token     suture.ServiceToken
	ID        int
	Block     *core.Block
}

// GetID returns the block's ID
func (b BlockLedger) GetID() int {
	return b.ID
}

// GetName returns the block's Name
func (b BlockLedger) GetName() string {
	return b.Name
}

// MarshalJSON returns the JSON representation of the block ledger
func (b BlockLedger) MarshalJSON() ([]byte, error) {
	out := map[string]interface{}{
		"id":        strconv.Itoa(b.ID),
		"name":      b.Name,
		"blockType": b.BlockType,
	}
	return json.Marshal(out)
}

func (s *Server) ListBlocks() []BlockLedger {
	var blocks []BlockLedger
	s.Lock()
	for _, b := range s.blocks {
		blocks = append(blocks, *b)
	}
	s.Unlock()
	return blocks
}

func (s *Server) BlockIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(s.ListBlocks()); err != nil {
		panic(err)
	}
}

// CreateBlockHandler responds to a POST request to instantiate a new block and add it to the Server.
// TODO currently all blocks start off life in the root group, which may be a bit limiting.
func (s *Server) BlockCreate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(body)
	}
	var m BlockLedger
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Panic(err)
	}
	blockSpec, ok := s.library[m.BlockType]
	if !ok {
		log.Fatal("unknown block type")
	}
	b := core.NewBlock(blockSpec)
	m.ID = s.GetNextID()
	m.Block = b
	// we need to introduce the block to our running supervisor
	m.Token = s.supervisor.Add(b)
	// and we need to assign it to a group
	s.groups[0].AddNode(&m)
	s.blocks[m.ID] = &m
	//broadcast for state update
	out, _ := json.Marshal(m)
	s.broadcast <- out
	// write HTTP response
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
func (s *Server) BlockDelete(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) BlockModify(w http.ResponseWriter, r *http.Request) {
}
