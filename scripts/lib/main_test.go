package slacklog_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func cleanupTmpDir(t *testing.T, path string) {
	t.Helper()

	err := os.RemoveAll(path)
	if err != nil {
		t.Fatalf("failed to cleanupTmpDir: %s", err)
	}
}

func createTmpDir(t *testing.T) string {
	t.Helper()

	path, err := ioutil.TempDir("testdata", "slacklog")
	if err != nil {
		t.Fatalf("failed to createTmpDir: %s", err)
	}
	return path
}

func dirDiff(a, b string) error {
	aInfos, err := ioutil.ReadDir(a)
	if err != nil {
		return err
	}
	bInfos, err := ioutil.ReadDir(b)
	if err != nil {
		return err
	}

	if len(aInfos) != len(bInfos) {
		return fmt.Errorf(
			"the number of files in the directory is different: (%s: %d) (%s: %d)",
			a, len(aInfos),
			b, len(bInfos),
		)
	}

	sort.Slice(aInfos, func(i, j int) bool {
		return aInfos[i].Name() >= aInfos[i].Name()
	})
	sort.Slice(bInfos, func(i, j int) bool {
		return bInfos[i].Name() >= bInfos[i].Name()
	})

	for i := range aInfos {
		if aInfos[i].Name() != bInfos[i].Name() {
			return fmt.Errorf(
				"the file name is different: %s != %s",
				filepath.Join(a, aInfos[i].Name()),
				filepath.Join(b, bInfos[i].Name()),
			)
		}
		if aInfos[i].Size() != bInfos[i].Size() {
			return fmt.Errorf(
				"the file size is different: (%s: %d) (%s: %d)",
				filepath.Join(a, aInfos[i].Name()), aInfos[i].Size(),
				filepath.Join(b, bInfos[i].Name()), bInfos[i].Size(),
			)
		}
	}
	return nil
}
