package server

import (
	"net"
	"net/http"
	"polley/db"
	"polley/models"
)

type pollResponse struct {
	*models.Poll
	VoteAllowed bool `json:"voteAllowed"`
}

func newPollResponse(s *Server, r *http.Request, poll *models.Poll) pollResponse {
	return pollResponse{poll, isVoteAllowed(poll.UUID, poll.Filter, r, s.ipsDB)}
}

func isVoteAllowed(uuid string, filter string, r *http.Request, ipsDB db.IPsDB) bool {
	switch filter {
	case "ip":
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return true
		}
		return ipsDB.IsVoteAllowedForIP(uuid, ip)

	case "cookie":
		_, err := r.Cookie(uuid)
		if err != nil {
			return true
		}
		return false

	}
	return true
}
