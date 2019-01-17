
import { auth_url, refresh_url, enrol_url, assessment_url, user_url } from "api/routes";

export const fetchUserData = async (email, token) => {
	// Make sure request completes successfully
	let url = `${user_url}?email=${email}`;
	let resp = await fetch(url, {
		method: "GET",
		headers: {
			"Authorization": token	
		}
	});
	if (!resp.ok) {
		console.log(resp.statusText);
		return null;	
	}
	// Check for data
	let body = await resp.json();
	if (body.user === null || body.user === undefined) {
		console.log("body.user is null");
		return null;
	}
	return body.user;
};

// fetchCourses returns a list of courses a user is enrolled in
export const fetchCourses = async (uid, token) => {
	// Make sure request completes successfully
	let url = `${enrol_url}?uid=${uid}`;
	let resp = await fetch(url, {
		method: "GET",
		headers: {
			"Authorization": token	
		}
	});
	if (!resp.ok) {
		console.log(resp.statusText);
		return null;	
	}
	// Check for data
	let body = await resp.json();
	if (body.courses === null || body.courses === undefined) {
		return null;	
	}
	return body.courses;
}

export const fetchAssessmentsForCourse = async (course_id, token) => {
	// Make sure request completes successfully
	let url = `${assessment_url}?course_id=${course_id}`;
	let resp = await fetch(url, {
		method: "GET",
		headers: {
			"Authorization": token	
		}
	});
	if (!resp.ok) {
		console.log(resp.statusText);
		return null;	
	}
	// Check for data
	let body = await resp.json();
	if (body.assessment === null || body.assessment == null) {
		return null;	
	}
	return {
		assignments: body.assessment.filter(item => item.type == "assignment"),
		labs: body.assessment.filter(item => item.type == "lab"),
	}
}
