package mongorepo_test

import (
	"context"
	"github.com/bitstored/repository/pkg/mongorepo"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

type Test struct {
	Name  string
	Field string
	Fails bool
}

func TestCreate(t *testing.T) {
	repo := mongorepo.NewRepository("mongodb://localhost:27017", "test")
	test := Test{"TestCreate", "fieldCreate", false}
	rsp, err := repo.Create(context.Background(), "test_create", test)
	require.NoError(t, err)
	t.Log(rsp)
}

func TestRead(t *testing.T) {
	repo := mongorepo.NewRepository("mongodb://localhost:27017", "test")
	test := bson.M{"Name": "TestRead", "Field": "fieldRead", "Fails": false}
	rsp, err := repo.Create(context.Background(), "test_read", test)
	require.NoError(t, err)
	t.Log(rsp)
	res, err := repo.Read(context.Background(), "test_read", [][]string{{"Name", "TestRead"}})
	require.NoError(t, err)
	res.Next(context.Background())
	var p Test
	res.Decode(&p)
	t.Log(p)

	res, err = repo.Read(context.Background(), "test_read", [][]string{{"Name", "TestRead"}, {"Field", "fieldRead"}})
	require.NoError(t, err)
	t.Log(res)

	res, err = repo.Read(context.Background(), "test_read", [][]string{{"Name", "NotExists"}})
	require.NoError(t, err)
	t.Log(res)
}
