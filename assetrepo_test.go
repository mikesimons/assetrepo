package assetrepo_test

import (
	"fmt"
	"os"
	"time"

	. "github.com/mikesimons/assetrepo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (f *FakeFileInfo) Name() string       { return f.name }
func (f *FakeFileInfo) Size() int64        { return f.size }
func (f *FakeFileInfo) Mode() os.FileMode  { return f.mode }
func (f *FakeFileInfo) ModTime() time.Time { return f.modTime }
func (f *FakeFileInfo) IsDir() bool        { return f.isDir }
func (f *FakeFileInfo) Sys() interface{}   { return nil }

func genericadapter() AssetRepo {
	n := []string{"test", "n1", "n2", "n3"}
	return adapter(n)
}

func adapter(names []string) AssetRepo {
	getfunc := func(name string) ([]byte, error) {
		for _, n := range names {
			if n == name {
				return []byte(name), nil
			}
		}
		return nil, fmt.Errorf(name)
	}
	mustgetfunc := func(name string) []byte {
		return []byte(name)
	}
	namesfunc := func() []string {
		return names
	}
	infofunc := func(name string) (os.FileInfo, error) {
		f := &FakeFileInfo{
			name:    name,
			size:    0,
			mode:    0777,
			modTime: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
			isDir:   false,
		}
		return f, fmt.Errorf(name)
	}
	dirfunc := func(name string) ([]string, error) {
		dir := []string{"f1", "f2", "f3"}
		return dir, fmt.Errorf(name)
	}

	return NewAdapter(getfunc, mustgetfunc, namesfunc, infofunc, dirfunc)
}

func testrepos() []AssetRepo {
	r1 := []string{"repo1/file1", "repo1/file2", "repo1/file3"}
	r2 := []string{"repo1/repo2/test", "repo2/test"}
	r3 := []string{"repo1/repo2/test/repo3", "repo3/test", "repo1/repo3"}
	ret := []AssetRepo{
		adapter(r1),
		adapter(r2),
		adapter(r3),
	}
	return ret
}

var _ = Describe("Assetrepo", func() {
	Describe("NewAdapter", func() {
		It("should produce adapter that proxies Get to getfunc", func() {
			adapter := genericadapter()
			valStr := "test"
			valBytes := []byte(valStr)

			result, err := adapter.Get(valStr)
			Expect(result).Should(Equal(valBytes))
			Expect(err).Should(BeNil())

			result, err = adapter.Get("not/exist")
			Expect(result).Should(BeNil())
			Expect(err.Error()).Should(Equal("not/exist"))
		})

		It("should produce adapter that proxies MustGet to mustgetfunc", func() {
			adapter := genericadapter()
			valStr := "test"
			valBytes := []byte(valStr)

			result := adapter.MustGet(valStr)
			Expect(result).Should(Equal(valBytes))
		})

		It("should produce adapter that proxies Names to namesfunc", func() {
			adapter := genericadapter()
			expected := []string{"test", "n1", "n2", "n3"}
			result := adapter.Names()
			Expect(result).Should(Equal(expected))
		})

		It("should produce adapter that proxies Info to infofunc", func() {
			adapter := genericadapter()
			result, err := adapter.Info("test")
			Expect(result.Name()).Should(Equal("test"))
			Expect(err.Error()).Should(Equal("test"))
		})

		It("should produce adapter that proxies Dir to dirfunc", func() {
			adapter := genericadapter()
			expected := []string{"f1", "f2", "f3"}
			result, err := adapter.Dir("test")
			Expect(result).Should(Equal(expected))
			Expect(err.Error()).Should(Equal("test"))
		})
	})

	Describe("Layered", func() {
		Describe("NewLayered", func() {
			It("should produce a layered repo", func() {
				Expect(NewLayered()).Should(BeAssignableToTypeOf(&Layered{}))
			})
		})

		Describe("NewLayeredWithRepos", func() {
			It("should produce a layered repo", func() {
				Expect(NewLayered()).Should(BeAssignableToTypeOf(&Layered{}))
			})

			It("should add repos specified", func() {
				repo := NewLayeredWithRepos(testrepos())
				subject, _ := repo.Get("repo1/file1")
				Expect(subject).Should(Equal([]byte("repo1/file1")))
				subject, _ = repo.Get("repo2/test")
				Expect(subject).Should(Equal([]byte("repo2/test")))
				subject, _ = repo.Get("repo3/test")
				Expect(subject).Should(Equal([]byte("repo3/test")))
				subject, _ = repo.Get("not/exist")
				Expect(subject).Should(BeNil())
			})
		})

		Describe("AddRepo", func() {
			It("should add specified repo to repolist", func() {
				repo := NewLayeredWithRepos(testrepos())
				names := []string{"addrepo/test"}
				newRepo := adapter(names)
				repo.AddRepo(newRepo)
				subject, _ := repo.Get("addrepo/test")
				Expect(subject).Should(Equal([]byte("addrepo/test")))
			})
		})

		Describe("Get", func() {
			It("should get asset contents from string handle", func() {
				repo := NewLayeredWithRepos(testrepos())
				content, _ := repo.Get("repo2/test")
				Expect(content).Should(Equal([]byte("repo2/test")))
			})

			It("should return nil for content and error if handle could not be found", func() {
				repo := NewLayeredWithRepos(testrepos())
				content, err := repo.Get("not/exist")
				Expect(content).Should(BeNil())
				Expect(err.Error()).Should(Equal("Asset not/exist not found"))
			})
		})

		Describe("MustGet", func() {
			It("should get asset contents from string handle", func() {
				repo := NewLayeredWithRepos(testrepos())
				content := repo.MustGet("repo2/test")
				Expect(content).Should(Equal([]byte("repo2/test")))
			})

			It("should panic if asset could not be found", func() {
				repo := NewLayeredWithRepos(testrepos())
				Expect(func() { repo.MustGet("not/exist") }).Should(Panic())
			})
		})

		Describe("Info", func() {
			It("should return file info for given handle", func() {
				repo := NewLayeredWithRepos(testrepos())
				info, _ := repo.Info("repo2/test")
				Expect(info.Name()).Should(Equal("repo2/test"))
			})

			It("should return nil and an error if handle could not be found", func() {
				repo := NewLayeredWithRepos(testrepos())
				info, err := repo.Info("not/exist")
				Expect(info).Should(BeNil())
				Expect(err.Error()).Should(Equal("Asset info for not/exist not found"))
			})
		})

		Describe("Names", func() {
			It("should return a sorted list of files", func() {
				repo := NewLayeredWithRepos(testrepos())
				expected := []string{
					"repo1/file1",
					"repo1/file2",
					"repo1/file3",
					"repo1/repo2/test",
					"repo1/repo2/test/repo3",
					"repo1/repo3",
					"repo2/test",
					"repo3/test",
				}

				names := repo.Names()
				Expect(names).Should(Equal(expected))
			})

			It("should update when a new repo is added", func() {
				repo := NewLayeredWithRepos(testrepos())
				expected := []string{
					"repo1/file1",
					"repo1/file2",
					"repo1/file3",
					"repo1/repo2/test",
					"repo1/repo2/test/repo3",
					"repo1/repo3",
					"repo2/test",
					"repo3/test",
				}

				names := repo.Names()
				Expect(names).Should(Equal(expected))

				extraNames := []string{"addrepo/test"}
				newRepo := adapter(extraNames)
				repo.AddRepo(newRepo)

				expected = []string{
					"addrepo/test",
					"repo1/file1",
					"repo1/file2",
					"repo1/file3",
					"repo1/repo2/test",
					"repo1/repo2/test/repo3",
					"repo1/repo3",
					"repo2/test",
					"repo3/test",
				}

				newNames := repo.Names()
				Expect(newNames).Should(Equal(expected))
			})
		})

		Describe("Dir", func() {
			It("should return a 'directory' listing for a handle prefix", func() {
				repo := NewLayeredWithRepos(testrepos())
				expected := []string{
					"file1",
					"file2",
					"file3",
					"repo2",
					"repo3",
				}

				result, _ := repo.Dir("repo1/")
				Expect(result).Should(Equal(expected))
			})
		})
	})
})
