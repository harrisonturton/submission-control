import React, { Component } from "react";
import search from "./assets/search.png";
import Header from "./common/Header.js";

export default class App extends Component {
	render() {
		return (
			<div className="app">
				<Header
					currentCourse={"Concurrent & Distributed Programming"}
					courses={[
						"Computer Networks",
						"Structured Programming",
						"Programming as Problem Solving"
					]}
				/>
				<h1>Test</h1>
			</div>
		);
	}
}
