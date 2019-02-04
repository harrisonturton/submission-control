
import * as routes from "api/routes";

export const fetchStudentState = async (uid, token) => {
	let url  = `${routes.studentState}?uid=${uid}`;
	let body = await get(url, token, [
		"enrolled",
		"user",
		"submissions",
		"assessment"
	]);
	return {
		user: body.user,
		submissions: body.submissions,
		courses: body.enrolled.map(enrolment => ({
			id: enrolment.course.id,
			name: enrolment.course.name,
			course_code: enrolment.course.course_code,
			period: enrolment.course.period,
			year: enrolment.course.year,
			role: enrolment.role
		})),
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
