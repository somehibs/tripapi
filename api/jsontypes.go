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

type DrugReply struct {
	Err string
	Data []map[string]Drug
}
