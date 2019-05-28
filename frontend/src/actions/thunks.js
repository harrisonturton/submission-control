
import * as auth from "api/auth";
import * as api from "api/api";
import * as auth_action from "actions/auth";
import * as api_action from "actions/api";
import { forgetState } from "util/state";

// attemptSignIn will attempt to authenticate with the backend server,
// and receives a JWT token if successful. It will also dispatch
// thunks to fetch the initial state and populate the redux store.
export const attemptSignIn = (uid, password) => dispatch => {
	dispatch(auth_action.loginRequest())
	return auth.signIn(uid, password).then(token => {
		// If the token is null, then we failed to sign in.
		// Wipe any tokens in localStorage.
		if (token === null) {
			dispatch(auth_action.loginFailure());
			return;
		}
		// Store the token for later access, and refresh the
		// token in the future. Get the initial state for the store.
		dispatch(auth_action.loginSuccess());
		dispatch(fetchInitialState(token, uid));
		setTimeout(() => dispatch(attemptRefreshToken()), auth.refresh_time);
	});
}

// attemptRefreshToken will exchange our current token with a new one, so we
// can stay authenticated (since each token is time-limited).
// It will try and retreive the token from localStorage. If it does not
// exist, then we assume the user has logged out, and we do not refresh.
export const attemptRefreshToken = () => (dispatch, getState) => {
	let { is_authenticated, token, timestamp }	= getState().auth;
	// If already logged out, then do not refresh the token.
	if (!is_authenticated) {
		return;
	}
	// If the token doesn't exist, then assume logged out
	if (token === undefined || token === null) {
		return;	
	}
	// Check if the token has timed out
	let time_since_refresh = new Date() - timestamp;
	if (time_since_refresh >= auth.token_timeout) {
		dispatch(logout());
		return;	
	}
	// Make the refresh request. This is very similar to the behaviour
	// in signIn().
	return auth.refreshToken(token).then(token => {
		if (token === null) {
			dispatch(logout());
			return;
		}
		dispatch(auth_action.refreshToken(token));
		setTimeout(() => dispatch(attemptRefreshToken()), auth.refresh_time);
	});
};

// logout will remove persisted state (in localStorage) and wipe
// our in-memory state.
export const logout = () => dispatch => {
	forgetState();
	dispatch(auth_action.logout());
}

// fetchInitialState will make multiple requests to the backend
// to build up a state object, and then send this to our store.
// It needs the users email to begin collecting user-specific data.
export const fetchInitialState = (token, uid) => dispatch => {
	dispatch(api_action.dataRequest());
	// First, get the user data
	return api.fetchStudentState(uid, token).then(data => {
		if (data === null) {
			dispatch(api_action.dataFailure());
			return;
		}
		console.log("recieved data on login:\n", JSON.stringify(data));
		dispatch(api_action.dataSuccess(data));
	});
};
