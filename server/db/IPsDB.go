package db

// IPsDB stores IPs for each poll.
type IPsDB interface {
	AddIPForPoll(uuid string, ip string) error
	IsVoteAllowedForIP(uuid string, ip string) bool
}
