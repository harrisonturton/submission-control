import React, { Component } from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import { Redirect } from "react-router-dom";
import { hasToken } from "api/auth";
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
		console.log("Rendering home...", JSON.stringify(labs));
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

Home.propTypes = {
	header:      PropTypes.object.isRequired,
	assignments: PropTypes.array.isRequired,
	labs:        PropTypes.array.isRequired,
	feedback:    PropTypes.array.isRequired
};

const mapStateToProps = state => ({
	header: mapHeaderState(state),
	assignments: state.data.assessment.assignments,
	labs: state.data.assessment.labs,
	feedback: []
});

const mapHeaderState = state => {
	let courses = state.data.courses;
	return {
		current_course: courses[0] === undefined ? "" : courses[0].name,
		courses:        courses[0] === undefined ? [] : courses,
	};
};

const mapAssignmentState = state => {
	if (state.data.is_fetching || state.data.failed) {
		return [];
	}
	return state.data.assessment.assignments;
	let assignments = []
	for (var i = 0; i < state.data.assessment.length; i++) {
		let assessment = state.data.assessment[i];
		if (assessment.type === "assignment") {
			assignments.push(assessment)	
		}
	}
	return assignments;
}

const mapLabState = state => {
	if (state.data.is_fetching || state.data.failed) {
		return [];
	}
	let labs = []
	for (var i = 0; i < state.data.assessment.length; i++) {
		let assessment = state.data.assessment[i];
		if (assessment.type === "assignment") {
			labs.push(assessment)	
		}
	}
	return labs;
}

const HomeScreen = connect(
	mapStateToProps,
	null
)(Home);

export default HomeScreen;
