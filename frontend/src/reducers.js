
import { token_timeout } from "auth";
import {
	LOGIN_REQUEST,
	LOGIN_SUCCESS,
	LOGIN_FAILURE,
	LOGOUT,
	REFRESH_TOKEN,
	CHECK_TOKEN_TIMEOUT
} from "actions";

let INITIAL_STATE = {
	auth: {
		is_authenticated: false,
		is_fetching: false,
		token: null,
		timestamp: null,
	},
	user: {
		email: null,
		first_name: null,
		last_name: null,
		current_course: null,
		courses: []
	},
	assignments_by_course: {},
	labs_by_course: {}
};

let appReducer = (prev_state = INITIAL_STATE, action) => ({
	...prev_state,
	auth: authReducer(prev_state.auth, action)
});

let authReducer = (prev_state = INITIAL_STATE.auth, action) => {
	console.log(action.type);
	switch (action.type) {
		case CHECK_TOKEN_TIMEOUT:
			console.log("Reducing CHECK_TOKEN_TIMEOUT")
			return checkTokenTimeout(prev_state, action)
		case LOGIN_REQUEST:
		case LOGIN_SUCCESS:
		case LOGIN_FAILURE:
		case LOGOUT:
		case REFRESH_TOKEN:
			return { ...prev_state, ...action.auth };
		default:
			return prev_state;
	}
};

// checkTokenTimeout is called on CHECK_TOKEN_TIMEOUT actions.
// It compares the tokens timestamp with the current time to
// determine if it has timed out or not.
function checkTokenTimeout(prev_state, action) {
	console.log("Checking again....")
	let time_since_token = action.timestamp - prev_state.timestamp;
	if (time_since_token >= token_timeout) {
		return {
			is_authenticated: false,
			is_fetching: false,
			token: null,
			timestamp: action.timestamp
		};
	}
	return prev_state;
}

export {
	INITIAL_STATE,
	appReducer
}
