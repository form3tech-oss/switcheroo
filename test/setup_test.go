package test

import (
	"os"
	"testing"
)

var sharedContext *testContext

func TestMain(m *testing.M) {
	given := newTestContext()
	given.
		a_local_instance_of_kubernetes().
		and().
		the_webhook_api_ready_to_receive_mutations()
	sharedContext = given
	code := m.Run()
	sharedContext.tearDown()
	os.Exit(code)
}
