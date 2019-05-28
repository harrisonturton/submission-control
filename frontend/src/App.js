import React from "react";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import {
	HomeScreen,
	LoginScreen,
	CourseScreen,
	AssessmentScreen,
	SubmissionScreen,
	ManageStudentsScreen,
	ManageTutorialsScreen
} from "containers";

const App = () => (
	<Router>
		<Switch>
			<Route
				path="/login"
				component={LoginScreen}
			/>
			<Route
				path="/course/:course_id/admin/students"
				component={ManageStudentsScreen}
			/>
			<Route
				path="/course/:course_id/admin/tutorials"
				component={ManageTutorialsScreen}
			/>
			<Route
				path="/course/:course_id/assessment/:assessment_id/:submission_id"
				component={SubmissionScreen}
			/>
			<Route
				path="/course/:course_id/assessment/:assessment_id"
				component={AssessmentScreen}
			/>
			<Route
				path="/course/:course_id"
				component={CourseScreen}
			/>
			<Route
				exact path="/"
				component={HomeScreen}
			/>
		</Switch>
	</Router>
);

export default App;
