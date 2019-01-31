
import * as routes from "api/routes";

// fetchUserData will retreive the information about a user
// uniquely identified by their email.
// { uid, first_name, last_name, email }
export const fetchUserData = async (email, token) => {
	let url = `${routes.user}?email=${email}`;
	let body = await get(url, token, [
		"user"
	]);
	return body === null ? null : body.user;
};

// fetchCourses returns a list of courses a user is enrolled in.
// [ { id, name, course_code, period, year } ]
export const fetchCourses = async (uid, token) => {
	let url = `${routes.enrolled}?uid=${uid}`;
	let body = await get(url, token, [
		"courses"	
	]);
	return body === null ? null : body.courses;
};

export const fetchAssessment = async (uid, token) => {
	let url = `${routes.assessment}?uid=${uid}`;
	let body = await get(url, token, [
		"assessment"	
	]);
	return body === null ? null : {
		assignments: body.assessment.filter(item => item.type === "assignment"),
		labs: body.assessment.filter(item => item.type === "lab"),
	};
};

// fetchSubmissions will return an array of all submissions the user
// has ever made.
export const fetchSubmissions = async (uid, token) => {
	// Make sure request completes successfully
	let url = `${routes.submissions}?uid=${uid}`;
	let body = await get(url, token, [
		"submissions"
	]);
	return body.submissions === null ? null : body.submissions;
}

export const fetchStudentState = async (email, token) => {
	let url  = `${routes.studentState}?email=${email}`;
	let body = await get(url, token, [
		"enrolled",
		"user",
		"submissions",
		"assessment"
	]);
	return {
		user: body.user,
		submissions: body.submissions,
		courses: body.enrolled,
		assessment: {
			assignments: body.assessment.filter(item => item.type === "assignment"),
			labs: body.assessment.filter(item => item.type === "lab"),
		}
	}
};

const get = async (url, token, responseKeys) => {
	let resp = await fetch(url, {
		method: "GET",
		headers: {
			"Authorization": token	
		}
	});
	if (!resp.ok) {
		return null;
	}
	let body = await resp.json();
	return keysAreValid(body, responseKeys) ? body : null;
}

// Will check that every key has a value that is not null or undefined.
const keysAreValid = (obj, keys) =>
	Object.keys(obj).every(key => obj[key] !== null && obj[key] !== undefined);
