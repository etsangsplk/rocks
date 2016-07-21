package consistency

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ind9/rocks/cmd"
	"github.com/ind9/rocks/cmd/backup"
	"github.com/ind9/rocks/cmd/restore"
	"github.com/ind9/rocks/cmd/testutils"
	"github.com/stretchr/testify/assert"
)

func TestConsitency(t *testing.T) {
	backupDir, err := ioutil.TempDir("", "ind9-rocks-backup")
	assert.NoError(t, err)
	defer os.RemoveAll(backupDir)
	restoreDir, err := ioutil.TempDir("", "ind9-rocks-restore")
	assert.NoError(t, err)
	defer os.RemoveAll(restoreDir)
	dataDir, err := ioutil.TempDir("", "ind9-rocks")
	assert.NoError(t, err)
	defer os.RemoveAll(dataDir)

	testutils.WriteTestDB(t, dataDir)

	err = backup.DoBackup(dataDir, backupDir)
	assert.NoError(t, err)
	assert.True(t, testutils.Exists(filepath.Join(backupDir, cmd.LatestBackup)))

	err = restore.DoRestore(backupDir, restoreDir, restoreDir, false)
	assert.NoError(t, err)
	assert.True(t, testutils.Exists(filepath.Join(restoreDir, cmd.Current)))

	err = DoConsistency(dataDir, restoreDir)
	assert.NoError(t, err)
}

func TestRecursiveConsistency(t *testing.T) {
	const ShardCount = 3
	const DBsInEachShard = 3
	var paths []string
	for shard := 0; shard < ShardCount; shard++ {
		for db := 0; db < DBsInEachShard; db++ {
			path := fmt.Sprintf("%d/store_%d/", shard, db)
			paths = append(paths, path)
		}
	}

	baseDataDir, err := ioutil.TempDir("", "baseDataDir")
	assert.NoError(t, err)
	defer os.RemoveAll(baseDataDir)

	baseBackupDir, err := ioutil.TempDir("", "baseBackupDir")
	assert.NoError(t, err)
	defer os.RemoveAll(baseBackupDir)

	baseRestoreDir, err := ioutil.TempDir("", "baseRestoreDir")
	assert.NoError(t, err)
	defer os.RemoveAll(baseRestoreDir)

	// recursively write data
	for _, relLocation := range paths {
		err = os.MkdirAll(filepath.Join(baseDataDir, relLocation), os.ModePerm)
		assert.NoError(t, err)
		testutils.WriteTestDB(t, filepath.Join(baseDataDir, relLocation))
	}

	// recursive backup + assert it
	err = backup.DoRecursiveBackup(baseDataDir, baseBackupDir, 1)
	assert.NoError(t, err)
	for _, relLocation := range paths {
		assert.True(t, testutils.Exists(filepath.Join(baseBackupDir, relLocation, cmd.LatestBackup)))
	}

	err = restore.DoRecursiveRestore(baseBackupDir, baseRestoreDir, baseRestoreDir, 5, true)
	assert.NoError(t, err)
	for _, relLocation := range paths {
		assert.True(t, testutils.Exists(filepath.Join(baseRestoreDir, relLocation, cmd.Current)))
	}

	err = DoRecursiveConsistency(baseDataDir, baseRestoreDir, 5)
	assert.NoError(t, err)
}