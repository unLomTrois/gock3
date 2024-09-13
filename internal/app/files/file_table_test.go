package files

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
)

func TestGetPathTableInstance(t *testing.T) {
	resetPathTable()

	tests := []struct {
		name string
		want *pathTable
	}{
		{
			name: "Check Singleton Instance",
			want: GetPathTableInstance(), // Expect the same instance to be returned
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPathTableInstance()
			if got != tt.want { // Compare instances (pointer equality check)
				t.Errorf("GetPathTableInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pathTable_Store(t *testing.T) {
	// Reset the singleton before the test starts
	resetPathTable()

	type args struct {
		local    string
		fullpath string
	}
	tests := []struct {
		name string
		args args
		want PathTableIndex
	}{
		{
			name: "Store new path",
			args: args{
				local:    filepath.Join("local", "path"),
				fullpath: filepath.Join("full", "path"),
			},
			want: PathTableIndex{index: 0}, // Expected index after first insertion
		},
		{
			name: "Store second path",
			args: args{
				local:    filepath.Join("local2", "path"),
				fullpath: filepath.Join("full2", "path"),
			},
			want: PathTableIndex{index: 1}, // Expected index after second insertion
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use PATHTABLE to store paths, which will call GetPathTableInstance() internally
			got := PATHTABLE.Store(tt.args.local, tt.args.fullpath)
			if got.index != tt.want.index {
				t.Errorf("PATHTABLE.Store() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pathTable_store(t *testing.T) {
	resetPathTable()

	type args struct {
		local    string
		fullpath string
	}
	tests := []struct {
		name string
		pt   *pathTable
		args args
		want PathTableIndex
	}{
		{
			name: "Store new path",
			pt:   GetPathTableInstance(),
			args: args{
				local:    filepath.Join("local", "path"),
				fullpath: filepath.Join("full", "path"),
			},
			want: PathTableIndex{index: 0}, // Expected index after first insertion
		},
		{
			name: "Store second path",
			pt:   GetPathTableInstance(),
			args: args{
				local:    filepath.Join("local2", "path"),
				fullpath: filepath.Join("full2", "path"),
			},
			want: PathTableIndex{index: 1}, // Expected index after second insertion
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pt.store(tt.args.local, tt.args.fullpath)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pathTable.store() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pathTable_Store_Concurrent(t *testing.T) {
	resetPathTable()

	var wg sync.WaitGroup
	numGoroutines := 10

	// Expected number of paths stored should match the number of goroutines
	expectedPaths := numGoroutines

	// Run Store method concurrently in multiple goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			localPath := filepath.Join(fmt.Sprintf("local%d", i), "path")
			fullPath := filepath.Join(fmt.Sprintf("full%d", i), "path")

			// Use PATHTABLE.Store to store paths concurrently
			PATHTABLE.Store(localPath, fullPath)
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Check the total number of stored paths in the pathTable
	pt := GetPathTableInstance() // Get the singleton instance
	pt.mu.RLock()                // Use read lock since we are only reading data
	defer pt.mu.RUnlock()

	if len(pt.paths) != expectedPaths {
		t.Errorf("Concurrent Store failed: expected %d paths, got %d", expectedPaths, len(pt.paths))
	}
}
