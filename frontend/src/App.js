import React, { Component } from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import LoginScreen from "./login/LoginScreen.js";
import Header from "./common/Header.js";

const Course = ({ match }) => {
	const course_code = match.params.course_code;
	if (course_code === undefined || course_code === null) {
		console.log("error, course_code is undefined!");
		return <h1>Course code is not found</h1>;
	}
	return <h1>Course: {course_code}</h1>;
};

const Home = () => (
	<Header
		currentCourse={"Concurrent & Distributed Programming"}
		courses={[
			"Computer Networks",
			"Structured Programming",
			"Programming as Problem Solving",
		]}
	/>
);

const App = () => (
	<Router>
		<Switch>
			<Route exact path="/"                    component={Home}/>
			<Route       path="/login"               component={LoginScreen}/>
			<Route       path="/course/:course_code" component={Course}/>
		</Switch>
	</Router>
);

export default App;
