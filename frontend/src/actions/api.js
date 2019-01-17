
export const SET_COURSES    = "SET_COURSES";
export const SET_USER       = "SET_USER";
export const SET_ASSESSMENT = "SET_ASSESSMENT";

export const setCourses = courses => ({
	type: SET_COURSES,
	courses: courses
});

export const setUser = user_data => ({
	type: SET_USER,
	user: user_data
});

export const setAssessment = assessment => ({
	type: SET_ASSESSMENT,
	assessment: assessment
});
