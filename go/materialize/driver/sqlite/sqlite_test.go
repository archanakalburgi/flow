package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/estuary/flow/go/bindings"
	"github.com/estuary/flow/go/fdb/tuple"
	sqlDriver "github.com/estuary/flow/go/materialize/driver/sql2"
	pf "github.com/estuary/flow/go/protocols/flow"
	pm "github.com/estuary/flow/go/protocols/materialize"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func TestSQLiteDriver(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	var ctx = context.Background()
	const bufSize = 1024 * 1024

	var lis *bufconn.Listener
	var bufDialer = func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pm.RegisterDriverServer(s, NewSQLiteDriver())
	var done = make(chan error, 1)
	go func() {
		var e = s.Serve(lis)
		done <- e
	}()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	var client = pm.NewDriverClient(conn)

	doTestSQLite(t, client)
	s.GracefulStop()
	err = <-done
	require.NoError(t, err)
}

func doTestSQLite(t *testing.T, driver pm.DriverClient) {
	var built, err = bindings.BuildCatalog(bindings.BuildArgs{
		FileRoot: "./testdata",
		BuildAPI_Config: pf.BuildAPI_Config{
			Directory:   "testdata",
			Source:      "file:///flow.yaml",
			SourceType:  pf.ContentType_CATALOG_SPEC,
			CatalogPath: filepath.Join(t.TempDir(), "catalog.db"),
		}})
	require.NoError(t, err)
	require.Empty(t, built.Errors)

	// Config fixture which matches schema of ParseConfig.
	var spec = struct {
		Path  string
		Table string
	}{
		Path:  "file://" + path.Join(t.TempDir(), "target.db"),
		Table: "test_target",
	}
	var specJSON, _ = json.Marshal(spec)

	collection := &built.Collections[0]

	// Validate should return constraints for a non-existant materialization
	var validateReq = pm.ValidateRequest{
		EndpointName:     "an/endpoint",
		EndpointType:     pf.EndpointType_SQLITE,
		EndpointSpecJson: string(specJSON),
		Collection:       collection,
	}

	var ctx = context.Background()
	validateResp, err := driver.Validate(ctx, &validateReq)
	require.NoError(t, err)
	// There should be a constraint for every projection
	require.Equal(t, &pm.ValidateResponse{
		Constraints: map[string]*pm.Constraint{
			"array":         {Type: pm.Constraint_FIELD_OPTIONAL, Reason: "This field is able to be materialized"},
			"bool":          {Type: pm.Constraint_LOCATION_RECOMMENDED, Reason: "The projection has a single scalar type"},
			"flow_document": {Type: pm.Constraint_LOCATION_REQUIRED, Reason: "The root document must be materialized"},
			"int":           {Type: pm.Constraint_LOCATION_RECOMMENDED, Reason: "The projection has a single scalar type"},
			"number":        {Type: pm.Constraint_LOCATION_RECOMMENDED, Reason: "The projection has a single scalar type"},
			"object":        {Type: pm.Constraint_FIELD_OPTIONAL, Reason: "This field is able to be materialized"},
			"string":        {Type: pm.Constraint_LOCATION_RECOMMENDED, Reason: "The projection has a single scalar type"},
			"theKey":        {Type: pm.Constraint_LOCATION_REQUIRED, Reason: "All Locations that are part of the collections key are required"},
		},
		ResourcePath: []string{"test_target"},
	}, validateResp)

	// Select some fields and Apply the materialization
	var fields = pf.FieldSelection{
		Keys:     []string{"theKey"},
		Values:   []string{"bool", "int", "string"}, // intentionally missing "number" field
		Document: "flow_document",
	}
	var applyReq = pm.ApplyRequest{
		Materialization: &pf.MaterializationSpec{
			Materialization:      "a/materialization",
			Collection:           *collection,
			EndpointName:         "an/endpoint",
			EndpointType:         pf.EndpointType_SQLITE,
			EndpointSpecJson:     string(specJSON),
			EndpointResourcePath: []string{"test_target"},
			FieldSelection:       fields,
		},
		DryRun: true,
	}

	applyResp, err := driver.Apply(ctx, &applyReq)
	require.NoError(t, err)
	require.NotEmpty(t, applyResp.ActionDescription)

	applyReq.DryRun = false
	applyResp, err = driver.Apply(ctx, &applyReq)
	require.NoError(t, err)
	require.NotEmpty(t, applyResp.ActionDescription)

	// Now that we've applied, call Validate again to ensure the existing fields are accounted for
	validateResp, err = driver.Validate(ctx, &validateReq)
	require.NoError(t, err)
	require.Equal(t, len(collection.Projections), len(validateResp.Constraints))
	for _, field := range fields.AllFields() {
		var actual = validateResp.Constraints[field].Type
		require.Equal(
			t,
			pm.Constraint_FIELD_REQUIRED,
			actual,
			"wrong constraint for field: %s, expected FIELD_REQUIRED, got %s",
			field,
			actual,
		)
	}
	// The "number" field should be forbidden because it was not included in the FieldSelection that
	// was applied.
	require.Equal(t, pm.Constraint_FIELD_FORBIDDEN, validateResp.Constraints["number"].Type)

	transaction, err := driver.Transactions(ctx)
	require.NoError(t, err)

	// Send open.
	err = transaction.Send(&pm.TransactionRequest{
		Open: &pm.TransactionRequest_Open{
			Materialization:  applyReq.Materialization,
			DriverCheckpoint: nil,
			ShardFqn:         "canary",
		},
	})
	require.NoError(t, err)

	// Receive Opened.
	opened, err := transaction.Recv()
	require.NoError(t, err)
	require.Equal(t, &pm.TransactionResponse_Opened{
		DeltaUpdates:   false,
		FlowCheckpoint: pm.ExplicitZeroCheckpoint,
	}, opened.Opened)

	// Test Load with keys that don't exist yet
	var key1 = tuple.Tuple{"key1Value"}
	var key2 = tuple.Tuple{"key2Value"}
	var loadReq = newLoadReq(key1.Pack(), key2.Pack())
	err = transaction.Send(&pm.TransactionRequest{
		Load: &loadReq,
	})
	require.NoError(t, err)
	var key3 = tuple.Tuple{"key3Value"}
	loadReq = newLoadReq(key3.Pack())
	err = transaction.Send(&pm.TransactionRequest{
		Load: &loadReq,
	})
	require.NoError(t, err)

	var checkpoint1 = []byte("first checkpoint value")
	err = transaction.Send(&pm.TransactionRequest{
		Prepare: &pm.TransactionRequest_Prepare{
			FlowCheckpoint: checkpoint1,
		},
	})
	require.NoError(t, err)

	// Receive Prepared, which indicates that none of the documents exist
	prepared, err := transaction.Recv()
	require.NoError(t, err)
	require.NotNil(t, prepared.Prepared, "unexpected message: %v+", prepared)
	require.Empty(t, prepared.Prepared.DriverCheckpoint)

	// Test Store to add those keys
	require.NoError(t, err)
	var doc1 = `{ "theKey": "key1Value", "string": "foo", "bool": true, "int": 77, "number": 12.34 }`
	var doc2 = `{ "theKey": "key2Value", "string": "bar", "bool": false, "int": 88, "number": 56.78 }`
	var doc3 = `{ "theKey": "key3Value", "string": "baz", "bool": false, "int": 99, "number": 0 }`

	var store1 = pm.TransactionRequest_Store{}
	store1.DocsJson = store1.Arena.AddAll([]byte(doc1), []byte(doc2))
	store1.PackedKeys = store1.Arena.AddAll(key1.Pack(), key2.Pack())
	store1.PackedValues = store1.Arena.AddAll(
		tuple.Tuple{"foo", true, 77}.Pack(),
		tuple.Tuple{"bar", false, 88}.Pack(),
	)
	store1.Exists = []bool{false, false}
	err = transaction.Send(&pm.TransactionRequest{
		Store: &store1,
	})
	require.NoError(t, err)

	var store2 = pm.TransactionRequest_Store{}
	store2.DocsJson = store2.Arena.AddAll([]byte(doc3))
	store2.PackedKeys = store2.Arena.AddAll(key3.Pack())
	store2.PackedValues = store2.Arena.AddAll(
		tuple.Tuple{"baz", false, 99}.Pack(),
	)
	store2.Exists = []bool{false}
	err = transaction.Send(&pm.TransactionRequest{
		Store: &store2,
	})
	require.NoError(t, err)

	err = transaction.Send(&pm.TransactionRequest{
		Commit: &pm.TransactionRequest_Commit{},
	})
	require.NoError(t, err)

	committed, err := transaction.Recv()
	require.NoError(t, err)
	require.NotNil(t, committed.Committed)

	// Next transaction.

	loadReq = newLoadReq(key1.Pack(), key2.Pack(), key3.Pack())
	err = transaction.Send(&pm.TransactionRequest{
		Load: &loadReq,
	})
	require.NoError(t, err)

	var checkpoint2 = []byte("second checkpoint value")
	err = transaction.Send(&pm.TransactionRequest{
		Prepare: &pm.TransactionRequest_Prepare{
			FlowCheckpoint: checkpoint2,
		},
	})
	require.NoError(t, err)

	// Receive Loaded response, which is expected to contain our 3 documents.
	loaded, err := transaction.Recv()
	require.NoError(t, err)
	require.NotNil(t, loaded.Loaded)
	require.Equal(t, 3, len(loaded.Loaded.DocsJson))

	for i, expected := range []string{doc1, doc2, doc3} {
		var actual = loaded.Loaded.Arena.Bytes(loaded.Loaded.DocsJson[i])
		require.Equal(t, expected, string(actual))
	}

	// Receive Prepared
	prepared, err = transaction.Recv()
	require.NoError(t, err)
	require.NotNil(t, prepared.Prepared, "unexpected message: %v+", prepared)

	// This store will update one document and add a new one.
	var newDoc1 = `{ "theKey": "key1Value", "string": "notthesame", "bool": false, "int": 33, "number": 2 }`
	var key4 = tuple.Tuple{"key4Value"}
	var doc4 = `{ "theKey": "key4Value" }`

	var storeReq = pm.TransactionRequest_Store{}
	storeReq.Exists = []bool{true, false}
	storeReq.PackedKeys = storeReq.Arena.AddAll(key1.Pack(), key4.Pack())
	storeReq.PackedValues = storeReq.Arena.AddAll(
		tuple.Tuple{"totally different", false, 33}.Pack(),
		tuple.Tuple{nil, nil, nil}.Pack(),
	)
	storeReq.DocsJson = storeReq.Arena.AddAll([]byte(newDoc1), []byte(doc4))

	err = transaction.Send(&pm.TransactionRequest{
		Store: &storeReq,
	})
	require.NoError(t, err)

	// Commit transaction and assert we get a Committed.
	err = transaction.Send(&pm.TransactionRequest{
		Commit: &pm.TransactionRequest_Commit{},
	})
	require.NoError(t, err)

	committed, err = transaction.Recv()
	require.NoError(t, err)
	require.NotNil(t, committed.Committed)

	// One more transaction just to verify the updated documents

	loadReq = newLoadReq(key1.Pack(), key2.Pack(), key3.Pack(), key4.Pack())
	err = transaction.Send(&pm.TransactionRequest{
		Load: &loadReq,
	})
	require.NoError(t, err)

	var checkpoint3 = []byte("third checkpoint value")
	err = transaction.Send(&pm.TransactionRequest{
		Prepare: &pm.TransactionRequest_Prepare{
			FlowCheckpoint: checkpoint3,
		},
	})
	require.NoError(t, err)

	// Receive LoadResponse, which is expected to contain 4 documents.

	loaded, err = transaction.Recv()
	require.NoError(t, err)
	require.NotNil(t, loaded.Loaded)
	require.Equal(t, 4, len(loaded.Loaded.DocsJson))

	for i, expected := range []string{newDoc1, doc2, doc3, doc4} {
		var actual = loaded.Loaded.Arena.Bytes(loaded.Loaded.DocsJson[i])
		require.Equal(t, expected, string(actual))
	}

	// Receive Prepared
	prepared, err = transaction.Recv()
	require.NoError(t, err)
	require.NotNil(t, prepared.Prepared, "unexpected message: %v+", prepared)

	// Send and receive Commit / Committed.
	err = transaction.Send(&pm.TransactionRequest{
		Commit: &pm.TransactionRequest_Commit{},
	})
	require.NoError(t, err)

	committed, err = transaction.Recv()
	require.NoError(t, err)
	require.NotNil(t, committed.Committed)

	// Shut down stream.
	require.NoError(t, transaction.CloseSend())
	_, err = transaction.Recv()
	require.Equal(t, io.EOF, err)

	// Last thing is to snapshot the database tables we care about
	var quotes = sqlDriver.DoubleQuotes()
	var tab = sqlDriver.TableForMaterialization(spec.Table, "", &quotes,
		&pf.MaterializationSpec{
			Collection:     *collection,
			FieldSelection: fields,
		})
	var dump = dumpTables(t, spec.Path, tab)
	cupaloy.SnapshotT(t, dump)
}

