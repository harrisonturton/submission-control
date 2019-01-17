
export const DATA_REQUEST = "DATA_REQUEST";
export const DATA_SUCCESS = "DATA_SUCCESS";
export const DATA_FAILURE = "DATA_FAILURE";

export const dataRequest = data => ({
	type: DATA_REQUEST,
	data: {
		is_fetching: true,
		failed: false
	}
});

export const dataSuccess = data => ({
	type: DATA_SUCCESS,
	data: {
		is_fetching: false,
		failed: false,
		user: data.user,
		courses: data.courses,
		assessment: data.assessment
	}
});

export const dataFailure = () => ({
	type: DATA_FAILURE,
	data: {
		is_fetching: false,
		failed: true,
	}
});
