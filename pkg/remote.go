package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mikkeloscar/aur"
	pkgbuild "github.com/mikkeloscar/gopkgbuild"
)

// NewRemotePkgs creates a Pkg slice from information returned from AUR RPC.
func NewRemotePkgs(pkg []aur.Pkg) []Pkg {
	var r []Pkg
	for i := range pkg {
		r = append(r, &remotePkg{&pkg[i]})
	}
	return r
}

type remotePkg struct {
	pkg *aur.Pkg
}

func (p *remotePkg) Name() string {
	return p.pkg.Name
}

func (p *remotePkg) Version() *pkgbuild.CompleteVersion {
	version, _ := pkgbuild.NewCompleteVersion(p.pkg.Version)
	return version
}

func (p *remotePkg) LocalPKGBUILD() string {
	return ""
}

func (p *remotePkg) URL() string {
	return p.pkg.URL
}

func (p *remotePkg) Sources() ([]string, error) {
	resp, err := http.Get("https://aur.archlinux.org/cgit/aur.git/plain/.SRCINFO?h=" + p.pkg.PackageBase)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch .SRCINFO for %s: %w", p.pkg.Name, err)
	}
	defer resp.Body.Close()
	pkgbuildBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch .SRCINFO for %s: %w", p.pkg.Name, err)
	}
	pkg, err := pkgbuild.ParseSRCINFOContent(pkgbuildBytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse .SRCINFO for %s: %w", p.pkg.Name, err)
	}
	return pkg.Source, nil
}

func (p *remotePkg) OutOfDate() bool {
	return p.pkg.OutOfDate > 0
}
