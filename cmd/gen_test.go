package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCmdParamsRevise_DefaultsAndTrimTables(t *testing.T) {
	params := &CmdParams{
		Tables: []string{" users ", "", " orders", "   "},
	}

	got := params.revise()
	if got.DB != string(dbMySQL) {
		t.Fatalf("expected default db %q, got %q", dbMySQL, got.DB)
	}
	if got.OutPath != defaultQueryPath {
		t.Fatalf("expected default outPath %q, got %q", defaultQueryPath, got.OutPath)
	}

	wantTables := []string{"users", "orders"}
	if !reflect.DeepEqual(got.Tables, wantTables) {
		t.Fatalf("unexpected tables, want %v, got %v", wantTables, got.Tables)
	}
}

func TestParseCmdFromYaml(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "gen.yml")
	content := []byte(`version: "0.1"
database:
  dsn: "user:pass@tcp(127.0.0.1:3306)/demo"
  db: "mysql"
  tables: ["users", "orders"]
  outPath: "./dao/query"
`)
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("write yaml failed: %v", err)
	}

	got, err := parseCmdFromYaml(configPath)
	if err != nil {
		t.Fatalf("parseCmdFromYaml failed: %v", err)
	}
	if got.DSN != "user:pass@tcp(127.0.0.1:3306)/demo" {
		t.Fatalf("unexpected dsn: %q", got.DSN)
	}
	if !reflect.DeepEqual(got.Tables, []string{"users", "orders"}) {
		t.Fatalf("unexpected tables: %v", got.Tables)
	}
}

func TestConnectDB_EmptyDSN(t *testing.T) {
	if _, err := connectDB(dbMySQL, ""); err == nil {
		t.Fatal("expected error when dsn is empty")
	}
}

func TestConnectDB_UnsupportedDB(t *testing.T) {
	if _, err := connectDB(DBType("postgres"), "dsn"); err == nil {
		t.Fatal("expected error for unsupported db")
	}
}
