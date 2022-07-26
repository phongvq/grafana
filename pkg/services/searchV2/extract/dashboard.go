package extract

import (
	"io"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/services/searchV2/dslookup"
)

func logf(format string, a ...interface{}) {
	//fmt.Printf(format, a...)
}

type templateVariable struct {
	current struct {
		value interface{}
	}
	name         string
	query        interface{}
	variableType string
}

type datasourceVariableLookup struct {
	variableNameToRefs map[string][]dslookup.DataSourceRef
	dsLookup           dslookup.DatasourceLookup
}

func (d *datasourceVariableLookup) getDsRefsByTemplateVariableValue(value string, datasourceType string) []dslookup.DataSourceRef {
	switch value {
	case "default":
		// can be the default DS, or a DS with UID="default"
		candidateDs := d.dsLookup.ByRef(&dslookup.DataSourceRef{UID: value})
		if candidateDs == nil {
			// get the actual default DS
			candidateDs = d.dsLookup.ByRef(nil)
		}

		if candidateDs != nil {
			return []dslookup.DataSourceRef{*candidateDs}
		}
		return []dslookup.DataSourceRef{}
	case "$__all":
		// TODO: filter datasources by template variable's regex
		return d.dsLookup.ByType(datasourceType)
	case "":
		return []dslookup.DataSourceRef{}
	case "No data sources found":
		return []dslookup.DataSourceRef{}
	default:
		// some variables use `ds.name` rather `ds.uid`
		if ref := d.dsLookup.ByRef(&dslookup.DataSourceRef{
			UID: value,
		}); ref != nil {
			return []dslookup.DataSourceRef{*ref}
		}

		// discard variable
		return []dslookup.DataSourceRef{}
	}
}

func (d *datasourceVariableLookup) add(templateVariable templateVariable) {
	var refs []dslookup.DataSourceRef

	datasourceType, isDataSourceTypeValid := templateVariable.query.(string)
	if !isDataSourceTypeValid {
		d.variableNameToRefs[templateVariable.name] = refs
		return
	}

	if values, multiValueVariable := templateVariable.current.value.([]interface{}); multiValueVariable {
		for _, value := range values {
			if valueAsString, ok := value.(string); ok {
				refs = append(refs, d.getDsRefsByTemplateVariableValue(valueAsString, datasourceType)...)
			}
		}
	}

	if value, stringValue := templateVariable.current.value.(string); stringValue {
		refs = append(refs, d.getDsRefsByTemplateVariableValue(value, datasourceType)...)
	}

	d.variableNameToRefs[templateVariable.name] = unique(refs)
}

func unique(refs []dslookup.DataSourceRef) []dslookup.DataSourceRef {
	var uniqueRefs []dslookup.DataSourceRef
	uidPresence := make(map[string]bool)
	for _, ref := range refs {
		if !uidPresence[ref.UID] {
			uidPresence[ref.UID] = true
			uniqueRefs = append(uniqueRefs, ref)
		}
	}
	return uniqueRefs
}

func (d *datasourceVariableLookup) getDatasourceRefs(name string) []dslookup.DataSourceRef {
	refs, ok := d.variableNameToRefs[name]
	if ok {
		return refs
	}

	return []dslookup.DataSourceRef{}
}

func newDatasourceVariableLookup(dsLookup dslookup.DatasourceLookup) *datasourceVariableLookup {
	return &datasourceVariableLookup{
		variableNameToRefs: make(map[string][]dslookup.DataSourceRef),
		dsLookup:           dsLookup,
	}
}

