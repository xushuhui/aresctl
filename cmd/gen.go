package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// DBType represents supported database types.
type DBType string

const (
	dbMySQL DBType = "mysql"
)

const (
	defaultQueryPath      = "./dao/query"
	defaultGormConfigPath = "gorm.yaml"
)

// CmdParams defines gen command config.
type CmdParams struct {
	DSN               string   `yaml:"dsn"`
	DB                string   `yaml:"db"`
	Tables            []string `yaml:"tables"`
	OnlyModel         bool     `yaml:"onlyModel"`
	OutPath           string   `yaml:"outPath"`
	OutFile           string   `yaml:"outFile"`
	WithUnitTest      bool     `yaml:"withUnitTest"`
	ModelPkgName      string   `yaml:"modelPkgName"`
	FieldNullable     bool     `yaml:"fieldNullable"`
	FieldCoverable    bool     `yaml:"fieldCoverable"`
	FieldWithIndexTag bool     `yaml:"fieldWithIndexTag"`
	FieldWithTypeTag  bool     `yaml:"fieldWithTypeTag"`
	FieldSignable     bool     `yaml:"fieldSignable"`
}

func (c *CmdParams) revise() *CmdParams {
	if c == nil {
		return c
	}
	if c.DB == "" {
		c.DB = string(dbMySQL)
	}
	if c.OutPath == "" {
		c.OutPath = defaultQueryPath
	}
	if len(c.Tables) == 0 {
		return c
	}

	tableList := make([]string, 0, len(c.Tables))
	for _, tableName := range c.Tables {
		trimmed := strings.TrimSpace(tableName)
		if trimmed == "" {
			continue
		}
		tableList = append(tableList, trimmed)
	}
	c.Tables = tableList
	return c
}

// YamlConfig defines gen yaml structure.
type YamlConfig struct {
	Version  string     `yaml:"version"`
	Database *CmdParams `yaml:"database"`
}

func parseCmdFromYaml(path string) (*CmdParams, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config failed: %w", err)
	}
	defer file.Close()

	var yamlConfig YamlConfig
	if err := yaml.NewDecoder(file).Decode(&yamlConfig); err != nil {
		return nil, fmt.Errorf("decode config failed: %w", err)
	}
	if yamlConfig.Database == nil {
		return nil, fmt.Errorf("config missing database section")
	}
	return yamlConfig.Database, nil
}

func connectDB(t DBType, dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn cannot be empty")
	}
	if t != dbMySQL {
		return nil, fmt.Errorf("unsupported db %q (only mysql is supported)", t)
	}
	return gorm.Open(mysql.Open(dsn))
}

func genModels(g *gen.Generator, db *gorm.DB, tables []string) ([]any, error) {
	if len(tables) == 0 {
		allTables, err := db.Migrator().GetTables()
		if err != nil {
			return nil, fmt.Errorf("get tables failed: %w", err)
		}
		tables = allTables
	}

	models := make([]any, len(tables))
	for i, tableName := range tables {
		models[i] = g.GenerateModel(tableName)
	}
	return models, nil
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate GORM models and query code",
	Long:  "Generate GORM models and query code by database schema using gorm.yaml in current directory",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := parseCmdFromYaml(defaultGormConfigPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", defaultGormConfigPath, err)
			os.Exit(1)
		}

		if err := runGen(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error running gen: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✓ Generated gorm code successfully")
	},
}

func runGen(config *CmdParams) error {
	config = config.revise()
	if config == nil {
		return fmt.Errorf("parse config failed")
	}

	db, err := connectDB(DBType(config.DB), config.DSN)
	if err != nil {
		return fmt.Errorf("connect db failed: %w", err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           config.OutPath,
		OutFile:           config.OutFile,
		ModelPkgPath:      config.ModelPkgName,
		WithUnitTest:      config.WithUnitTest,
		FieldNullable:     config.FieldNullable,
		FieldCoverable:    config.FieldCoverable,
		FieldWithIndexTag: config.FieldWithIndexTag,
		FieldWithTypeTag:  config.FieldWithTypeTag,
		FieldSignable:     config.FieldSignable,
	})

	g.UseDB(db)

	g.WithDataTypeMap(map[string]func(gorm.ColumnType) (dataType string){
		"int": func(columnType gorm.ColumnType) (dataType string) {
			return "int32"
		},

		// bool mapping
		"tinyint": func(columnType gorm.ColumnType) (dataType string) {
			return "int32"
		},
	})
	models, err := genModels(g, db, config.Tables)
	if err != nil {
		return err
	}

	if !config.OnlyModel {
		g.ApplyBasic(models...)
	}
	g.Execute()
	return nil
}

func init() {
	rootCmd.AddCommand(genCmd)
}
