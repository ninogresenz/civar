package gitlab

type CiVariable struct {
	Key              string `json:"key"`
	VariableType     string `json:"variable_type"`
	Value            string `json:"value"`
	Protected        bool   `json:"protected"`
	Masked           bool   `json:"masked"`
	EnvironmentScope string `json:"environment_scope"`
}

type CiVariableList []CiVariable

func (c *CiVariableList) Includes(needle CiVariable) bool {
	for _, variable := range *c {
		if variable.Key == needle.Key && variable.EnvironmentScope == needle.EnvironmentScope {
			return true
		}
	}
	return false
}
func (c *CiVariableList) Push(newElement ...CiVariable) {
	*c = append(*c, newElement...)
}

type Project struct {
	PathWithNamespace string `json:"path_with_namespace"`
	WebUrl            string `json:"web_url"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type UpdateBody struct {
	Value  string `url:"value"`
	Filter Filter `url:"filter"`
}

type Filter struct {
	EnvironmentScope string `url:"environment_scope"`
}
