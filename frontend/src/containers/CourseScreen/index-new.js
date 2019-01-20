import React, { Component } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { Redirect } from "react-router-dom";
import { Header } from "containers";
import { AssessmentList, FeedbackList } from "components";
import "./style.css";

const renderAssessment = (assignments, labs) => (
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

const renderFeedback = submissions => (
	<div className="feedback-wrapper">
		<FeedbackList
			title="Assessment Feedback"
			subtitle=""
			submissions={submissions}
		/>
	</div>
);

const _CourseScreen = ({
	match,
	is_authenticated,
	courses,
	assignments,
	labs,
	submissions,
	current_course
}) => {
	if (!is_authenticated) {
		return <Redirect to="/login"/>;
	}
	if (assignments === undefined) {
		console.log("Assignments are undefined");
	}
	return (
		<div className="two-column-wrapper">
			<Header
				currentCourse={current_course}
				courses={courses}
			/>
			{renderAssessment(assignments, labs)}
			{renderFeedback(submissions)}
		</div>
	);
}

_CourseScreen.propTypes = {
	is_authenticated: PropTypes.bool.isRequired,	
	current_course:   PropTypes.string.isRequired,
	courses:          PropTypes.array.isRequired,
	assignments:      PropTypes.array.isRequired,
	labs:             PropTypes.array.isRequired,
	submissions:      PropTypes.array.isRequired
};

// hasCourseID is used to filter an array of assessment/submissions to
// find the ones applicable for a specific course.
const hasCourseID = course_id => item => item.course_id == course_id;

const mapStateToProps = (state, ownProps) => {
	let { is_authenticated } = state.auth;
	let { courses, submissions } = state.data;
	let { assignments, labs } = state.data.assessment;
	let { course_id } = ownProps.match.params.course_id;
	return {
		is_authenticated: is_authenticated,
		labs:           labs.filter(hasCourseID(course_id)),
		asssignments:   assignments.filter(hasCourseID(course_id)),
		submissions:    submissions.filter(hasCourseID(course_id)),
		current_course: courses.find(hasCourseID(course_id)).name,
		courses:        courses.filter(course => course.id != course_id),
	};
};

const CourseScreen = connect(
	mapStateToProps,
	null
)(_CourseScreen);

export default CourseScreen;
