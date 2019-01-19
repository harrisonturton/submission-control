
import * as routes from "api/routes";

// fetchUserData will retreive the information about a user
// uniquely identified by their email.
// { uid, first_name, last_name, email }
export const fetchUserData = async (email, token) => {
	// Make sure request completes successfully
	let url = `${routes.user}?email=${email}`;
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

// fetchCourses returns a list of courses a user is enrolled in.
// [ { id, name, course_code, period, year } ]
export const fetchCourses = async (uid, token) => {
	// Make sure request completes successfully
	let url = `${routes.enrolled}?uid=${uid}`;
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

// fetchAssessment will return a list of assessments for a user.
export const fetchAssessment = async (uid, token) => {
	// Make sure request completes successfully
	let url = `${routes.assessment}?uid=${uid}`;
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
		assignments: body.assessment.filter(item => item.type === "assignment"),
		labs: body.assessment.filter(item => item.type === "lab"),
	}
}

// fetchSubmissions will return an array of all submissions the user
// has ever made.
export const fetchSubmissions = async (uid, token) => {
	// Make sure request completes successfully
	let url = `${routes.submissions}?uid=${uid}`;
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
	if (body.submissions === null || body.submissions == null) {
		return null;	
	}
	return body.submissions;
}
