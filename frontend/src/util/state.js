
const state_key = "state";

export const loadState = () => {
	try {
		const serialized_state = localStorage.getItem(state_key);
		if (serialized_state === null) {
			return undefined;	
		}
		return JSON.parse(serialized_state);
	} catch (err) {
		return undefined;
	}
};

export const saveState = state => {
	try {
		const serialized_state = JSON.stringify(state);
		localStorage.setItem(state_key, serialized_state);
	} catch (err) {
		// Ignore errors
		console.log("Error saving state: ", err);
	}
};
