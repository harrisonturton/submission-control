import thunk from "redux-thunk";
import { createStore, applyMiddleware } from "redux";
import { rootReducer } from "reducer";
import { loadState, saveState } from "util/state";

const configureStore = () => {
	const persisted_state = loadState();
	const store = createStore(
		rootReducer,
		persisted_state,
		applyMiddleware(
			thunk,
			logger
		)
	);
	store.subscribe(() => saveState(store.getState()));
	return store;
};

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
	
export default configureStore;
