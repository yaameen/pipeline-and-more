package pipeline

import "testing"

func Test_String_Pipeline(t *testing.T) {

	pipe := NewPipeline(nil)

	if pipe.Return() != nil {
		t.Fail()
	}

	pipe.Send("testing")

	if pipe.Return() != "testing" {
		t.Fail()
	}

	pipe.Through(func(value interface{}, next Handler) {
		next(value.(string)+" amigo", nil)
	})

	if pipe.Return() != "testing amig" {
		t.Fatalf("Expected testing but returned %v", pipe.value)
	}

}

// TODO: add more tests