// nolint:gocyclo
// ReadDashboard will take a byte stream and return dashboard info
func ReadDashboard(stream io.Reader, lookup dslookup.DatasourceLookup) (*DashboardInfo, error) {
	dash := &DashboardInfo{}

	iter := jsoniter.Parse(jsoniter.ConfigDefault, stream, 1024)

	datasourceVariablesLookup := newDatasourceVariableLookup(lookup)

	for l1Field := iter.ReadObject(); l1Field != ""; l1Field = iter.ReadObject() {
		// Skip null values so we don't need special int handling
		if iter.WhatIsNext() == jsoniter.NilValue {
			iter.Skip()
			continue
		}

		switch l1Field {
		case "id":
			dash.ID = iter.ReadInt64()

		case "uid":
			iter.ReadString()

		case "title":
			dash.Title = iter.ReadString()

		case "description":
			dash.Description = iter.ReadString()

		case "schemaVersion":
			switch iter.WhatIsNext() {
			case jsoniter.NumberValue:
				dash.SchemaVersion = iter.ReadInt64()
			case jsoniter.StringValue:
				val := iter.ReadString()
				if v, err := strconv.ParseInt(val, 10, 64); err == nil {
					dash.SchemaVersion = v
				}
			default:
				iter.Skip()
			}
		case "timezone":
			dash.TimeZone = iter.ReadString()

		case "editable":
			dash.ReadOnly = !iter.ReadBool()

		case "refresh":
			nxt := iter.WhatIsNext()
			if nxt == jsoniter.StringValue {
				dash.Refresh = iter.ReadString()
			} else {
				iter.Skip()
			}

		case "tags":
			for iter.ReadArray() {
				dash.Tags = append(dash.Tags, iter.ReadString())
			}

		case "links":
			for iter.ReadArray() {
				iter.Skip()
				dash.LinkCount++
			}

		case "time":
			obj, ok := iter.Read().(map[string]interface{})
			if ok {
				if timeFrom, ok := obj["from"].(string); ok {
					dash.TimeFrom = timeFrom
				}
				if timeTo, ok := obj["to"].(string); ok {
					dash.TimeTo = timeTo
				}
			}
		case "panels":
			for iter.ReadArray() {
				dash.Panels = append(dash.Panels, readPanelInfo(iter, lookup))
			}

		case "rows":
			for iter.ReadArray() {
				v := iter.Read()
				logf("[DASHBOARD.ROW???] id=%s // %v\n", dash.ID, v)
			}

		case "annotations":
			for sub := iter.ReadObject(); sub != ""; sub = iter.ReadObject() {
				if sub == "list" {
					for iter.ReadArray() {
						v := iter.Read()
						logf("[dash.anno] %v\n", v)
					}
				} else {
					iter.Skip()
				}
			}

		case "templating":
			for sub := iter.ReadObject(); sub != ""; sub = iter.ReadObject() {
				if sub == "list" {
					for iter.ReadArray() {
						templateVariable := templateVariable{}

						for k := iter.ReadObject(); k != ""; k = iter.ReadObject() {
							switch k {
							case "name":
								name := iter.ReadString()
								dash.TemplateVars = append(dash.TemplateVars, name)
								templateVariable.name = name
							case "type":
								templateVariable.variableType = iter.ReadString()
							case "query":
								templateVariable.query = iter.Read()
							case "current":
								for c := iter.ReadObject(); c != ""; c = iter.ReadObject() {
									if c == "value" {
										templateVariable.current.value = iter.Read()
									} else {
										iter.Skip()
									}
								}
							default:
								iter.Skip()
							}
						}

						if templateVariable.variableType == "datasource" {
							log.New("dash").Info("Adding new template variable", "var", templateVariable)
							datasourceVariablesLookup.add(templateVariable)
						}
					}
				} else {
					iter.Skip()
				}
			}

		// Ignore these properties
		case "timepicker":
			fallthrough
		case "version":
			fallthrough
		case "iteration":
			iter.Skip()

		default:
			v := iter.Read()
			logf("[DASHBOARD] support key: %s / %v\n", l1Field, v)
		}
	}

	replaceDatasourceVariables(dash, datasourceVariablesLookup)
	fillDefaultDatasources(dash, lookup)
	filterOutSpecialDatasources(dash)

	targets := newTargetInfo(lookup)
	for _, panel := range dash.Panels {
		targets.addPanel(panel)
	}
	dash.Datasource = targets.GetDatasourceInfo()

	return dash, iter.Error
}

func panelRequiresDatasource(panel PanelInfo) bool {
	return panel.Type != "row"
}

func fillDefaultDatasources(dash *DashboardInfo, lookup dslookup.DatasourceLookup) {
	for i, panel := range dash.Panels {
		if len(panel.Datasource) != 0 || !panelRequiresDatasource(panel) {
			continue
		}

		defaultDs := lookup.ByRef(nil)
		if defaultDs != nil {
			dash.Panels[i].Datasource = []dslookup.DataSourceRef{*defaultDs}
		}
	}
}

