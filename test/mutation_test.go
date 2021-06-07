package test

import (
	"testing"
)

func TestPodCreateMutation(t *testing.T) {
	t.Run("", func(t *testing.T) {
		given, when, then := NewMutationScenario(t, sharedContext.client)

		given.preparedPod(testPodWithName("create-test-pod"))

		when.podCreated()

		then.podImagesHasRequiredRegistry()
	})
}

func TestPodUpdateMutation(t *testing.T) {
	t.Run("", func(t *testing.T) {
		given, when, then := NewMutationScenario(t, sharedContext.client)

		given.preparedPod(testPodWithName("update-test-pod")).
			and().podCreated().
			and().prepareUpdatedPod(updateImageOfPod)

		when.podUpdated()

		then.podImagesHasRequiredRegistry()
	})
}
