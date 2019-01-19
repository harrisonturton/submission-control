import React, { Component } from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import { Redirect } from "react-router-dom";
import "./style.css";
import {
	Header,
	AssessmentItemList,
	AssessmentFeedbackList,
} from "components";

class Home extends Component {
	renderAssessment(assignments, labs) {
		return (
			<div className="assessment-wrapper">
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
					submissions={feedback}
				/>
			</div>
		);
	}
	render() {
		let { is_authenticated, header, assignments, labs, feedback } = this.props;
		if (!is_authenticated) {
			return <Redirect to="/login"/>;
		}
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
	is_authenticated: PropTypes.bool.isRequired,
	header:      PropTypes.object.isRequired,
	assignments: PropTypes.array.isRequired,
	labs:        PropTypes.array.isRequired,
	feedback:    PropTypes.array.isRequired
};

const mapStateToProps = state => ({
	is_authenticated: state.auth.is_authenticated,
	header: mapHeaderState(state),
	assignments: state.data.assessment.assignments,
	labs: state.data.assessment.labs,
	feedback: state.data.submissions === undefined ? [] : state.data.submissions
});

const mapHeaderState = state => {
	let courses = state.data.courses;
	return {
		current_course: courses[0] === undefined ? "" : courses[0].name,
		courses:        courses[0] === undefined ? [] : courses,
	};
};

const HomeScreen = connect(
	mapStateToProps,
	null
)(Home);

export default HomeScreen;
