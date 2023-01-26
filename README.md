# randfs: file and directory tree generator for Go

*Status: WIP, library design in progress, no ready for use API yet*

The purpose of the library is help with testing of applications that
deal with large file trees. Also it could be useful for filesystems
load testing.

The library offers easy ways to generate:
- directory trees of arbitrary depth
  - set templates for directory names or randomize them
- fill directories with files
  - set templates for filenames or randomize them
  - limit filesizes
  - set file access masks and owners/groups
  - fill files with zeroes or random data

It could be used with functions of [io/fs](https://pkg.go.dev/io/fs) package from Go standard
library.

Simple examples:

```go
import (
	"github.com/grafov/randfs"
	"github.com/grafov/randfs/name"
	"github.com/grafov/randfs/dir"
)

// Makes directory tree with default settings for directories and
// files. Directories created with random names of vary length that
// consist of digits and latin letters in mixed case. Files created
// with zero size. Filenames by default use the same rules as the
// directories. Example above set 3 levels of directory tree with
// empty files placed on each level.
func makeTree() {
	randfs.Make("/tmp", randfs.Depth(3), randfs.LevelLimit(3))
}

// Makes directory tree with 3 levels deep and 3 directories per
// level. Each directory has random name that lowercased and uses from only
// letters used for hexadecimal numbers (0..f).
func makeEmptyTree() {
	// New template for a name. It could be name for file or directory
	// no matter. In this case it used for directories naming.
	n := name.New(name.LowerCase(), name.UseHex(), name.Limits(4, 8))
	// New directory optionally sets custom access mode and uses
	d := dir.New(n, dir.Mode(0o770))
	randfs.Make("/tmp", randfs.Entries(d), randfs.Depth(3), randfs.Limit(3))
}

// Makes tree with a mix of directory and file entries on each
// level. With depth of 3 levels.
func makeTreeWithFiles() {
	// file template
	f := file.New(
		name.New(name.LowerCase()),
		file.SizeRand(0, 1024),
		file.Mode(0o660),
	)
	// directory template
	d := dir.New(
		name.New(name.UpperCase(), name.NameLength(3, 8), name.Limits(3)),
		dir.Mode(0o770),
	)
	// Make tree with a mix of directory and file entries on each level.
	randfs.Make("/tmp", randfs.Entries(d, f), randfs.Depth(3), randfs.LevelLimit(3))
}
```

## Full features list

- setup templates for file and directory names
  - set fixed or random length of names
  - limit alphabet used for names
  - upper/lower/mixed case for names
  - optional file extensions
- set access and modes for files and directories
- fill files with random data or zeroes up to limit
- make directories trees with limits
  - limit depth of directory tree
  - limit number of subdirectories per level
  - limit number of files in a directory
- empty directory tree or filled with files
- implement function similar to os.WalkDir() from standard library
  but for making directory tree

## Complex use cases

```go
func makeCustomTree() {
	fn := func(parent string, info randfs.TreeStat, randfs.DirEntry, err error) error {
		switch {
		case info.Level == 0:
			return randfs.MakeDirs("usr/local", "var", "lib", "home", "opt", "tmp")
		case info.Level > 0 && info.Level < 3:
			// make dirs
			fname := name.New(name.UpperCase(), name.UseHex(), name.LimitRand(4, 9))
			d := dir.New(fname, dir.Mode(0x770))
			return d.MakeLimit(100)
		case info.Level == 3:
			// for the deepest level fill dirs with files
			fname := name.New(randfs.UseNumbers(), randfs.LowerCase(), randfs.LimitRand(3, 16), randfs.UseAlphabet("abcdefghijklmn"))
			f := file.New(fname, randfs.PrefixRand(3), randfs.NameLength(3, 8), randfs.RandomExt(), randfs.LowerCase())
			return f.MakeRand(0, 100)
		}
	}
	randfs.Make("/tmp", randfs.WalkTree(f))
}
```

See more detailed cases in the unit tests.

## Dependencies

The library uses only Go standard library in runtime. There are
external dependencies used only for unit tests:

- https://github.com/smartystreets/goconvey
- https://github.com/spf13/afero

## Similar packages

A couple of abandoned projects:
- https://github.com/jbenet/go-random-files
- https://github.com/techjacker/randomfs

Also there is excellent https://github.com/brianvoe/gofakeit library
but for other purposes. As a part of API it has several functions
for files generation.
