package admissionserver

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"k8s.io/api/admission/v1beta1"
)

func TestMutationHookAdmit(t *testing.T) {
	hook := &mutationHook{imageScannerEvaluator: &mockImageScannerEvaluator{Accepted: true}}

	review := &v1beta1.AdmissionRequest{}
	if b, err := ioutil.ReadFile("./assets/admission-review.json"); err != nil {
		t.Fatal(err)
	} else {
		json.Unmarshal(b, review)
	}

	response := hook.Admit(review)

	if !response.Allowed {
		t.Fatalf("Admission should not be rejected")
	}

	if response.UID != review.UID {
		t.Fatalf("Unexpected UID: %s", response.UID)
	}

	if response.Result != nil {
		t.Fatalf("Response Result should be nil")
	}

	if *response.PatchType != v1beta1.PatchType("JSONPatch") {
		t.Fatalf("Unexpected PatchType: %s", *response.PatchType)
	}

	var a []interface{}
	if err := json.Unmarshal(response.Patch, &a); err != nil {
		t.Fatalf("Unmarshal error for JSON Patch: %s\n%v", string(response.Patch), err)
	}

	if string(response.Patch) != `[{"op": "add", "path": "/metadata/annotations", "value": {}}, {"op": "replace", "path": "/spec/containers/0/image", "value": "image1@digest1"}, {"op": "add", "path": "/metadata/annotations/admission.sysdig.com~1container-1-original-image", "value": "image1:tag1"}, {"op": "add", "path": "/metadata/annotations/admission.sysdig.com~1container-1-mutated-image", "value": "image1@digest1"}, {"op": "replace", "path": "/spec/containers/1/image", "value": "image2@digest2"}, {"op": "add", "path": "/metadata/annotations/admission.sysdig.com~1container-2-original-image", "value": "image2:tag2"}, {"op": "add", "path": "/metadata/annotations/admission.sysdig.com~1container-2-mutated-image", "value": "image2@digest2"}]` {
		t.Fatalf("Unexpected JSON Patch: %s", string(response.Patch))
	}
}

func TestMutationHookPreserveAnnotations(t *testing.T) {
	hook := &mutationHook{imageScannerEvaluator: &mockImageScannerEvaluator{Accepted: true}}

	review := &v1beta1.AdmissionRequest{}
	if b, err := ioutil.ReadFile("./assets/admission-review-with-annotations.json"); err != nil {
		t.Fatal(err)
	} else {
		json.Unmarshal(b, review)
	}

	response := hook.Admit(review)

	if !response.Allowed {
		t.Fatalf("Admission should not be rejected")
	}

	if response.UID != review.UID {
		t.Fatalf("Unexpected UID: %s", response.UID)
	}

	if response.Result != nil {
		t.Fatalf("Response Result should be nil")
	}

	if *response.PatchType != v1beta1.PatchType("JSONPatch") {
		t.Fatalf("Unexpected PatchType: %s", *response.PatchType)
	}

	var a []interface{}
	if err := json.Unmarshal(response.Patch, &a); err != nil {
		t.Fatalf("Unmarshal error for JSON Patch: %s\n%v", string(response.Patch), err)
	}

	if string(response.Patch) != `[{"op": "replace", "path": "/spec/containers/0/image", "value": "image1@digest1"}, {"op": "add", "path": "/metadata/annotations/admission.sysdig.com~1container-1-original-image", "value": "image1:tag1"}, {"op": "add", "path": "/metadata/annotations/admission.sysdig.com~1container-1-mutated-image", "value": "image1@digest1"}, {"op": "replace", "path": "/spec/containers/1/image", "value": "image2@digest2"}, {"op": "add", "path": "/metadata/annotations/admission.sysdig.com~1container-2-original-image", "value": "image2:tag2"}, {"op": "add", "path": "/metadata/annotations/admission.sysdig.com~1container-2-mutated-image", "value": "image2@digest2"}]` {
		t.Fatalf("Unexpected JSON Patch: %s", string(response.Patch))
	}
}
