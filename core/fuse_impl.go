//go:build !windows

package core

import (
	"os"
	"os/exec"

	"github.com/jacobsa/fuse"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isFuseTInstalled() bool {
	// fusego checks FUSE_NFSSRV_PATH first.
	if env := os.Getenv("FUSE_NFSSRV_PATH"); env != "" {
		if fileExists(env) {
			return true
		}
	}

	// Prefer whatever is in PATH (Homebrew often installs to /opt/homebrew/bin).
	if p, err := exec.LookPath("go-nfsv4"); err == nil && p != "" {
		_ = os.Setenv("FUSE_NFSSRV_PATH", p)
		return true
	}

	// Common fallback locations.
	for _, p := range []string{
		"/usr/local/bin/go-nfsv4",
		"/opt/homebrew/bin/go-nfsv4",
	} {
		if fileExists(p) {
			_ = os.Setenv("FUSE_NFSSRV_PATH", p)
			return true
		}
	}

	return false
}

func preferredFuseImpl() fuse.FUSEImpl {
	// Always use FUSE-T. The function still exists to detect and set
	// FUSE_NFSSRV_PATH if go-nfsv4 is found in non-standard locations.
	isFuseTInstalled()
	return fuse.FUSEImplFuseT
}
