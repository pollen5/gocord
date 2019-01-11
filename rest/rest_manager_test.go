package rest

import (
	"net/http"
	"testing"
)

func TestParseRoute(t *testing.T) {
	t.Run("normal routes", func(t *testing.T) {
		sample := "/channels/372539957824323584/messages/532935925194555392"
		expected := "/channels/372539957824323584/messages/:id"
		if ParseRoute(http.MethodGet, sample) != expected {
			t.Errorf("Test failed, expected: %s", expected)
		}
	})

	t.Run("individual routes", func(t *testing.T) {
		sample := "/channels/372539957824323584"
		if ParseRoute(http.MethodPut, sample) != sample {
			t.Errorf("Test failed, expected: %s", sample)
		}
	})

	t.Run("bulk deletes", func(t *testing.T) {
		sample := "/channels/372539957824323584/messages/532935925194555392"
		expected := "DELETE /channels/372539957824323584/messages/:id"
		if ParseRoute(http.MethodDelete, sample) != expected {
			t.Errorf("Test failed, expected: %s", expected)
		}
	})
}
