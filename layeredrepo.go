package assetrepo

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type Layered struct {
	repos []AssetRepo
	index map[string]AssetRepo
	names []string
}

func NewLayered() *Layered {
	return &Layered{}
}

func NewLayeredWithRepos(repos []AssetRepo) *Layered {
	ret := &Layered{}
	for _, repo := range repos {
		ret.addRepo(repo)
	}
	ret.reindex()
	return ret
}

func (s *Layered) AddRepo(r AssetRepo) {
	s.addRepo(r)
	s.reindex()
}

func (s *Layered) addRepo(r AssetRepo) {
	s.repos = append(s.repos, r)
}

func (s *Layered) Get(name string) ([]byte, error) {
	r, found := s.index[name]
	if !found {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	return r.Get(name)
}

func (s *Layered) MustGet(name string) []byte {
	asset, err := s.Get(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}
	return asset
}

func (s *Layered) Info(name string) (os.FileInfo, error) {
	r, found := s.index[name]
	if !found {
		return nil, fmt.Errorf("Asset info for %s not found", name)
	}
	return r.Info(name)
}

func (s *Layered) Dir(dir string) ([]string, error) {
	dir = strings.Replace(dir, "\\", "/", -1)
	dirStrLen := len(dir)
	entries := make(map[string]bool)
	for _, name := range s.names {
		if len(name) < dirStrLen || !strings.HasPrefix(name, dir) {
			continue
		}

		elem := strings.SplitN(name[dirStrLen:], "/", 2)[0]
		entries[elem] = true
	}

	result := make([]string, 0, len(entries))
	for k := range entries {
		result = append(result, k)
	}
	sort.Strings(result)

	return result, nil
}

func (s *Layered) Names() []string {
	return s.names
}

func (s *Layered) reindex() {
	s.index = make(map[string]AssetRepo)
	s.names = []string{}
	for _, repo := range s.repos {
		for _, name := range repo.Names() {
			s.index[name] = repo
			s.names = append(s.names, name)
		}
	}
	sort.Strings(s.names)
}
