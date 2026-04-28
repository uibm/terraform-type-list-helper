package typelist

import (
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TypeListHelper provides functions to manage TypeList attributes.
type TypeListHelper struct {
	FieldName string // Name of the constant field used for matching
}

// New creates a new TypeListHelper instance.
func New(fieldName string) *TypeListHelper {
	return &TypeListHelper{
		FieldName: fieldName,
	}
}

// DiffSuppressFunc returns a DiffSuppressFunc to be used in a TypeList schema.
func (h *TypeListHelper) DiffSuppressFunc(k string, old, new string, d *schema.ResourceData) bool {
	// If the resource is new, there's no diff to suppress.
	if len(d.Id()) == 0 {
		return false
	}

	// Handle cases where the TypeList is set or not.
	_, oldOK := d.GetChange(k)
	if !oldOK.(bool) {
		return false
	}
	_, newOK := d.GetOk(k)
	if !newOK {
		return false
	}

	// If both old and new are present, compare them deeply.
	toAdd, toUpdate, toRemove := h.CalculateChanges(k, d)

	// Suppress the diff if there are no changes.
	return len(toAdd) == 0 && len(toUpdate) == 0 && len(toRemove) == 0
}

// CalculateChanges determines the changes between old and new TypeList values.
func (h *TypeListHelper) CalculateChanges(k string, d *schema.ResourceData) (toAdd, toUpdate, toRemove []interface{}) {
	old, new := d.GetChange(k)
	oldList := old.([]interface{})
	newList := new.([]interface{})

	oldMap := h.listToMap(oldList)
	newMap := h.listToMap(newList)

	for key, newVal := range newMap {
		if oldVal, ok := oldMap[key]; ok {
			// Element exists in both old and new, check for updates.
			if !reflect.DeepEqual(oldVal, newVal) {
				toUpdate = append(toUpdate, newVal)
			}
		} else {
			// Element is new.
			toAdd = append(toAdd, newVal)
		}
	}

	for key, oldVal := range oldMap {
		if _, ok := newMap[key]; !ok {
			// Element exists in old but not in new, mark for removal.
			toRemove = append(toRemove, oldVal)
		}
	}
	// We need to keep the order as specified in the configuration
	// Thus we sort it with respect to the new list
	sort.SliceStable(toUpdate, func(i, j int) bool {
		return h.getElementIndexInList(toUpdate[i], newList) < h.getElementIndexInList(toUpdate[j], newList)
	})

	sort.SliceStable(toRemove, func(i, j int) bool {
		return h.getElementIndexInList(toRemove[i], newList) < h.getElementIndexInList(toRemove[j], newList)
	})

	sort.SliceStable(toAdd, func(i, j int) bool {
		return h.getElementIndexInList(toAdd[i], newList) < h.getElementIndexInList(toAdd[j], newList)
	})

	return toAdd, toUpdate, toRemove
}

// listToMap converts a TypeList to a map using the constant field as the key.
func (h *TypeListHelper) listToMap(list []interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, item := range list {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue // Skip elements that are not maps.
		}
		if key, ok := itemMap[h.FieldName].(string); ok && key != "" {
			result[key] = item
		}
	}
	return result
}

func (h *TypeListHelper) getElementIndexInList(element interface{}, list []interface{}) int {
	for i, item := range list {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue // Skip elements that are not maps.
		}
		elementMap, ok := element.(map[string]interface{})
		if !ok {
			return -1 // the element to find is not a map
		}
		if itemMap[h.FieldName] == elementMap[h.FieldName] {
			return i
		}
	}
	return -1
}

// ApplyOnlyOnce ensures that a field is applied only during resource creation.
func ApplyOnlyOnce(k, o, n string, d *schema.ResourceData) bool {
	// For new resources, allow the first value to be set
	if len(d.Id()) == 0 {
		return false
	}

	// For existing resources, don't allow changes (keep the original value)
	if o == "" {
		return false
	}
	return true
}
