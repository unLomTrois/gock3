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
				fullpath: filepath.Join("full", "path"),
			},
			want: PathTableIndex{index: 0}, // Expected index after first insertion
		},
		{
			name: "Store second path",
			args: args{
				fullpath: filepath.Join("full2", "path"),
			},
			want: PathTableIndex{index: 1}, // Expected index after second insertion
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use PATHTABLE to store paths, which will call GetPathTableInstance() internally
			got := PATHTABLE.Store(tt.args.fullpath)
			if got.index != tt.want.index {
				t.Errorf("PATHTABLE.Store() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pathTable_store(t *testing.T) {
	resetPathTable()

	type args struct {
		fullpath string
	}
	tests := []struct {
		name string
		pt   *pathTable
		args args
		want *PathTableIndex
	}{
		{
			name: "Store new path",
			pt:   GetPathTableInstance(),
			args: args{
				fullpath: filepath.Join("full", "path"),
			},
			want: &PathTableIndex{index: 0}, // Expected index after first insertion
		},
		{
			name: "Store second path",
			pt:   GetPathTableInstance(),
			args: args{
				fullpath: filepath.Join("full2", "path"),
			},
			want: &PathTableIndex{index: 1}, // Expected index after second insertion
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pt.store(tt.args.fullpath)
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

			fullPath := filepath.Join(fmt.Sprintf("full%d", i), "path")

			// Use PATHTABLE.Store to store paths concurrently
			PATHTABLE.Store(fullPath)
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

func Test_pathTable_LookupFullpath(t *testing.T) {
	resetPathTable()

	// Store some paths for testing.
	PATHTABLE.Store(filepath.Join("full1", "path"))
	PATHTABLE.Store(filepath.Join("full2", "path"))

	type args struct {
		index PathTableIndex
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Lookup first full path",
			args: args{
				index: PathTableIndex{index: 0},
			},
			want:    filepath.Join("full1", "path"),
			wantErr: false,
		},
		{
			name: "Lookup second full path",
			args: args{
				index: PathTableIndex{index: 1},
			},
			want:    filepath.Join("full2", "path"),
			wantErr: false,
		},
		{
			name: "Lookup out of bounds",
			args: args{
				index: PathTableIndex{index: 2},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PATHTABLE.LookupFullpath(tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("pathTable.LookupFullpath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("pathTable.LookupFullpath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pathTable_Concurrent_Read(t *testing.T) {
	// Reset the singleton before the test starts
	resetPathTable()

	// Store some paths for testing.
	PATHTABLE.Store(filepath.Join("full1", "path"))
	PATHTABLE.Store(filepath.Join("full2", "path"))
	PATHTABLE.Store(filepath.Join("full3", "path"))

	// Number of goroutines to read concurrently
	numGoroutines := 10
	var wg sync.WaitGroup

	// Expected values
	expectedPaths := []string{
		filepath.Join("full1", "path"),
		filepath.Join("full2", "path"),
		filepath.Join("full3", "path"),
	}

	// Run concurrent readers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// Cyclically pick one of the stored indexes (0, 1, or 2)
			idx := PathTableIndex{index: uint32(i % 3)}

			// Read from the path table
			localPath, err := PATHTABLE.LookupFullpath(idx)
			if err != nil {
				t.Errorf("Error reading path at index %d: %v", idx.index, err)
				return
			}

			// Check that the path read matches the expected path
			expectedPath := expectedPaths[idx.index]
			if localPath != expectedPath {
				t.Errorf("Expected path %v, but got %v", expectedPath, localPath)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
}
