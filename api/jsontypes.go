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

var drugFields = []string {"effects", "onset", "dose", "aftereffects", "after-effects", "aliases", "categories"}

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
	unit := (*m2)["_unit"]
	for x, y := range *m2 {
		if x == "_unit" {
			continue
		}
		str = str + " `" + x + "` " + y + " " + unit
	}
	(*m1)[name] = str
}

func trim(s string) string {
	return s[:len(s)-2]
}

func (d Drug) Fields() *map[string]string {
	p := d.StringProperties()
	if d.DoseNote != "" {
		p["dose note"] = d.DoseNote
	}
	sp := &p
	p["factsheet"] = "https://drugs.tripsit.me/" + d.PrettyName
	addAll(sp, &d.Onset, "onset")
	addAll(sp, &d.Duration, "duration")
	addAll(sp, &d.Aftereffects, "after-effects")
	return sp
}

func (d Drug) TableFields() *map[string]map[string]map[string]string {
	cast := (map[string]map[string]string(d.Dose))
	ret := map[string]map[string]map[string]string{}
	ret["Dose"] = cast
	return &ret
}

func (d Drug) MultipleFields() *map[string][]string {
	a := map[string][]string {"effects": d.Effects, "categories": d.Categories}
	return &a
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
