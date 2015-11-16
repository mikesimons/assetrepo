package assetrepo

import "os"

type AssetRepo interface {
	Get(name string) ([]byte, error)
	MustGet(name string) []byte
	Names() []string
	Info(name string) (os.FileInfo, error)
	Dir(name string) ([]string, error)
}

type getFunc func(name string) ([]byte, error)
type mustGetFunc func(name string) []byte
type namesFunc func() []string
type infoFunc func(name string) (os.FileInfo, error)
type dirFunc func(name string) ([]string, error)

func NewAdapter(f1 getFunc, f2 mustGetFunc, f3 namesFunc, f4 infoFunc, f5 dirFunc) AssetRepo {
	return &AssetRepoAdapter{f1, f2, f3, f4, f5}
}

type AssetRepoAdapter struct {
	get     getFunc
	mustGet mustGetFunc
	names   namesFunc
	info    infoFunc
	dir     dirFunc
}

func (arf *AssetRepoAdapter) Get(name string) ([]byte, error) {
	return arf.get(name)
}

func (arf *AssetRepoAdapter) MustGet(name string) []byte {
	return arf.mustGet(name)
}

func (arf *AssetRepoAdapter) Names() []string {
	return arf.names()
}

func (arf *AssetRepoAdapter) Info(name string) (os.FileInfo, error) {
	return arf.info(name)
}

func (arf *AssetRepoAdapter) Dir(name string) ([]string, error) {
	return arf.dir(name)
}
