package couchbaserepo_test

import (
	"context"
	"github.com/couchbase/gocb"
	"github.com/dianabejan/repository/pkg/couchbaserepo"
	_ "github.com/dianabejan/repository/pkg/repo"
	"testing"
)

func TestSuite(t *testing.T) {
	// connect to cluster
	cluster, _ := gocb.Connect("http://localhost")
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})
	// open bucket
	bucket, _ := cluster.OpenBucket("bitstored", "")
	r := couchbaserepo.NewRepository(bucket)
	// init object
	beer := make(map[string]interface{})
	beer["comment"] = "Random beer from Norway"
	// write object into the data base
	cas, err := r.Create(context.Background(), "jsonbeer", &beer)
	if err != nil {
		t.Fatal(err.Error())
	}
	// read object
	var value interface{}
	_, err = r.Read(context.Background(), "jsonbeer", &value)
	if err != nil {
		t.Fatal(err.Error())
	}
	if value == nil {
		t.Fatal("no data")
	}
	// update object
	beer["volume"] = "300 ml"
	cas, err = r.Update(context.Background(), "jsonbeer", cas, &beer)
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = r.Read(context.Background(), "jsonbeer", &value)
	if value.(map[string]interface{})["volume"] == nil {
		t.Fatalf("not updated")
	}
	// delete object
	_, err = r.Delete(context.Background(), "jsonbeer", cas)
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = r.Read(context.Background(), "jsonbeer", &value)
	if err == nil {
		t.Fatal(err.Error())
	}
}
