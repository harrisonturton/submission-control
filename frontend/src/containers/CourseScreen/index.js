import React, { Component } from "react";
import { connect } from "react-redux";
import { Link, Redirect } from "react-router-dom";
import { Header, AssessmentItem } from "components";
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
				<div className="assessment-wrapper">
					<div className="assignment-wrapper">
						<span className="subtitle">Upcoming Assignments</span>
						<Link to="/">
							<AssessmentItem title="Harmony & Money" due={new Date()} comments="2 submissions with 8 warnings."/>
						</Link>
						<Link to="/">
							<AssessmentItem title="Distributed Server" due={new Date()} comments="2 submissions with 8 warnings on the most recent build."/>
						</Link>
						<Link to="/">
							<AssessmentItem title="Harmony & Money" due={new Date()} comments="1 submission that passes all tests!"/>
						</Link>
					</div>
					<div className="tutorial-wrapper">
						<span className="subtitle">Upcoming Tutorials</span>
						<Link to="/">
							<AssessmentItem title="Distributed Server" due={new Date()} comments="2 submissions with 8 warnings on the most recent build."/>
						</Link>
						<Link to="/">
							<AssessmentItem title="Harmony & Money" due={new Date()} comments="1 submission that passes all tests!"/>
						</Link>
					</div>
				</div>
			</div>
		);
	}
}

const CourseScreen = connect(
	null,
	null
)(Course);

export default CourseScreen;
