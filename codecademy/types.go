package codecademy

// Course is a single catalog entry from Codecademy.
type Course struct {
	Rank             int    `json:"rank"              csv:"rank"              tsv:"rank"`
	Slug             string `json:"slug"              csv:"slug"              tsv:"slug"`
	Title            string `json:"title"             csv:"title"             tsv:"title"`
	Type             string `json:"type"              csv:"type"              tsv:"type"`
	Difficulty       string `json:"difficulty"        csv:"difficulty"        tsv:"difficulty"`
	LessonCount      int    `json:"lesson_count"      csv:"lesson_count"      tsv:"lesson_count"`
	TimeToComplete   int    `json:"time_to_complete"  csv:"time_to_complete"  tsv:"time_to_complete"`
	Pro              bool   `json:"pro"               csv:"pro"               tsv:"pro"`
	Certificate      bool   `json:"certificate"       csv:"certificate"       tsv:"certificate"`
	ShortDescription string `json:"short_description" csv:"short_description" tsv:"short_description"`
	URL              string `json:"url"               csv:"url"               tsv:"url"`
}
