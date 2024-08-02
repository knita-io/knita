package file

import (
	"io/fs"
	"sort"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/knita-io/knita/api/executor/v1"
)

type outFile struct {
	dir  bool
	path string
	dest string
}

type testSendTransport struct {
	sends []*v1.FileTransfer
}

func (t *testSendTransport) Send(transfer *v1.FileTransfer) error {
	t.sends = append(t.sends, transfer)
	return nil
}

func TestSend(t *testing.T) {

	var testFS = fstest.MapFS{
		"one/a_file.txt":               &fstest.MapFile{Data: []byte("a")},
		"one/b_file.txt":               &fstest.MapFile{Data: []byte("b")},
		"nested_1/1_file.txt":          &fstest.MapFile{Data: []byte("1")},
		"nested_1/nested_2/2_file.txt": &fstest.MapFile{Data: []byte("2")},
	}

	var table = []struct {
		fs       fs.FS
		src      string
		dest     string
		excludes []string
		out      []outFile
	}{
		{
			fs:   testFS,
			src:  "one/*",
			dest: ".",
			out: []outFile{
				{
					dir:  false,
					path: "one/a_file.txt",
					dest: "a_file.txt",
				},
				{
					dir:  false,
					path: "one/b_file.txt",
					dest: "b_file.txt",
				},
			},
		},
		{
			fs:   testFS,
			src:  "one/",
			dest: ".",
			out: []outFile{
				{
					dir:  false,
					path: "one/a_file.txt",
					dest: "a_file.txt",
				},
				{
					dir:  false,
					path: "one/b_file.txt",
					dest: "b_file.txt",
				},
			},
		},
		{
			fs:   testFS,
			src:  "one",
			dest: "two",
			out: []outFile{
				{
					dir:  true,
					path: "one",
					dest: "two",
				},
				{
					dir:  false,
					path: "one/a_file.txt",
					dest: "two/a_file.txt",
				},
				{
					dir:  false,
					path: "one/b_file.txt",
					dest: "two/b_file.txt",
				},
			},
		},
		{
			fs:   testFS,
			src:  "one",
			dest: "two/",
			out: []outFile{
				{
					dir:  true,
					path: "one",
					dest: "two/one",
				},
				{
					dir:  false,
					path: "one/a_file.txt",
					dest: "two/one/a_file.txt",
				},
				{
					dir:  false,
					path: "one/b_file.txt",
					dest: "two/one/b_file.txt",
				},
			},
		},
		{
			fs:   testFS,
			src:  "one/*",
			dest: "two/",
			out: []outFile{
				{
					dir:  false,
					path: "one/a_file.txt",
					dest: "two/a_file.txt",
				},
				{
					dir:  false,
					path: "one/b_file.txt",
					dest: "two/b_file.txt",
				},
			},
		},
		{
			fs:   testFS,
			src:  "one/",
			dest: "two/",
			out: []outFile{
				{
					dir:  false,
					path: "one/a_file.txt",
					dest: "two/a_file.txt",
				},
				{
					dir:  false,
					path: "one/b_file.txt",
					dest: "two/b_file.txt",
				},
			},
		},
		{
			fs:   testFS,
			src:  "nested_1/**",
			dest: ".",
			out: []outFile{
				{
					dir:  true,
					path: "nested_1",
					dest: "nested_1",
				},
				{
					dir:  false,
					path: "nested_1/1_file.txt",
					dest: "nested_1/1_file.txt",
				},
				{
					dir:  true,
					path: "nested_1/nested_2",
					dest: "nested_1/nested_2",
				},
				{
					dir:  false,
					path: "nested_1/nested_2/2_file.txt",
					dest: "nested_1/nested_2/2_file.txt",
				},
				{
					dir:  false,
					path: "nested_1/1_file.txt",
					dest: "1_file.txt",
				},
				{
					dir:  true,
					path: "nested_1/nested_2",
					dest: "nested_2",
				},
				{
					dir:  false,
					path: "nested_1/nested_2/2_file.txt",
					dest: "nested_2/2_file.txt",
				},
				{
					dir:  false,
					path: "nested_1/nested_2/2_file.txt",
					dest: "2_file.txt",
				},
			},
		},
		{
			fs:   testFS,
			src:  "nested_1/**/*.*",
			dest: ".",
			out: []outFile{
				{
					dir:  false,
					path: "nested_1/1_file.txt",
					dest: "1_file.txt",
				},
				{
					dir:  false,
					path: "nested_1/nested_2/2_file.txt",
					dest: "2_file.txt",
				},
			},
		},
		{
			fs:   testFS,
			src:  "one/a_file.txt",
			dest: "two/b_file.txt",
			out: []outFile{
				{
					dir:  false,
					path: "one/a_file.txt",
					dest: "two/b_file.txt",
				},
			},
		},
		{
			fs:   testFS,
			src:  "one/a_file.txt",
			dest: "two/",
			out: []outFile{
				{
					dir:  false,
					path: "one/a_file.txt",
					dest: "two/a_file.txt",
				},
			},
		},
		{
			fs:   testFS,
			src:  "one/*.txt",
			dest: "foo",
			out: []outFile{
				{
					dir:  false,
					path: "one/a_file.txt",
					dest: "foo/a_file.txt",
				},
				{
					dir:  false,
					path: "one/b_file.txt",
					dest: "foo/b_file.txt",
				},
			},
		},
		{
			fs: fstest.MapFS{
				"a_file.txt": &fstest.MapFile{Data: []byte("a")},
			},
			src:  ".",
			dest: ".",
			out: []outFile{
				{
					dir:  false,
					path: "a_file.txt",
					dest: "a_file.txt",
				},
			},
		},
		{
			fs: fstest.MapFS{
				".hidden/hidden.txt": &fstest.MapFile{Data: []byte("hidden")},
			},
			src:  "**",
			dest: ".",
			out: []outFile{
				{
					dir:  false,
					path: ".hidden/hidden.txt",
					dest: ".hidden/hidden.txt",
				},
				{
					dir:  true,
					path: ".hidden",
					dest: ".hidden",
				},
				{
					dir:  false,
					path: ".hidden/hidden.txt",
					dest: ".hidden/hidden.txt",
				},
				{
					dir:  true,
					path: ".hidden",
					dest: ".hidden",
				},
				{
					dir:  false,
					path: ".hidden/hidden.txt",
					dest: "hidden.txt",
				},
			},
		},
		{
			fs: fstest.MapFS{
				".hidden/hidden.txt": &fstest.MapFile{Data: []byte("hidden")},
			},
			src:  "**",
			dest: "foo/",
			out: []outFile{
				{
					dir:  true,
					path: ".",
					dest: "foo",
				},
				{
					dir:  false,
					path: ".hidden/hidden.txt",
					dest: "foo/.hidden/hidden.txt",
				},
				{
					dir:  true,
					path: ".hidden",
					dest: "foo/.hidden",
				},
				{
					dir:  false,
					path: ".hidden/hidden.txt",
					dest: "foo/.hidden/hidden.txt",
				},
				{
					dir:  true,
					path: ".hidden",
					dest: "foo/.hidden",
				},
				{
					dir:  false,
					path: ".hidden/hidden.txt",
					dest: "foo/hidden.txt",
				},
			},
		},
		///////////////////////////////////
		// Default dest
		///////////////////////////////////
		{
			fs:   testFS,
			src:  "one/*.txt",
			dest: "",
			out: []outFile{
				{
					dir:  false,
					path: "one/a_file.txt",
					dest: "one/a_file.txt",
				},
				{
					dir:  false,
					path: "one/b_file.txt",
					dest: "one/b_file.txt",
				},
			},
		},
		{
			fs:   testFS,
			src:  "one/",
			dest: "",
			out: []outFile{
				{
					dir:  false,
					path: "one/a_file.txt",
					dest: "one/a_file.txt",
				},
				{
					dir:  false,
					path: "one/b_file.txt",
					dest: "one/b_file.txt",
				},
			},
		},
		{
			fs:   testFS,
			src:  "one",
			dest: "",
			out: []outFile{
				{
					dir:  true,
					path: "one",
					dest: "one",
				},
				{
					dir:  false,
					path: "one/a_file.txt",
					dest: "one/a_file.txt",
				},
				{
					dir:  false,
					path: "one/b_file.txt",
					dest: "one/b_file.txt",
				},
			},
		},
		{
			fs:   testFS,
			src:  "one/a_file.txt",
			dest: "",
			out: []outFile{
				{
					dir:  false,
					path: "one/a_file.txt",
					dest: "one/a_file.txt",
				},
			},
		},
		///////////////////////////////////
		// Excludes
		///////////////////////////////////
		{
			fs: fstest.MapFS{
				"a/a.txt":      &fstest.MapFile{Data: []byte("a")},
				"b/b.txt":      &fstest.MapFile{Data: []byte("b")},
				"aa/aa.txt":    &fstest.MapFile{Data: []byte("aa")},
				"bb/bb.txt":    &fstest.MapFile{Data: []byte("bb")},
				"cc/cc.txt":    &fstest.MapFile{Data: []byte("cc")},
				"cc/dd/dd.txt": &fstest.MapFile{Data: []byte("dd")},
				"ee/ee.txt":    &fstest.MapFile{Data: []byte("cc")},
				"ee/ee/ee.txt": &fstest.MapFile{Data: []byte("dd")},
			},
			src:  ".",
			dest: "",
			excludes: []string{
				"a/a.txt", // exact file
				"ee/ee",   // exact dir
				"bb*",     // single-star glob
				"cc/*/**", // doublestar glob
			},
			out: []outFile{
				{
					dir:  true,
					path: "a",
					dest: "a",
				},
				{
					dir:  true,
					path: "b",
					dest: "b",
				},
				{
					dir:  false,
					path: "b/b.txt",
					dest: "b/b.txt",
				},
				{
					dir:  true,
					path: "aa",
					dest: "aa",
				},
				{
					dir:  false,
					path: "aa/aa.txt",
					dest: "aa/aa.txt",
				},
				{
					dir:  true,
					path: "cc",
					dest: "cc",
				},
				{
					dir:  true,
					path: "ee",
					dest: "ee",
				},
				{
					dir:  false,
					path: "ee/ee.txt",
					dest: "ee/ee.txt",
				},
			},
		},
	}

	logger, _ := zap.NewDevelopment()
	log := logger.Sugar()
	for i, scenario := range table {
		log.Infof("Testing scenario %d: src: %s, dest: %s", i, scenario.src, scenario.dest)
		trans := &testSendTransport{}
		sender := NewSender(log, scenario.fs, trans, "test", "test", WithDest(scenario.dest), WithExcludes(scenario.excludes))
		_, err := sender.Send(scenario.src)
		require.NoError(t, err)
		require.Equal(t, len(scenario.out), len(trans.sends))
		sort.Slice(scenario.out, func(i, j int) bool {
			return scenario.out[i].dest < scenario.out[j].dest
		})
		sort.Slice(trans.sends, func(i, j int) bool {
			return trans.sends[i].Header.DestPath < trans.sends[j].Header.DestPath
		})
		for i, out := range scenario.out {
			send := trans.sends[i]
			require.Equal(t, out.dest, send.Header.DestPath)
		}
	}
}
