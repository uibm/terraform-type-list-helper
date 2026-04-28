package typelist

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDiffSuppressFunc(t *testing.T) {
	helper := New("id")
	resourceData := &schema.ResourceData{} // Mock or use a proper mock library

	// Test case 1: New resource, should not suppress
	resourceData.MarkNewResource()
	if helper.DiffSuppressFunc("list", "", "", resourceData) {
		t.Errorf("Diff should not be suppressed for new resources")
	}

	// Test case 2: Resource exists, no changes, should suppress
	resourceData = schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"list": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id":   {Type: schema.TypeString},
					"name": {Type: schema.TypeString},
				},
			},
		},
	}, map[string]interface{}{
		"list": []interface{}{
			map[string]interface{}{"id": "1", "name": "foo"},
			map[string]interface{}{"id": "2", "name": "bar"},
		},
	})

	if !helper.DiffSuppressFunc("list", "", "", resourceData) {
		t.Errorf("Diff should be suppressed when there are no changes")
	}

	// Test case 3: Resource exists, changes in list, should not suppress
	resourceData = schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"list": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id":   {Type: schema.TypeString},
					"name": {Type: schema.TypeString},
				},
			},
		},
	}, map[string]interface{}{
		"list": []interface{}{
			map[string]interface{}{"id": "1", "name": "foo"},
			map[string]interface{}{"id": "2", "name": "bar"},
		},
	})
	resourceData.Set("list", []interface{}{
		map[string]interface{}{"id": "1", "name": "updated"},
		map[string]interface{}{"id": "2", "name": "bar"},
	})

	if helper.DiffSuppressFunc("list", "", "", resourceData) {
		t.Errorf("Diff should not be suppressed when there are changes in the list")
	}

}

func TestCalculateChanges(t *testing.T) {
	helper := New("id")

	tests := []struct {
		name     string
		old      []interface{}
		new      []interface{}
		toAdd    []interface{}
		toUpdate []interface{}
		toRemove []interface{}
	}{
		{
			name: "no changes",
			old: []interface{}{
				map[string]interface{}{"id": "1", "name": "foo"},
				map[string]interface{}{"id": "2", "name": "bar"},
			},
			new: []interface{}{
				map[string]interface{}{"id": "1", "name": "foo"},
				map[string]interface{}{"id": "2", "name": "bar"},
			},
			toAdd:    []interface{}{},
			toUpdate: []interface{}{},
			toRemove: []interface{}{},
		},
		{
			name: "add element",
			old: []interface{}{
				map[string]interface{}{"id": "1", "name": "foo"},
			},
			new: []interface{}{
				map[string]interface{}{"id": "1", "name": "foo"},
				map[string]interface{}{"id": "2", "name": "bar"},
			},
			toAdd:    []interface{}{map[string]interface{}{"id": "2", "name": "bar"}},
			toUpdate: []interface{}{},
			toRemove: []interface{}{},
		},
		{
			name: "remove element",
			old: []interface{}{
				map[string]interface{}{"id": "1", "name": "foo"},
				map[string]interface{}{"id": "2", "name": "bar"},
			},
			new: []interface{}{
				map[string]interface{}{"id": "1", "name": "foo"},
			},
			toAdd:    []interface{}{},
			toUpdate: []interface{}{},
			toRemove: []interface{}{map[string]interface{}{"id": "2", "name": "bar"}},
		},
		{
			name: "update element",
			old: []interface{}{
				map[string]interface{}{"id": "1", "name": "foo"},
				map[string]interface{}{"id": "2", "name": "bar"},
			},
			new: []interface{}{
				map[string]interface{}{"id": "1", "name": "updated"},
				map[string]interface{}{"id": "2", "name": "bar"},
			},
			toAdd:    []interface{}{},
			toUpdate: []interface{}{map[string]interface{}{"id": "1", "name": "updated"}},
			toRemove: []interface{}{},
		},
		{
			name: "reorder elements",
			old: []interface{}{
				map[string]interface{}{"id": "1", "name": "foo"},
				map[string]interface{}{"id": "2", "name": "bar"},
			},
			new: []interface{}{
				map[string]interface{}{"id": "2", "name": "bar"},
				map[string]interface{}{"id": "1", "name": "foo"},
			},
			toAdd:    []interface{}{},
			toUpdate: []interface{}{},
			toRemove: []interface{}{},
		},
		{
			name: "complex changes",
			old: []interface{}{
				map[string]interface{}{"id": "1", "name": "foo"},
				map[string]interface{}{"id": "2", "name": "bar"},
				map[string]interface{}{"id": "3", "name": "baz"},
			},
			new: []interface{}{
				map[string]interface{}{"id": "2", "name": "updated"},
				map[string]interface{}{"id": "4", "name": "new"},
				map[string]interface{}{"id": "1", "name": "foo"},
			},
			toAdd:    []interface{}{map[string]interface{}{"id": "4", "name": "new"}},
			toUpdate: []interface{}{map[string]interface{}{"id": "2", "name": "updated"}},
			toRemove: []interface{}{map[string]interface{}{"id": "3", "name": "baz"}},
		},
		{
			name:     "empty list",
			old:      []interface{}{},
			new:      []interface{}{},
			toAdd:    []interface{}{},
			toUpdate: []interface{}{},
			toRemove: []interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceData := &schema.ResourceData{} // Mock or use a proper mock library
			resourceData = schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"list": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id":   {Type: schema.TypeString},
							"name": {Type: schema.TypeString},
						},
					},
				},
			}, map[string]interface{}{
				"list": tt.old,
			})
			resourceData.Set("list", tt.new)
			resourceData.SetId("id")
			toAdd, toUpdate, toRemove := helper.CalculateChanges("list", resourceData)

			if !reflect.DeepEqual(toAdd, tt.toAdd) {
				t.Errorf("toAdd mismatch: got %v, want %v", toAdd, tt.toAdd)
			}
			if !reflect.DeepEqual(toUpdate, tt.toUpdate) {
				t.Errorf("toUpdate mismatch: got %v, want %v", toUpdate, tt.toUpdate)
			}
			if !reflect.DeepEqual(toRemove, tt.toRemove) {
				t.Errorf("toRemove mismatch: got %v, want %v", toRemove, tt.toRemove)
			}
		})
	}
}

func TestApplyOnlyOnce(t *testing.T) {
	resourceData := &schema.ResourceData{} // Mock or use a proper mock library
	tests := []struct {
		name         string
		resourceId   string
		oldValue     string
		newValue     string
		expectResult bool
	}{
		{
			name:         "new resource",
			resourceId:   "",
			oldValue:     "",
			newValue:     "new",
			expectResult: false,
		},
		{
			name:         "existing resource with empty old value",
			resourceId:   "id",
			oldValue:     "",
			newValue:     "new",
			expectResult: false,
		},
		{
			name:         "existing resource with changes not allowed",
			resourceId:   "id",
			oldValue:     "old",
			newValue:     "new",
			expectResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceData.SetId(tt.resourceId)
			result := ApplyOnlyOnce("key", tt.oldValue, tt.newValue, resourceData)

			if result != tt.expectResult {
				t.Errorf("ApplyOnlyOnce test failed for case: %s, expected: %v, got: %v", tt.name, tt.expectResult, result)
			}
		})
	}
}