func filterOutSpecialDatasources(dash *DashboardInfo) {
	for i, panel := range dash.Panels {
		var dsRefs []dslookup.DataSourceRef

		// partition into actual datasource references and variables
		for _, ds := range panel.Datasource {
			switch ds.UID {
			case "-- Mixed --":
				// The actual datasources used as targets will remain
				continue
			case "-- Dashboard --":
				// The `Dashboard` datasource refers to the results of the query used in another panel
				continue
			default:
				dsRefs = append(dsRefs, ds)
			}
		}

		dash.Panels[i].Datasource = dsRefs
	}
}

func replaceDatasourceVariables(dash *DashboardInfo, datasourceVariablesLookup *datasourceVariableLookup) {
	for i, panel := range dash.Panels {
		var dsVariableRefs []dslookup.DataSourceRef
		var dsRefs []dslookup.DataSourceRef

		// partition into actual datasource references and variables
		for i := range panel.Datasource {
			uid := panel.Datasource[i].UID
			if isVariableRef(uid) {
				dsVariableRefs = append(dsVariableRefs, panel.Datasource[i])
			} else {
				dsRefs = append(dsRefs, panel.Datasource[i])
			}
		}

		variables := findDatasourceRefsForVariables(dsVariableRefs, datasourceVariablesLookup)
		dash.Panels[i].Datasource = append(dsRefs, variables...)
	}
}

func isSpecialDatasource(uid string) bool {
	return uid == "-- Mixed --" || uid == "-- Dashboard --"
}

func isVariableRef(uid string) bool {
	return strings.HasPrefix(uid, "$")
}

func getDataSourceVariableName(dsVariableRef dslookup.DataSourceRef) string {
	if strings.HasPrefix(dsVariableRef.UID, "${") {
		return strings.TrimPrefix(strings.TrimSuffix(dsVariableRef.UID, "}"), "${")
	}

	return strings.TrimPrefix(dsVariableRef.UID, "$")
}

func findDatasourceRefsForVariables(dsVariableRefs []dslookup.DataSourceRef, datasourceVariablesLookup *datasourceVariableLookup) []dslookup.DataSourceRef {
	var referencedDs []dslookup.DataSourceRef
	for _, dsVariableRef := range dsVariableRefs {
		variableName := getDataSourceVariableName(dsVariableRef)
		refs := datasourceVariablesLookup.getDatasourceRefs(variableName)
		referencedDs = append(referencedDs, refs...)
	}
	return referencedDs
}

// will always return strings for now
func readPanelInfo(iter *jsoniter.Iterator, lookup dslookup.DatasourceLookup) PanelInfo {
	panel := PanelInfo{}

	targets := newTargetInfo(lookup)

	for l1Field := iter.ReadObject(); l1Field != ""; l1Field = iter.ReadObject() {
		if iter.WhatIsNext() == jsoniter.NilValue {
			if l1Field == "datasource" {
				targets.addDatasource(iter)
				continue
			}

			// Skip null values so we don't need special int handling
			iter.Skip()
			continue
		}

		switch l1Field {
		case "id":
			panel.ID = iter.ReadInt64()

		case "type":
			panel.Type = iter.ReadString()

		case "title":
			panel.Title = iter.ReadString()

		case "description":
			panel.Description = iter.ReadString()

		case "pluginVersion":
			panel.PluginVersion = iter.ReadString() // since 7x (the saved version for the plugin model)

		case "datasource":
			targets.addDatasource(iter)

		case "targets":
			switch iter.WhatIsNext() {
			case jsoniter.ArrayValue:
				for iter.ReadArray() {
					targets.addTarget(iter)
				}
			case jsoniter.ObjectValue:
				for f := iter.ReadObject(); f != ""; f = iter.ReadObject() {
					targets.addTarget(iter)
				}
			default:
				iter.Skip()
			}

		case "transformations":
			for iter.ReadArray() {
				for sub := iter.ReadObject(); sub != ""; sub = iter.ReadObject() {
					if sub == "id" {
						panel.Transformer = append(panel.Transformer, iter.ReadString())
					} else {
						iter.Skip()
					}
				}
			}

		// Rows have nested panels
		case "panels":
			for iter.ReadArray() {
				panel.Collapsed = append(panel.Collapsed, readPanelInfo(iter, lookup))
			}

		case "options":
			fallthrough

		case "gridPos":
			fallthrough

		case "fieldConfig":
			iter.Skip()

		default:
			v := iter.Read()
			logf("[PANEL] support key: %s / %v\n", l1Field, v)
		}
	}

	panel.Datasource = targets.GetDatasourceInfo()

	return panel
}
