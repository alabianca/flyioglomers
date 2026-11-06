package snowflake_test

import (
	"errors"
	"glomers/internal/snowflake"
	"testing"
	"time"
)

func TestID(t *testing.T) {
	gen := snowflake.Generator{}
	date := time.Date(2025, time.February, 10, 13, 34, 39, 256000000, time.UTC)
	id := gen.MustID(date, 360)
	if !id.Timestamp().Equal(date) {
		t.Fatalf("expected timestamp to equal %s, but got %s\n", date, id.Timestamp())
	}
	if id.NodeID() != 360 {
		t.Fatalf("Expected node id to be %d, but got %d\n", 360, id.NodeID())
	}
}

func TestIDsSameTimestamps(t *testing.T) {
	gen := snowflake.Generator{}
	date := time.Date(2025, time.February, 10, 13, 34, 39, 256000000, time.UTC)
	id := gen.MustID(date, 360)
	id2 := gen.MustID(date, 360)
	if id2 == id {
		t.Fatal("ids are equal, but expected them to be different")
	}
	if id.SequenceId() != 1 {
		t.Fatalf("Expected id to have sequence id %d, but got %d\n", 1, id.SequenceId())
	}
	if id2.SequenceId() != 2 {
		t.Fatalf("Expected id2 to have sequence id %d, but got %d\n", 2, id2.SequenceId())
	}
}

func TestResetSequence(t *testing.T) {
	gen := snowflake.Generator{}
	date := time.Date(2025, time.February, 10, 13, 34, 39, 256000000, time.UTC)
	date2 := time.Date(2025, time.February, 10, 13, 34, 39, 257000000, time.UTC) // 1ms later

	id := gen.MustID(date, 360)
	id2 := gen.MustID(date2, 360)
	if id2 == id {
		t.Fatal("ids are equal, but expected them to be different")
	}
	if id.SequenceId() != 1 {
		t.Fatalf("Expected id to have sequence id %d, but got %d\n", 1, id.SequenceId())
	}
	if id2.SequenceId() != 1 {
		t.Fatalf("Expected id2 to have sequence id %d, but got %d\n", 1, id2.SequenceId())
	}
}

func TestSequenceOverflow(t *testing.T) {
	gen := snowflake.Generator{}
	ids := make(map[snowflake.ID]bool)
	for i := 0; i < 4095; i++ {
		date := time.Date(2025, time.February, 10, 13, 34, 39, 256000000, time.UTC)
		id := gen.MustID(date, 360)
		if id.SequenceId() != int64(i+1) {
			t.Fatalf("Expected sequence to be %d, but got %d\n", i+1, id.SequenceId())
		}
		if _, ok := ids[id]; ok {
			t.Fatalf("Got duplicate id %d\n", id)
		}
		ids[id] = true
	}

	date := time.Date(2025, time.February, 10, 13, 34, 39, 256000000, time.UTC)
	id := gen.MustID(date, 360)
	if _, ok := ids[id]; ok {
		t.Fatalf("Got duplicate id %d\n", id)
	}
	ids[id] = true
	if id.SequenceId() != int64(1) {
		t.Fatalf("Expected sequence to reset to 1, but sequence is %d\n", id.SequenceId())
	}
	if id.Timestamp().Equal(date) {
		t.Fatalf("Expected timestamp to correct after sequence overflow, but it equals provided timestamp (%s)\n", id.Timestamp())
	}
	if id.Timestamp().Before(date) {
		t.Fatalf("Expected timestamp to be larger than provided date, but got %s", id.Timestamp())
	}

	id = gen.MustID(date, 360)
	if _, ok := ids[id]; ok {
		t.Fatalf("Got duplicate id %d\n", id)
	}
	ids[id] = true
	if id.SequenceId() != int64(2) {
		t.Fatalf("Expected sequence to reset to 1, but sequence is %d\n", id.SequenceId())
	}
	if id.Timestamp().Equal(date) {
		t.Fatalf("Expected timestamp to correct after sequence overflow, but it equals provided timestamp (%s)\n", id.Timestamp())
	}
	if id.Timestamp().Before(date) {
		t.Fatalf("Expected timestamp to be larger than provided date, but got %s", id.Timestamp())
	}
}

func TestNodeIdMustBePositive(t *testing.T) {
	gen := snowflake.Generator{}
	date := time.Date(2025, time.February, 10, 13, 34, 39, 256000000, time.UTC)
	_, err := gen.NewID(date, -1)
	if err == nil {
		t.Fatalf("Expected to get an error, but got nil")
	}
	if !errors.Is(err, snowflake.ErrNodeIdOverflow) {
		t.Fatalf("Expected error: %s, but got %s", snowflake.ErrNodeIdOverflow, err)
	}
}

func TestNodeIdMustBeLessThan1024(t *testing.T) {
	gen := snowflake.Generator{}
	date := time.Date(2025, time.February, 10, 13, 34, 39, 256000000, time.UTC)
	_, err := gen.NewID(date, 1024)
	if err == nil {
		t.Fatalf("Expected to get an error, but got nil")
	}
	if !errors.Is(err, snowflake.ErrNodeIdOverflow) {
		t.Fatalf("Expected error: %s, but got %s", snowflake.ErrNodeIdOverflow, err)
	}
}

func TestDecode(t *testing.T) {
	id := snowflake.ID(2118526850467627000)
	t.Logf("Timestamp %s\n", id.Timestamp())
	t.Logf("Machine ID %d\n", id.NodeID())
	t.Logf("Sequence ID: %d\n", id.SequenceId())
}

func TestString(t *testing.T) {
	id := snowflake.ID(2118526850467627000)
	if "2118526850467627000" != id.String() {
		t.Fatal("string conversion failed")
	}
}
