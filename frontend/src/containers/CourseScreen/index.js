import React, { Component } from "react";
import { connect } from "react-redux";
import { Redirect } from "react-router-dom";
import { Header } from "components";
import { hasToken } from "auth";
import "./style.css";

class Course extends Component {
	render() {
		let is_authenticated = hasToken();
		if (!is_authenticated) {
			return <Redirect to="/login"/>;	
		}
		return (
			<div className="course-wrapper">
				<Header
					currentCourse="Concurrent & Distributed Programming"
					courses={[
						"Computer Networks",
						"Structured Programming",
						"Programming as Problem Solving",
					]}
				/>
				{this.props.match.params.course_no}
			</div>
		);
	}
}

const CourseScreen = connect(
	null,
	null
)(Course);

export default CourseScreen;
