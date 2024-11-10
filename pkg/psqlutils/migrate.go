package psqlutils

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
)

type ConnectionDetails struct {
	MigAppName string
	Host       string
	Port       int
	Username   string
	Password   string
	Dbname     string
	Sslmode    bool
	Options    map[string]string
}

func MigrateUp(details ConnectionDetails, migsFS embed.FS) error {
	processedUrl := BuildPostgresURL(
		details.Host,
		details.Port,
		details.Username,
		details.Password,
		details.Dbname,
		details.Sslmode,
	)
	if details.MigAppName != "" {
		processedUrl = fmt.Sprintf("%s&application_name=%s", processedUrl, "gi_auth_migrations")
	}

	deets := &pop.ConnectionDetails{
		Dialect: "postgres",
		URL:     processedUrl,
		Options: details.Options,
	}

	pop.Debug = true
	db, err := pop.NewConnection(deets)
	if err != nil {
		return fmt.Errorf("%+v", errors.Wrap(err, "opening db connection"))
	}
	defer db.Close()

	if err := db.Open(); err != nil {
		return fmt.Errorf("%+v", errors.Wrap(err, "checking database connection"))
	}

	mig, err := pop.NewMigrationBox(migsFS, db)
	if err != nil {
		return fmt.Errorf("%+v", errors.Wrap(err, "creating db migrator"))
	}

	err = mig.Status(os.Stdout)
	if err != nil {
		return fmt.Errorf("%+v", errors.Wrap(err, "migration status"))
	}

	// turn off schema dump
	mig.SchemaPath = ""

	err = mig.Up()
	if err != nil {
		return fmt.Errorf("%v", errors.Wrap(err, "running db migrations"))
	}

	err = mig.Status(os.Stdout)
	if err != nil {
		return fmt.Errorf("%+v", errors.Wrap(err, "migration status"))
	}

	return nil
}

func MigrationsContent(
	p string,
	migsFS fs.FS,
	options map[string]string,
) (map[string]string, error) {
	deets := &pop.ConnectionDetails{
		Dialect: "postgres",
		URL:     "postgres://",
		Options: options,
	}
	conn, err := pop.NewConnection(deets)
	if err != nil {
		return nil, err
	}

	contentMap := make(map[string]string)
	dirEntries, err := fs.ReadDir(migsFS, p)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		// dont handle nested for now
		if entry.IsDir() {
			continue
		}

		filePath := path.Join(p, entry.Name())
		file, err := migsFS.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		content, err := pop.MigrationContent(
			pop.Migration{
				Path:      filePath,
				Version:   "",
				Name:      entry.Name(),
				Direction: "up",
				Type:      "sql",
				DBType:    "postgres",
				Runner: func(pop.Migration, *pop.Connection) error {
					fmt.Println("running?!")
					// TODO: ...
					return nil
				},
			},
			conn,
			file,
			true,
		)
		if err != nil {
			return nil, err
		}

		// cache the content
		contentMap[filePath] = content
	}

	return contentMap, nil
}

func removeFirstDirectory(path string) string {
	parts := strings.SplitN(filepath.Clean(path), string(filepath.Separator), 2)
	if len(parts) > 1 {
		return parts[1]
	}

	return path
}

func GeneratePopTemplatedMigrations(
	sourceFolder string,
	destFolder string,
	options map[string]string,
) error {
	base := filepath.Base(filepath.Clean(sourceFolder))
	fileSystem := os.DirFS(filepath.Dir(filepath.Clean(sourceFolder)))
	contentMap, err := MigrationsContent(base, fileSystem, options)
	if err != nil {
		return err
	}

	// save on the destination now
	for filePath, content := range contentMap {
		newPath := path.Join(filepath.Clean(destFolder), removeFirstDirectory(filePath))

		// ensure the directories exist for the new path
		dir := filepath.Dir(newPath)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directories for path %s: %v", newPath, err)
		}

		// write the content to the new path
		err = os.WriteFile(newPath, []byte(content), 0644)
		if err != nil {
			return fmt.Errorf("error writing file %s: %v", newPath, err)
		}
	}

	return nil
}
