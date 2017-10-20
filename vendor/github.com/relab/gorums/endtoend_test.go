package gorums

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const (
	gorumsBaseImport = "github.com/relab/gorums"
	devImport        = gorumsBaseImport + "/" + devDir
	e2eTestDevImport = gorumsBaseImport + "/" + e2eTestDir + "/" + devDir

	devDir     = "dev"
	e2eTestDir = "e2etest"

	storageProtoFile = "storage.proto"

	devStorageProtoRelPath = devDir + "/" + storageProtoFile

	protoc        = "protoc"
	protocIFlag   = "-I=../../../:."
	protocOutFlag = "--gorums_out=plugins=grpc+gorums:"
)

func run(t *testing.T, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func runAndCaptureOutput(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%v\n%s", err, out)
	}
	return bytes.TrimSuffix(out, []byte{'\n'}), nil
}

const (
	protocVersionPrefix  = "libprotoc "
	currentProtocVersion = "3.3.2"
)

func checkProtocVersion(t *testing.T) {
	t.Helper()
	out, err := runAndCaptureOutput("protoc", "--version")
	if err != nil {
		t.Skipf("skipping test due to protoc error: %v", err)
	}
	gotVersion := string(out)
	gotVersion = strings.TrimPrefix(gotVersion, protocVersionPrefix)
	if gotVersion != currentProtocVersion {
		t.Skipf("skipping test due to old protoc version, got %q, required is %q", gotVersion, currentProtocVersion)
	}
}

var devFilesToCopy = []struct {
	name              string
	devPath, testPath string
	rewriteImport     bool
}{
	{
		"config_qc_test.go",
		"", "",
		true,
	},
	{
		"config_qc_secure_test.go",
		"", "",
		true,
	},
	{
		"config_qc_qspecs_test.go",
		"", "",
		true,
	},
	{
		"node_test.go",
		"", "",
		false,
	},
	{
		"mgr_test.go",
		"", "",
		true,
	},
	{
		"storage_server_udef.go",
		"", "",
		false,
	},
}

// copy copies the from file to the to file.
func copy(to, from string) error {
	toFd, err := os.Create(to)
	if err != nil {
		return err
	}
	defer toFd.Close()
	fromFd, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFd.Close()
	_, err = io.Copy(toFd, fromFd)
	return err
}

// Based on https://github.com/tools/godep/blob/master/rewrite.go.
// https://github.com/tools/godep/blob/master/License
func rewriteGoFile(name, old, new string) error {
	printerConfig := &printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	var changed bool
	for _, s := range f.Imports {
		iname, ierr := strconv.Unquote(s.Path.Value)
		if ierr != nil {
			return err // can't happen
		}
		if iname == old {
			s.Path.Value = strconv.Quote(new)
			changed = true
		}
	}

	if !changed {
		return errors.New("no import changed for file")
	}

	var buffer bytes.Buffer
	if err = printerConfig.Fprint(&buffer, fset, f); err != nil {
		return err
	}
	fset = token.NewFileSet()
	f, err = parser.ParseFile(fset, name, &buffer, parser.ParseComments)
	if err != nil {
		return err
	}
	ast.SortImports(fset, f)

	tpath := name + ".temp"
	t, err := os.Create(tpath)
	if err != nil {
		return err
	}
	if err = printerConfig.Fprint(t, fset, f); err != nil {
		return err
	}
	if err = t.Close(); err != nil {
		return err
	}
	// This is required before the rename on windows.
	if err = os.Remove(name); err != nil {
		return err
	}
	return os.Rename(tpath, name)
}

// TestEndToEnd runs the test suite from dev on the output generated by the gorums plugin.
func TestEndToEnd(t *testing.T) {
	checkProtocVersion(t)

	// Create temporary test dir.
	err := os.MkdirAll(e2eTestDir, 0755)
	if err != nil {
		t.Errorf("%v", err)
	}
	defer os.RemoveAll(e2eTestDir)

	// Run the proto compiler.
	run(t, protoc, protocIFlag, protocOutFlag+e2eTestDir, devStorageProtoRelPath)

	// Set file paths.
	e2eTestDevDirPath := filepath.Join(e2eTestDir, devDir)
	for i, file := range devFilesToCopy {
		devFilesToCopy[i].devPath = filepath.Join(devDir, file.name)
		devFilesToCopy[i].testPath = filepath.Join(e2eTestDevDirPath, file.name)
	}

	// Copy relevant files from dev dir.
	for _, file := range devFilesToCopy {
		err := copy(file.testPath, file.devPath)
		if err != nil {
			t.Fatalf("error copying file %q: %v", file.devPath, err)
		}
	}

	// Rewrite one import path for some test files.
	for _, file := range devFilesToCopy {
		if !file.rewriteImport {
			continue
		}
		err := rewriteGoFile(file.testPath, devImport, e2eTestDevImport)
		if err != nil {
			t.Fatalf("error rewriting import for file %q: %v", file.testPath, err)
		}
	}

	// Run go test. Use an explicit base for localhost ports to avoid
	// "address already in use" errors with regular 'dev' package tests
	// that may run concurrently.
	run(t, "go", "test", e2eTestDevImport, "-portbase=30000")
}