
import {
	SET_COURSES,
	SET_USER,
	SET_ASSESSMENT,
	setCourses,
	setUser,
	setAssessment
} from "actions/api";
import {
	LOGIN_REQUEST,
	LOGIN_SUCCESS,
	LOGIN_FAILURE,
	LOGOUT,
	REFRESH_TOKEN,
	CHECK_TOKEN_TIMEOUT,
	loginRequest,
	loginSuccess,
	loginFailure,
	logout,
	refreshToken,
	checkTokenTimeout
} from "actions/auth";
import {
	makeLoginRequest,
	makeRefreshRequest,
	makeEnrolRequest
} from "actions/thunks";

export {
	// Auth
	LOGIN_REQUEST,
	LOGIN_SUCCESS,
	LOGIN_FAILURE,
	LOGOUT,
	loginRequest,
	loginSuccess,
	loginFailure,
	logout,

	REFRESH_TOKEN,
	CHECK_TOKEN_TIMEOUT,
	refreshToken,
	checkTokenTimeout,

	// API
	SET_COURSES,
	SET_USER,
	SET_ASSESSMENT,
	setCourses,
	setUser,
	setAssessment,

	// Thunks
	makeLoginRequest,
	makeRefreshRequest,
	makeEnrolRequest
}
