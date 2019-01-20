import React from "react";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import {
	PrivateRoute,
	HomeScreen,
	LoginScreen,
	CourseScreen,
	AssessmentScreen,
	SubmissionScreen
} from "containers";

//const CourseScreen = ({ match }) => (
//	<p>{match.params.course_code}</p>
//);

const App = () => (
	<Router>
		<Switch>
			<Route
				path="/login"
				component={LoginScreen}
			/>
			<Route
				path="/course/:course_id/:assessment_id/:submission_id"
				component={SubmissionScreen}
			/>
			<Route
				path="/course/:course_id/:assessment_id"
				component={AssessmentScreen}
			/>
			<Route
				path="/course/:course_id"
				component={CourseScreen}
			/>
			<PrivateRoute
				exact path="/"
				component={HomeScreen}
			/>
		</Switch>
	</Router>
);

export default App;
