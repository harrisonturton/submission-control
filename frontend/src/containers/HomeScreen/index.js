import React, { Component } from "react";
import { connect } from "react-redux";
import { Redirect } from "react-router-dom";
import { hasToken } from "api";
import "./style.css";
import {
	Header,
	AssessmentItemList,
	AssessmentFeedbackList,
} from "components";

const RedirectToLogin = () => <Redirect to="/login"/>;

class Home extends Component {
	renderAssessment(assignments, labs) {
		return (
			<div className="assignment-wrapper">
				<AssessmentItemList
					title="Upcoming Assignments"
					subtitle=""
					items={assignments}
				/>
				<AssessmentItemList
					title="Upcoming Labs"
					subtitle=""
					items={labs}
				/>
			</div>
		);
	}
	renderFeedback(feedback) {
		return (
			<div className="feedback-wrapper">
				<AssessmentFeedbackList
					title="Assessment Feedback"
					subtitle=""
					items={feedback}
				/>
			</div>
		);
	}
	render() {
		let is_authenticated = hasToken();
		let { header, assignments, labs, feedback } = this.props;
		return (
			<div className="home-wrapper">
				<Header
					currentCourse={header.current_course}
					courses={header.courses}
				/>
				{this.renderAssessment(assignments, labs)}
				{this.renderFeedback(feedback)}
			</div>
		);
	}
}

const mapStateToProps = state => ({
	header: {
		current_course: "Test",
		courses: ["one", "two"],
	},
	assignments: [],
	labs: [],
	feedback: []
})

const HomeScreen = connect(
	mapStateToProps,
	null
)(Home);

const emptyIfUndefined = val => {
	switch (typeof val) {
		case "string":
			return "";
		case "array":
			return []
		default:
			return 0;
	}
}

export default HomeScreen;
