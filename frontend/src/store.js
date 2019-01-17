import thunk from "redux-thunk";
import { createStore, applyMiddleware } from "redux";
import { rootReducer, initial_state } from "reducer";
import { retreiveToken } from "api/auth";

let store = createStore(
	rootReducer,
	initial_state,
	// hydrateAuthentication(initial_state),
	applyMiddleware(thunk, logger)
);

// hydrateAuthentication looks for a token in localStorage,
// and sets is_authenticated accordingly. Used to keep authentication
// between page refreshes.
function hydrateAuthentication(initial_state) {
	let token = retreiveToken();
	if (token === undefined || token === null) {
		return initial_state;	
	}
	console.log(JSON.stringify(initial_state));
	return {
		...initial_state,
		auth: {
			...initial_state.auth,	
			is_authenticated: true,
			token: token
		}
	};
}

function logger({ getState })  {
	return next => action => {
		console.log(`%cDispatching ${action.type}`, "font-weight:bold")
		console.log(JSON.stringify(action));
		const returnVal = next(action);
		console.log("%cNew State", "font-weight:bold")
		console.log(JSON.stringify(getState()));
		return returnVal;
		// return next(action);
	};
}
	
export default store;
