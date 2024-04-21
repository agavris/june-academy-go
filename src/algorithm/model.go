package algorithm

type Request struct {
	Email     string `csv:"Email Address"`
	FirstName string `csv:"Students First Name"`
	LastName  string `csv:"Students Last Name"`
	Grade     string `csv:"Grade in school this year"`
	AMFD1     string `csv:"AM Course - 1st Choice. (Drop down option)"`
	AMFD2     string `csv:"AM Course - 2nd Choice. (Drop down option)"`
	AMFD3     string `csv:"AM Course - 3rd Choice. (Drop down option)"`
	AMFD4     string `csv:"AM Course - 4th Choice. (Drop down option)"`
	AMFD5     string `csv:"AM Course - 5th Choice. (Drop down option)"`
	PM1       string `csv:"PM Course - 1st Choice. (Drop down option)"`
	PM2       string `csv:"PM Course - 2nd Choice. (Drop down option)"`
	PM3       string `csv:"PM Course - 3rd Choice. (Drop down option)"`
	PM4       string `csv:"PM Course - 4th Choice. (Drop down option)"`
	PM5       string `csv:"PM Course - 5th Choice. (Drop down option)"`
}

func (r *Request) GetAMCourses() []string {
	return []string{r.AMFD1, r.AMFD2, r.AMFD3, r.AMFD4, r.AMFD5}
}

func (r *Request) GetPMCourses() []string {
	return []string{r.PM1, r.PM2, r.PM3, r.PM4, r.PM5}
}
