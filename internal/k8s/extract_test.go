package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/astro-walker/tilt/internal/k8s/testyaml"
)

func TestExtractSanchoContainers(t *testing.T) {
	entities, err := ParseYAMLFromString(testyaml.SanchoYAML)
	if err != nil {
		t.Fatal(err)
	}

	if len(entities) != 1 {
		t.Fatalf("Unexpected entities: %+v", entities)
	}

	entity := entities[0]
	containers, err := extractContainers(&entity)
	if err != nil {
		t.Fatal(err)
	}

	if len(containers) != 1 || containers[0].Image != "gcr.io/some-project-162817/sancho" {
		t.Errorf("Unexpected containers: %v", containers)
	}
}

func TestExtractSanchoPods(t *testing.T) {
	entities, err := ParseYAMLFromString(testyaml.SanchoYAML)
	if err != nil {
		t.Fatal(err)
	}

	if len(entities) != 1 {
		t.Fatalf("Unexpected entities: %+v", entities)
	}

	entity := entities[0]
	pods, err := ExtractPods(&entity)
	if err != nil {
		t.Fatal(err)
	}

	if len(pods) != 1 || pods[0].Containers[0].Name != "sancho" {
		t.Errorf("Unexpected pods: %v", pods)
	}
}

func TestExtractSanchoPodTemplateSpecs(t *testing.T) {
	entities, err := ParseYAMLFromString(testyaml.SanchoYAML)
	if err != nil {
		t.Fatal(err)
	}

	if len(entities) != 1 {
		t.Fatalf("Unexpected entities: %+v", entities)
	}

	entity := entities[0]
	tempSpecs, err := ExtractPodTemplateSpec(&entity)
	if err != nil {
		t.Fatal(err)
	}

	expectedLabels := map[string]string{"app": "sancho"}
	if assert.Equal(t, len(tempSpecs), 1) {
		assert.Equal(t, tempSpecs[0].ObjectMeta.Labels, expectedLabels)
	}
}
