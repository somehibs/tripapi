package tripapi

type FormattedItem map[string]string

type FormattedDose map[string]map[string]string

type Drug struct {
	Name string
	Categories []string
	Aliases []string
	Properties map[string]interface{} `json:"properties"`
	PrettyName string `json:"pretty_name"`
	DoseNote string `json:"dose_note"`
	Effects []string `json:"formatted_effects"`
	Onset FormattedItem `json:"formatted_onset"`
	Duration FormattedItem `json:"formatted_duration"`
	Dose FormattedDose `json:"formatted_dose"`
	Aftereffects FormattedItem `json:"formatted_aftereffects"`
}

var drugFields = []string {"effects", "onset", "duration", "dose", "aftereffects", "after-effects", "aliases", "categories"}

type DrugData map[string]Drug

type DrugReply struct {
	Err string
	Data []DrugData
}


func addAll(m1 *map[string]string, m2 *FormattedItem, name string) {
	str := ""
	if (*m2)["value"] != "" {
		(*m1)[name] = (*m2)["value"] + " " + (*m2)["_unit"]
		return
	}
	for x, y := range *m2 {
		str = str + " " + x + " " + y
	}
	(*m1)[name] = str
}

func (d Drug) Fields() *map[string]string {
	p := d.StringProperties()
	if d.DoseNote != "" {
		p["dose note"] = d.DoseNote
	}
	sp := &p
	addAll(sp, &d.Onset, "onset")
	addAll(sp, &d.Duration, "duration")
	addAll(sp, &d.Aftereffects, "after-effects")
	return sp
}

func (d Drug) TableFields() *map[string]map[string]map[string]string {
	cast := (map[string]map[string]string(d.Dose))
	ret := map[string]map[string]map[string]string{}
	ret["dose"] = cast
	return &ret
}

func (d Drug) ComplexFields() *map[string]map[string]string {
	return nil
}

func (d *Drug) StringProperties() map[string]string {
	props := map[string]string {}
	for k, v := range d.Properties {
		if d.FormattedField(k) {
			continue
		}
		if len(v.(string)) != 0 {
			props[k] = v.(string)
		}
	}
	return props
}

func (d *Drug) FormattedField(field string) bool {
	for _, x := range drugFields {
		if x == field {
			return true
		}
	}
	return false
}
