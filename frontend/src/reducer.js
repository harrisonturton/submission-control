
import { combineReducers } from "redux";
import * as auth from "actions/auth";
import * as data from "actions/api";

export const initial_state = {
	auth: {
		is_authenticated: false,
		is_fetching: false,
		failed: false,
		token: null,
		timestamp: null,
	},
	data: {
		is_fetching: false,
		failed: false,
		user: {
			email: null,
			first_name: null,
			last_name: null,
		},
		courses: [],
		assessment: []
	}
};

export const rootReducer = (state=initial_state, action) => ({
	auth: authReducer(state.auth, action),
	data: dataReducer(state.data, action)
});

const authReducer = (state=initial_state.auth, action) => {
	switch (action.type) {
		case auth.LOGIN_REQUEST:
		case auth.LOGIN_SUCCESS:
		case auth.LOGIN_FAILURE:
		case auth.LOGOUT:
		case auth.REFRESH_TOKEN:
			return { ...state, ...action.auth };
		default:
			return state;
	};
};

const dataReducer = (state=initial_state.data, action) => {
	switch (action.type) {
		case data.DATA_REQUEST:
		case data.DATA_SUCCESS:
		case data.DATA_FAILURE:
			return { ...state, ...action.data };
		default:
			return state;
	}
};
