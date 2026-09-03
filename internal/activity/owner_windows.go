//go:build windows

package activity

import "os"

// ownedByCurrentUser is not determined on Windows: ownership is not checked
// and permissions are left to the file system's ACLs.
func ownedByCurrentUser(os.FileInfo) (owned, known bool) { return false, false }
