package enum

const (
	FOLLOWING int = iota
	FOLLOWER
	BLOCKED
	MUTED
)

const (
	DO_FOLLOW int = iota
	DO_UNFOLLOW
	DO_BLOCK
	DO_UNBLOCK
	DO_MUTE
	DO_UNMUTE
)
