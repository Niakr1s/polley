package controllers

// IPsController stores IPs for each poll.
type IPsController interface {
	AddIPForPoll(uuid string, ip string) error
	IsVoteAllowedForIP(uuid string, ip string) bool
}
