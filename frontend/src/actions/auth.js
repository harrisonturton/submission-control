
// Authentication actions
export const LOGIN_REQUEST = "LOGIN_REQUEST";
export const LOGIN_SUCCESS = "LOGIN_SUCCESS";
export const LOGIN_FAILURE = "LOGIN_FAILURE";
export const LOGOUT        = "LOGOUT";

export const REFRESH_TOKEN = "REFRESH_TOKEN";

// Check the tokens current timestamp, and detect when it times out.
export const CHECK_TOKEN_TIMEOUT = "CHECK_TOKEN_TIMEOUT";

export const loginRequest = () => ({
	type: LOGIN_REQUEST,
	auth: {
		is_authenticated: false,
		is_fetching: true
	}
});

export const loginSuccess = token => ({
	type: LOGIN_SUCCESS,
	auth: {
		is_authenticated: true,
		is_fetching: false,
		token: token,
		timestamp: new Date()
	}
});

export const loginFailure = () => ({
	type: LOGIN_FAILURE,
	auth: {
		is_authenticated: false,
		is_fetching: false,
	}
})

export const logout = () => ({
	type: LOGOUT,
	auth: {
		is_authenticated: false,
		is_fetching: false,
		timestamp: new Date(),
		token: null
	}
});

export const refreshToken = token => ({
	type: REFRESH_TOKEN,
	auth: {
		is_authenticated: true,
		is_fetching: false,
		token: token,
		timestamp: new Date()
	}
});

export const checkTokenTimeout = () => ({
	type: CHECK_TOKEN_TIMEOUT,
	timestamp: new Date()
});
