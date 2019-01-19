
import * as routes from "api/routes";

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

const post = data => ({
	method: "POST",
	headers: {
		"Content-Type": "application/json",
	},
	body: JSON.stringify(data)
});
