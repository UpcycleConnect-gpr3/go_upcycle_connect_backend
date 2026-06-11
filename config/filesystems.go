package config

import (
	"go-upcycle_connect-backend/internal"
	"go-upcycle_connect-backend/var/filesystem"
	"os"
)

func InitFilesystem() {
	filesystem.PrivateStorage = internal.NewRustFS(
		"private",
		os.Getenv("RUSTFS_REGION"),
		os.Getenv("RUSTFS_ACCESS_KEY_ID"),
		os.Getenv("RUSTFS_SECRET_ACCESS_KEY"),
		os.Getenv("RUSTFS_ENDPOINT"),
	)

	filesystem.PublicStorage = internal.NewRustFS(
		"public",
		os.Getenv("RUSTFS_REGION"),
		os.Getenv("RUSTFS_ACCESS_KEY_ID"),
		os.Getenv("RUSTFS_SECRET_ACCESS_KEY"),
		os.Getenv("RUSTFS_ENDPOINT"),
	)
}
