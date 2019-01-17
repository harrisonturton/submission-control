
import * as routes from "api/routes";

// Where we store the JWT token in localStorage
const token_key = "submission-control-token";

// Time between token refreshes. Do it every 3 minutes.
export const refresh_time = 1000 * 60 * 3;

// How long until a JWT token is invalid. 10 minutes.
export const token_timeout = 1000 * 60 * 10;

export const signIn = async (email, password) => {
	let resp = await fetch(routes.auth, post({
		email: email,
		password: password,
	}));
	if (!resp.ok) {
		console.log(resp.statusText);
		return null;	
	}
	// Check for data
	let body = await resp.json();
	if (body.token === null || body.token === undefined) {
		return null;
	}
	return body.token;
};

export const refreshToken = async token => {
	let resp = await fetch(routes.refresh, post({
		token: token
	}));
	if (!resp.ok) {
		console.log(resp.statusText);
		return null;	
	}
	// Check for data
	let body = await resp.json();
	if (body.token === null || body.token === undefined) {
		return null;
	}
	return body.token;
};

// store_token will save a new token to localStorage
export function storeToken(token) {
	localStorage.setItem(token_key, token);
}

export function forgetToken() {
	localStorage.removeItem(token_key);
}

// retreive_token will retreive the JWT token from localStorage
export function retreiveToken() {
	return localStorage.getItem(token_key);
}

export function hasToken() {
	return retreiveToken() !== null && retreiveToken() !== undefined;
}

const post = data => ({
	method: "POST",
	headers: {
		"Content-Type": "application/json",
	},
	body: JSON.stringify(data)
});