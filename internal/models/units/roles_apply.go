// File: simfight-tactics/internal/models/units/roles_apply.go
package units

import (
	"fmt"
	"reflect"
	"strings"
)

// ApplyReport collects issues encountered while applying a JSON-like map onto Stats.
type ApplyReport struct {
	UnknownKeys []string // keys with no matching JSON tag in the target struct
	TypeErrors  []string // "path.to.key: expected <type>, got <actual>"
}

func (r *ApplyReport) empty() bool {
	return len(r.UnknownKeys) == 0 && len(r.TypeErrors) == 0
}

func (r *ApplyReport) appendUnknown(path string) {
	r.UnknownKeys = append(r.UnknownKeys, path)
}

func (r *ApplyReport) appendTypeErr(path, expected, got string) {
	r.TypeErrors = append(r.TypeErrors, fmt.Sprintf("%s: expected %s, got %s", path, expected, got))
}

// applyRoleMapToStats applies a JSON-like map[string]any onto *Stats using JSON tags,
// collecting unknown keys and type mismatches.
func applyRoleMapToStats(dst *Stats, roleDoc map[string]any) ApplyReport {
	var report ApplyReport
	if dst == nil || len(roleDoc) == 0 {
		return report
	}
	applyJSONMapByTags(reflect.ValueOf(dst).Elem(), roleDoc, &report, "")
	return report
}

// applyJSONMapByTags recursively traverses a struct and assigns fields from m by JSON tag.
// - Supports bools and numbers (float64/int/float32 → float64).
// - Reports unknown keys and type mismatches.
func applyJSONMapByTags(structV reflect.Value, m map[string]any, report *ApplyReport, path string) {
	if structV.Kind() != reflect.Struct {
		return
	}
	structT := structV.Type()

	tagIndex := make(map[string]int, structT.NumField())
	for i := 0; i < structT.NumField(); i++ {
		sf := structT.Field(i)
		if sf.PkgPath != "" { // unexported
			continue
		}
		tag := tagBase(sf.Tag.Get("json"))
		if tag == "" || tag == "-" {
			continue
		}
		tagIndex[tag] = i
	}

	for rawKey, rawVal := range m {
		key := strings.ToLower(rawKey)

		idx, ok := tagIndex[key]
		curPath := key
		if path != "" {
			curPath = path + "." + key
		}
		if !ok {
			report.appendUnknown(curPath)
			continue
		}
		fieldV := structV.Field(idx)
		fieldT := fieldV.Type()

		switch fieldV.Kind() {
		case reflect.Struct:
			childMap, ok := rawVal.(map[string]any)
			if !ok {
				report.appendTypeErr(curPath, "object", typeName(rawVal))
				continue
			}
			applyJSONMapByTags(fieldV, childMap, report, curPath)

		case reflect.Float64:
			if num, ok := asFloat64(rawVal); ok && fieldV.CanSet() {
				fieldV.SetFloat(num)
			} else {
				report.appendTypeErr(curPath, "number", typeName(rawVal))
			}

		case reflect.Bool:
			if b, ok := rawVal.(bool); ok && fieldV.CanSet() {
				fieldV.SetBool(b)
			} else {
				report.appendTypeErr(curPath, "boolean", typeName(rawVal))
			}

		default:
			// Not supported today: report type mismatch.
			report.appendTypeErr(curPath, fieldT.String(), typeName(rawVal))
		}
	}
}

func tagBase(tag string) string {
	if tag == "" {
		return ""
	}
	return strings.Split(tag, ",")[0]
}

func asFloat64(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case int32:
		return float64(t), true
	default:
		return 0, false
	}
}

func typeName(v any) string {
	if v == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%T", v)
}
