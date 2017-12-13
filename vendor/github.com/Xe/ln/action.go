package ln

// Action is a convenience helper for logging the "action" being performed as
// part of a log line.
//
// It is a convenience wrapper for the following:
//
//     ln.Log(ctx, fer, f, ln.Action("writing frozberry sales reports to database"))
func Action(act string) Fer {
	return F{"action": act}
}
