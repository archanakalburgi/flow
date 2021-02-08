package flow

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	pf "github.com/estuary/flow/go/protocols/flow"
	_ "github.com/mattn/go-sqlite3" // Import for registration side-effect.
)

// Catalog is a Catalog database which has been fetched to a local file.
type Catalog struct {
	dbPath string
	db     *sql.DB
}

// NewCatalog copies the catalog DB at the given |url| to the local |tempDir|, and opens it.
func NewCatalog(pathOrURL, tempDir string) (*Catalog, error) {
	if tempDir == "" {
		tempDir = os.TempDir()
	}

	var path = pathOrURL
	if url, err := url.Parse(pathOrURL); err == nil && url.Scheme != "" {
		if path, err = fetchRemote(url, tempDir); err != nil {
			return nil, fmt.Errorf("fetching remote catalog: %w", err)
		}
	}
	db, err := sql.Open("sqlite3", "file:"+path+"?immutable=true&mode=ro")
	if err != nil {
		return nil, fmt.Errorf("opening catalog database %v: %w", path, err)
	}

	return &Catalog{
		dbPath: path,
		db:     db,
	}, nil
}

func fetchRemote(url *url.URL, tempDir string) (string, error) {
	// Download catalog from URL to a local file.
	var dbFile, err = ioutil.TempFile(tempDir, "catalog-db")
	if err != nil {
		return "", fmt.Errorf("failed to create local DB tempfile: %w", err)
	}

	if resp, err := http.Get(url.String()); err != nil {
		return "", fmt.Errorf("failed to request catalog URL: %w", err)
	} else if _, err = io.Copy(dbFile, resp.Body); err != nil {
		return "", fmt.Errorf("failed to copy catalog to local file: %w", err)
	} else if err = resp.Body.Close(); err != nil {
		return "", fmt.Errorf("failed to close catalog response: %w", err)
	} else if err = dbFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close local catalog file: %w", err)
	}
	return dbFile.Name(), nil
}

// LocalPath returns the local path of the catalog.
func (catalog *Catalog) LocalPath() string { return catalog.dbPath }

// Close the Catalog database, rendering it unusable.
func (catalog *Catalog) Close() error {
	return catalog.db.Close()
}

// LoadCollection loads the collection with the given name from the catalog, or returns an error if
// one is not found
func (catalog *Catalog) LoadCollection(name string) (*pf.CollectionSpec, error) {
	var row = catalog.db.QueryRow(`
		SELECT spec FROM built_collections WHERE collection = ?;
		`, name)

	var b []byte
	var collection = new(pf.CollectionSpec)

	if err := row.Scan(&b); err != nil {
		return nil, fmt.Errorf("failed to load collection: %w", err)
	} else if err = collection.Unmarshal(b); err != nil {
		return nil, fmt.Errorf("failed to unmarshal collection: %w", err)
	} else if err = collection.Validate(); err != nil {
		return nil, fmt.Errorf("collection %q is invalid: %w", name, err)
	}
	return collection, nil
}

// LoadCapturedCollections loads all captured collections from the catalog.
func (catalog *Catalog) LoadCapturedCollections() (map[pf.Collection]*pf.CollectionSpec, error) {
	var rows, err = catalog.db.Query(`
		SELECT DISTINCT collection FROM captures WHERE allow_push;
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to read captured collections: %w", err)
	}

	var out = make(map[pf.Collection]*pf.CollectionSpec)

	defer rows.Close()
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan collection name: %w", err)
		} else if collection, err := catalog.LoadCollection(name); err != nil {
			return nil, err
		} else if err = collection.Validate(); err != nil {
			return nil, fmt.Errorf("collection %q is invalid: %w", name, err)
		} else {
			out[collection.Collection] = collection
		}
	}
	return out, nil
}

// LoadDerivedCollection loads the named derived collection from the catalog.
func (catalog *Catalog) LoadDerivedCollection(name string) (*pf.DerivationSpec, error) {
	var row = catalog.db.QueryRow(`
		SELECT d.spec, c.spec
			FROM built_collections AS c
			JOIN built_derivations AS d
			ON c.collection = d.derivation
			WHERE d.derivation = ?;
		`, name)

	var b1, b2 []byte
	if err := row.Scan(&b1, &b2); err != nil {
		return nil, fmt.Errorf("failed to load derivation: %w", err)
	}
	var derivation = new(pf.DerivationSpec)
	if err := derivation.Unmarshal(b1); err != nil {
		return nil, fmt.Errorf("failed to unmarshal derivation: %w", err)
	}
	derivation.Collection = new(pf.CollectionSpec)
	if err := derivation.Collection.Unmarshal(b2); err != nil {
		return nil, fmt.Errorf("failed to unmarshal collection: %w", err)
	}

	var rows, err = catalog.db.Query(`
		SELECT spec FROM built_transforms
		WHERE derivation = ?
		ORDER BY transform asc;
	`, name)

	if err != nil {
		return nil, fmt.Errorf("failed to load transforms: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var transform pf.TransformSpec

		if err = rows.Scan(&b1); err != nil {
			return nil, fmt.Errorf("failed to load transform: %w", err)
		} else if err = transform.Unmarshal(b1); err != nil {
			return nil, fmt.Errorf("failed to unmarshal transform: %w", err)
		}
		derivation.Transforms = append(derivation.Transforms, transform)
	}

	if err = derivation.Validate(); err != nil {
		return nil, fmt.Errorf("derivation %q is invalid: %w", name, err)
	}

	return derivation, nil
}

// LoadJournalRules loads the set of journal rules from the catalog.
func (catalog *Catalog) LoadJournalRules() (*pf.JournalRules, error) {
	var rows, err = catalog.db.Query(`SELECT spec FROM journal_rules ORDER BY rule ASC;`)
	if err != nil {
		return nil, fmt.Errorf("failed to query rules: %w", err)
	}
	var rules = new(pf.JournalRules)

	defer rows.Close()
	for rows.Next() {
		var b []byte
		var rule pf.JournalRules_Rule

		if err = rows.Scan(&b); err != nil {
			return nil, fmt.Errorf("failed to load rule: %w", err)
		} else if err = rule.Unmarshal(b); err != nil {
			return nil, fmt.Errorf("failed to unmarshal rule: %w", err)
		}
		rules.Rules = append(rules.Rules, rule)
	}

	if err = rules.Validate(); err != nil {
		return nil, fmt.Errorf("rules are invalid: %w", err)
	}
	return rules, nil
}

// LoadSchemaBundle loads the bundle of JSON schemas from the catalog.
func (catalog *Catalog) LoadSchemaBundle() (*pf.SchemaBundle, error) {
	var rows, err = catalog.db.Query(`SELECT schema, dom FROM schema_docs;`)
	if err != nil {
		return nil, fmt.Errorf("failed to query rules: %w", err)
	}
	var bundle = &pf.SchemaBundle{
		Bundle: make(map[string]string),
	}

	defer rows.Close()
	for rows.Next() {
		var url, dom string

		if err = rows.Scan(&url, &dom); err != nil {
			return nil, fmt.Errorf("failed to load schema document: %w", err)
		}
		bundle.Bundle[url] = dom
	}
	return bundle, nil
}

// LoadNPMPackage loads the NPM package from a catalog.
func (catalog *Catalog) LoadNPMPackage() ([]byte, error) {
	var row = catalog.db.QueryRow(`SELECT content FROM resources WHERE content_type = '"NpmPack"';`)
	var b []byte

	if err := row.Scan(&b); err != nil {
		return nil, fmt.Errorf("failed to query NPM package: %w", err)
	}
	return b, nil
}
