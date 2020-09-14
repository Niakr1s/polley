package storage

import "polley/controllers"

// Storage is unified storage, used by server.
type Storage struct {
	Polls controllers.PollController
	Ips   controllers.IPsController
}

// NewStorage constructs new Storage.
func NewStorage(pollController controllers.PollController, ipsController controllers.IPsController) *Storage {
	return &Storage{
		Polls: pollController,
		Ips:   ipsController,
	}
}
