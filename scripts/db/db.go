package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gookit/color"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

// ModelStruct for generating model file
type ModelStruct struct {
	ModelName string
}

func main() {
	_ = godotenv.Load("db.env")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"), os.Getenv("DB_SSL"))
	db, err := sql.Open("postgres", connectionString)
	defer db.Close()
	m, err := migrate.New(
		"file://migrations",
		connectionString,
	)
	if err != nil {
		color.Error.Tips(err.Error())
		os.Exit(1)
	}
	app := &cli.App{

		Name:  "db",
		Usage: "database related tasks!",
		Commands: []*cli.Command{
			{
				Name:    "migrate",
				Aliases: []string{"m"},
				Usage:   "Run migrations",
				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:    "version",
						Aliases: []string{"v"},
						Usage:   "versions upto which you want to migrate",
					},
					&cli.IntFlag{
						Name:    "steps",
						Aliases: []string{"s"},
						Usage:   "versions upto which you want to migrate",
					},
					&cli.BoolFlag{
						Name:    "down",
						Aliases: []string{"d"},
						Usage:   "if you have to migrate down",
					},
					&cli.BoolFlag{
						Name:    "force",
						Aliases: []string{"f"},
						Usage:   "force the specified version",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Bool("force") {
						if c.Uint("version") == 0 && c.Int("steps") == 0 {
							color.Error.Tips("version or steps required")
							os.Exit(1)
						} else if c.Uint("version") != 0 {
							if c.Int("steps") != 0 {
								color.Warn.Tips("Version specified, steps ignored")
							}
							err = m.Force(int(c.Uint("version")))
						} else {

							err = m.Force(c.Int("steps"))

						}
					} else if c.Uint("version") != 0 {
						if c.Int("steps") != 0 {
							color.Warn.Tips("Version specified, steps ignored")
						}
						err = m.Migrate(c.Uint("version"))
					} else if c.Int("steps") != 0 {
						err = m.Steps(c.Int("steps"))
					} else if c.Bool("down") {
						color.Info.Tips("Migrating down")
						err = m.Down()
					} else {
						color.Info.Tips("Migrating up")
						err = m.Up()
					}
					if err == nil {
						color.Info.Tips("Migration completed")
					} else {
						color.Error.Tips(err.Error())
					}
					return err
				},
			},
			{
				Name:    "create-migration",
				Aliases: []string{"cm"},
				Usage:   "Create a new migration",
				Action: func(c *cli.Context) error {
					name := strcase.ToSnake(c.Args().Get(0))
					if name == "" {
						color.Error.Tips("Migration name is required")
					} else {
						t := time.Now().UTC()
						timestamp := fmt.Sprintf("%d%02d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
						upMigrationFileName := fmt.Sprintf("migrations/%s_%s.up.sql", timestamp, name)
						_, err = os.Create(upMigrationFileName)
						if err != nil {
							color.Error.Tips(err.Error())
							os.Exit(1)
						}
						downMigrationFileName := fmt.Sprintf("migrations/%s_%s.down.sql", timestamp, name)
						_, err = os.Create(downMigrationFileName)
						if err != nil {
							os.Remove(upMigrationFileName)
							color.Error.Tips(err.Error())
						}
						color.Info.Tips("Migration files have been created")
					}

					return nil
				},
			},
			{
				Name:    "generate-model",
				Aliases: []string{"gm"},
				Usage:   "Generate a new model with/without migration",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "skip-migration",
						Aliases: []string{"sm"},
					},
				},
				Action: func(c *cli.Context) error {
					t := time.Now().UTC()
					modelName := inflection.Singular(strcase.ToCamel(c.Args().Get(0)))
					if !c.Bool("skip-migration") {
						tableName := inflection.Plural(strcase.ToSnake(modelName))
						upStmt := fmt.Sprintf("CREATE TABLE %s(id BIGSERIAL PRIMARY KEY);", tableName)
						downStmt := fmt.Sprintf("DROP TABLE %s;", tableName)
						timestamp := fmt.Sprintf("%d%02d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
						fileName := fmt.Sprintf("migrations/%s_create_%s_table.up.sql", timestamp, tableName)
						err = ioutil.WriteFile(fileName, []byte(upStmt), 0644)
						color.Info.Tips("Up migration created")
						fileName = fmt.Sprintf("migrations/%s_create_%s_table.down.sql", timestamp, tableName)
						err = ioutil.WriteFile(fileName, []byte(downStmt), 0644)
						color.Info.Tips("Down migration created")
					}
					f, err := os.Create(fmt.Sprintf("internal/models/%s.go", modelName))
					if err == nil {
						color.Info.Tips("Model file created")
					} else {
						color.Error.Tips("Error: ", err)
					}
					modelTemplate, _ := ioutil.ReadFile("templates/model.go.tmpl")
					modelBody := string(modelTemplate)
					model := ModelStruct{
						ModelName: modelName,
					}
					templateContents, _ := template.New("model body").Parse(modelBody)

					if err == nil {
						templateContents.Execute(f, model)
					}
					return nil
				},
			},
			{
				Name:    "reset",
				Aliases: []string{"r"},
				Usage:   "Remove the list of schema migrations",

				Action: func(c *cli.Context) error {
					_, err := db.Exec("DELETE  FROM schema_migrations")
					if err == nil {
						color.Info.Tips("schema_migrations entry reset")
					}
					return err
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		color.Error.Tips(err.Error())
		os.Exit(1)
	}
}
