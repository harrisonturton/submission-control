import React, { Component } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { Redirect } from "react-router-dom";
import { WithHeader, Header } from "containers";
import { AssessmentList, FeedbackList } from "components";
import "./style.css";

class Course extends Component {
	renderAssessment = (assignments, labs) => (
		<div className="assessment-wrapper">
			<AssessmentList
				title="Upcoming Assignments"
				subtitle=""
				items={assignments}
			/>
			<AssessmentList
				title="Upcoming Labs"
				subtitle=""
				items={labs}
			/>
		</div>
	);
	renderFeedback = submissions => (
		<div className="feedback-wrapper">
			<FeedbackList
				title="Assessment Feedback"
				subtitle=""
				submissions={submissions}
			/>
		</div>
	);
	render() {
		let { course_id } = this.props.match.params;
		let { is_authenticated, courses, assignments, labs, submissions } = this.props;
		let filtered_assignments = assignments.filter(ass => ass.course_id == course_id);
		let filtered_labs = labs.filter(lab => lab.course_id == course_id);
		let filtered_submissions = submissions.filter(sub => sub.course_id == course_id);
		let current_course = "";
		for (var i = 0; i < courses.length; i++) {
			if (courses[i].id == course_id) {
				current_course = courses[i].name;	
			}
		}
		if (!is_authenticated) {
			return <Redirect to="/login"/>;
		}
		return (
			<WithHeader className="assessment-screen" currentCourseID={course_id}>
				{this.renderAssessment(filtered_assignments, filtered_labs)}
				{this.renderFeedback(filtered_submissions)}
			</WithHeader>
		);
	}
}

Course.propTypes = {
	is_authenticated: PropTypes.bool.isRequired,
	courses:          PropTypes.array.isRequired,
	assignments:      PropTypes.array.isRequired,
	labs:             PropTypes.array.isRequired,
	submissions:      PropTypes.array.isRequired,
};

const mapStateToProps = state => {
	console.log("current state: ", JSON.stringify(state));
	return {
		is_authenticated: state.auth.is_authenticated,
		courses:     state.data.courses,
		assignments: state.data.assessment.assignments,
		labs:        state.data.assessment.labs,
		submissions: state.data.submissions
	};
};

const CourseScreen = connect(mapStateToProps, null)(Course);

export default CourseScreen;
