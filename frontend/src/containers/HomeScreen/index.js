import React, { Component } from "react";
import { Header } from "components";

export default class HomeScreen extends Component {
	render = () => (
		<Header
			currentCourse="Concurrent & Distributed Programming"
			courses={[
				"Computer Networks",
				"Structured Programming",
				"Programming as Problem Solving",
			]}
		/>
	);
}
