
import {
	// Authentication
	sendSignInRequest,
	sendRefreshRequest,
	storeToken,
	forgetToken,
	refresh_time,
	token_timeout,
	// Data
	fetchCourses,
	fetchUserData,
	fetchAssessmentsForCourse
} from "api";
import {
	loginRequest,
	loginSuccess,
	loginFailure,
	logout,
	refreshToken,
	setCourses,
	setUser,
	setAssessment
} from "actions";

// loginThunk sends a login request to the API
export const makeLoginRequest = (email, password) => dispatch => {
	console.log("LOGGING IN...");
	dispatch(loginRequest());
	return sendSignInRequest(email, password)
		.then(token => {
			if (token === null) {
				forgetToken()
				dispatch(loginFailure());
			} else {
				storeToken(token);
				dispatch(fetchState(email, token));
				setTimeout(() => dispatch(makeRefreshRequest()), refresh_time);
				dispatch(loginSuccess(token));
			}
		});
};

// refreshThunk sends a request to refresh out JWT token
export const makeRefreshRequest = () => (dispatch, getState) => {
	let { is_authenticated, token, timestamp } = getState().auth;
	if (!is_authenticated) {
		dispatch(logout());
		return;	
	}
	// Check if token has timed out
	let time_since_refresh = new Date() - timestamp;
	if (time_since_refresh >= token_timeout) {
		dispatch(logout());
		return;
	}
	// Make refresh request
	return sendRefreshRequest(token)
		.then(token => {
			if (token === null) {
				forgetToken();
				dispatch(logout());
			} else {
				storeToken(token);
				dispatch(refreshToken(token));
				setTimeout(() => dispatch(makeRefreshRequest()), refresh_time);
			}
		});
}

// fetchState will hydrate the entire store state from the database
export const fetchState = (email, token) => async dispatch => {
	console.log("INSIDE FETCH STATE")
	function onSuccess(user_data, courses, assessments) {
		console.log("FETCHING STATE...");
		console.log(JSON.stringify(user_data));
		console.log(JSON.stringify(courses));
		dispatch(setUser(user_data));
		dispatch(setCourses(courses));
		dispatch(setAssessment(assessments))
	}
	function onError(error) {
		return error;	
	}
	try {
		console.log("Trying...");
		const user_data = await fetchUserData(email, token);
		if (user_data === null) {
			return;	
		}
		console.log("User data", JSON.stringify(user_data));
		const courses = await fetchCourses(user_data.uid, token);
		if (courses === null) {
			return;
		}
		var assessments = {};
		for (var i = 0; i < courses.length; i++) {
			var course = courses[i];
			let assessment = await fetchAssessmentsForCourse(course.id, token);
			if (assessment === null) {
				continue;
			}
			assessments[course.id] = assessment;
		}
		return onSuccess(user_data, courses, assessments);
	} catch (err) {
		return onError(err);
	}
};

export const makeEnrolRequest = uid => (dispatch, getState) => {
	let { is_authenticated, token } = getState().auth;
	if (!is_authenticated) {
		dispatch(logout());
		return;
	}
	return fetchCourses(uid, token)
		.then(courses => {
			if (courses === null) {
				console.log("Courses is null");
				return;
			} else {
				dispatch(setCourses(courses))	
			}
		})
};

