import thunk from "redux-thunk";
import { createStore, applyMiddleware } from "redux";
import { appReducer, INITIAL_STATE } from "reducers"; 
import { retreiveToken } from "auth";

let appStore = createStore(
	appReducer,
	hydrateState(INITIAL_STATE),
	applyMiddleware(
		thunk,
		logger
	)
);

function hydrateState(initial_state) {
	let token = retreiveToken();
	if (token === undefined || token === null) {
		return initial_state;	
	}
	console.log(JSON.stringify(initial_state));
	return {
		...initial_state,
		is_authenticated: true,
		token: token
	};
}

function logger({ getState }) {
	return next => action => {
		console.log(`%cDispatching ${action.type}`, "font-weight:bold")
		console.log(JSON.stringify(action));
		const returnVal = next(action);
		console.log("%cNew State", "font-weight:bold")
		console.log(JSON.stringify(getState()));
		return returnVal;
	};
}

export default appStore;
