//go:build !windows

package activity

import (
	"os"
	"syscall"
)

// ownedByCurrentUser reports whether info's file belongs to the current
// user. known is false when the platform gives no owner.
func ownedByCurrentUser(info os.FileInfo) (owned, known bool) {
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return false, false
	}
	return int(st.Uid) == os.Getuid(), true
}
