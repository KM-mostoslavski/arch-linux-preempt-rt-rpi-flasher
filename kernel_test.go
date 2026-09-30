package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRTPkgInRepoBeside(t *testing.T) {
	parent := t.TempDir()
	touch := func(rel string) string {
		t.Helper()
		p := filepath.Join(parent, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, nil, 0o644); err != nil {
			t.Fatal(err)
		}
		return p
	}
	touch("flasher/rpi4-flash")
	if got, err := findRTPkg(repoDirs(parent)); err == nil {
		t.Fatalf("found %q with no package anywhere", got)
	}

	// the package built at the root of the PKGBUILD repository's checkout
	root := touch("linux-rt-arm/linux-rt-arm-7.2.8-1-aarch64.pkg.tar.xz")
	touch("linux-rt-arm/linux-rt-arm-headers-7.2.8-1-aarch64.pkg.tar.xz")
	if got, err := findRTPkg(repoDirs(parent)); err != nil || got != root {
		t.Fatalf("root layout: got %q, %v; want %q", got, err, root)
	}

	// 7.2.7 wins over a newer version, wherever it is
	sub := touch("linux-rt-arm/PKGBUILDs/linux-rt-arm/linux-rt-arm-7.2.7-2-aarch64.pkg.tar.zst")
	if got, err := findRTPkg(repoDirs(parent)); err != nil || got != sub {
		t.Fatalf("PKGBUILDs layout: got %q, %v; want %q", got, err, sub)
	}
}
