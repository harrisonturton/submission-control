
import {
	signIn,
	refresh,
	refresh_time,
	token_timeout,
	storeToken,
	forgetToken
} from "auth";

// Authentication actions
export const LOGIN_REQUEST = "LOGIN_REQUEST";
export const LOGIN_SUCCESS = "LOGIN_SUCCESS";
export const LOGIN_FAILURE = "LOGIN_FAILURE";
export const LOGOUT        = "LOGOUT";

// Dispatched on every PrivateRoute to ensure we're logged out
// when the token times out.
export const CHECK_TOKEN_TIMEOUT = "CHECK_TOKEN_TIMEOUT";

export const REFRESH_TOKEN = "REFRESH_TOKEN";

// makeLoginRequest will handle dispatching various actions to
// represent the API flow (not authenticated, fetching, authenticated).
export function makeLoginRequest(email, password) {
	return dispatch => {
		dispatch(loginRequest());
		return signIn(email, password)
			.then(token => {
				if (token === null) {
					forgetToken();
					dispatch(loginFailure());
					console.log("Login failed");
					return;
				}
				console.log("Login succeeded!");
				storeToken(token);
				dispatch(loginSuccess(token));
				/*setTimeout(() => {
					console.log("Dispatching from inside makeLoginRequest");
					dispatch(makeRefreshRequest());
				}, refresh_time);*/
			});
	};
}

// makeRefreshRequest will try to replace the current JWT token
// with a new one, allowing us to stay authenticated without
// prompting the user for login details every ~5minutes.
export const makeRefreshRequest = () => (dispatch, getState) => {
	console.log("Refreshing...");
	// Check that we can still refresh
	let { is_authenticated, token, timestamp } = getState().auth;
	if (!is_authenticated) {
		dispatch(logout())
		return;		
	}
	let time_since_refresh = new Date() - timestamp;
	if (time_since_refresh > token_timeout) {
		dispatch(logout())
		return;
	}
	// Refresh
	return refresh(token)
		.then(token => {
			if (token === null) {
				forgetToken();
				dispatch(logout());
				console.log("Refresh failed");
				return;
			}
			storeToken(token);
			dispatch(refreshToken(token));
			setTimeout(() => {
				console.log("Dispatching from inside makeRefreshRequest");
				dispatch(makeRefreshRequest());
			}, refresh_time);
		})
};


const loginRequest = () => ({
	type: LOGIN_REQUEST,
	auth: {
		is_authenticated: false,
		is_fetching: true
	}
});

const loginSuccess = token => ({
	type: LOGIN_SUCCESS,
	auth: {
		is_authenticated: true,
		is_fetching: false,
		token: token,
		timestamp: new Date()
	}
});

const loginFailure = () => ({
	type: LOGIN_FAILURE,
	auth: {
		is_authenticated: false,
		is_fetching: false,
	}
})

const logout = () => ({
	type: LOGOUT,
	auth: {
		is_authenticated: false,
		is_fetching: false,
		timestamp: new Date(),
		token: null
	}
});

const refreshToken = token => ({
	type: REFRESH_TOKEN,
	auth: {
		is_authenticated: true,
		is_fetching: false,
		token: token,
		timestamp: new Date()
	}
});

export const checkTokenTimeout = () => {
	console.log("Dispatching action")
	return {
		type: CHECK_TOKEN_TIMEOUT,
		timestamp: new Date()
	};
};
