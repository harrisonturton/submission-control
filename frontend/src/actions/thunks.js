
import * as auth from "api/auth";
import * as api from "api/api";
import * as auth_action from "actions/auth";
import * as api_action from "actions/api";

// attemptSignIn will attempt to authenticate with the backend server,
// and receives a JWT token if successful. It will also dispatch
// thunks to fetch the initial state and populate the redux store.
export const attemptSignIn = (email, password) => dispatch => {
	dispatch(auth_action.loginRequest())
	return auth.signIn(email, password).then(token => {
		// If the token is null, then we failed to sign in.
		// Wipe any tokens in localStorage.
		if (token === null) {
			auth.forgetToken();
			dispatch(auth_action.loginFailure());
			return;
		}
		// Store the token for later access, and refresh the
		// token in the future. Get the initial state for the store.
		auth.storeToken(token);
		dispatch(auth_action.loginSuccess());
		dispatch(fetchInitialState(token, email));
		setTimeout(() => dispatch(attemptRefreshToken()), auth.refresh_time);
	});
}

// attemptRefreshToken will exchange our current token with a new one, so we
// can stay authenticated (since each token is time-limited).
// It will try and retreive the token from localStorage. If it does not
// exist, then we assume the user has logged out, and we do not refresh.
export const attemptRefreshToken = () => (dispatch, getState) => {
	let { is_authenticated, timestamp }	= getState().auth;
	// If already logged out, then do not refresh the token.
	if (!is_authenticated) {
		return;
	}
	// If the token doesn't exist, then assume logged out
	let token = auth.retreiveToken();
	if (token === undefined || token === null) {
		return;	
	}
	// Check if the token has timed out
	let time_since_refresh = new Date() - timestamp;
	if (time_since_refresh >= auth.token_timeout) {
		dispatch(auth_action.logout());
		return;	
	}
	// Make the refresh request. This is very similar to the behaviour
	// in signIn().
	return auth.refreshToken(token).then(token => {
		if (token === null) {
			auth.forgetToken();
			dispatch(auth_action.logout());
			return;
		}
		auth.storeToken(token);
		dispatch(auth_action.refreshToken(token));
		setTimeout(() => dispatch(attemptRefreshToken()), auth.refresh_time);
	});
};

// fetchInitialState will make multiple requests to the backend
// to build up a state object, and then send this to our store.
// It needs the users email to begin collecting user-specific data.
export const fetchInitialState = (token, email) => dispatch => {
	dispatch(api_action.dataRequest());
	// First, get the user data
	return fetchState(token, email).then(data => {
		if (data === null) {
			dispatch(api_action.dataFailure());
			return;
		}
		dispatch(api_action.dataSuccess(data));
	});
};

const fetchState = async (token, email) => {
	let user = await api.fetchUserData(email, token);
	if (user === null || user === undefined) {
		return null;	
	}
	let courses = await api.fetchCourses(user.uid, token);
	if (courses === null || courses === undefined) {
		return null;
	}
	let assessment = await api.fetchAssessment(user.uid, token);
	if (assessment === null || assessment === undefined) {
		return null;	
	}
	return { user, courses, assessment };
};
