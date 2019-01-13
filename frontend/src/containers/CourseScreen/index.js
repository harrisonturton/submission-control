import React, { Component } from "react";
import { connect } from "react-redux";
import { Link, Redirect } from "react-router-dom";
import { Header, AssessmentItemList } from "components";
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
						<AssessmentItemList
							title="Upcoming Assignments"
							subtitle="View or submit"
							items={[
								{ title: "Harmony & Money",      due: new Date(), course_code: "COMP2300" },
							]}
						/>
					</div>
					<div className="tutorial-wrapper">
						<AssessmentItemList
							title="Upcoming Tutorials"
							subtitle="View or submit"
							items={[
								{ title: "Harmony & Money",      due: new Date(), course_code: "COMP2300" },
								{ title: "Distributed Server",   due: new Date(), course_code: "COMP2100" },
								{ title: "Implicit Concurrency", due: new Date(), course_code: "COMP2400" },
							]}
						/>
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