type AnyCol string

func (col *AnyCol) Scan(i interface{}) error {
	var sval string
	if b, ok := i.([]byte); ok {
		sval = string(b)
	} else {
		sval = fmt.Sprint(i)
	}
	*col = AnyCol(sval)
	return nil
}
func (col AnyCol) String() string {
	return string(col)
}

func dumpTables(t *testing.T, uri string, tables ...*sqlDriver.Table) string {
	uri = fmt.Sprintf("%s?mode=ro", uri)
	var db, err = sql.Open("sqlite3", uri)
	require.NoError(t, err)
	defer db.Close()

	var builder strings.Builder
	for tn, table := range tables {
		if tn > 0 {
			builder.WriteString("\n\n") // make it more readable
		}
		var colNames strings.Builder
		for i, col := range table.Columns {
			if i > 0 {
				colNames.WriteString(", ")
			}
			colNames.WriteString(col.Identifier)
		}

		var sql = fmt.Sprintf("SELECT %s FROM %s;", colNames.String(), table.Identifier)
		rows, err := db.Query(sql)
		require.NoError(t, err)
		defer rows.Close()

		fmt.Fprintf(&builder, "%s:\n", table.Identifier)
		builder.WriteString(colNames.String())

		for rows.Next() {
			var data = make([]AnyCol, len(table.Columns))
			var ptrs = make([]interface{}, len(table.Columns))
			for i := range data {
				ptrs[i] = &data[i]
			}
			err = rows.Scan(ptrs...)
			require.NoError(t, err)
			builder.WriteString("\n")
			for i, v := range ptrs {
				if i > 0 {
					builder.WriteString(", ")
				}
				var val = v.(*AnyCol)
				builder.WriteString(val.String())
			}
		}
	}
	return builder.String()
}

func newLoadReq(keys ...[]byte) pm.TransactionRequest_Load {
	var arena pf.Arena
	var packedKeys = arena.AddAll(keys...)
	return pm.TransactionRequest_Load{
		Arena:      arena,
		PackedKeys: packedKeys,
	}
}
