import React, { Component } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import {
	Header,
	AssessmentItemList,
	AssessmentFeedbackList,
} from "components";
import "./style.css";

class Course extends Component {
	renderAssessment = (course_id, assignments, labs) => (
		<div className="assessment-wrapper"	>
			<AssessmentItemList
				title="Upcoming Assignments"
				subtitle="View or Submit"
				items={assignments.filter(ass => ass.course_id === course_id)}
			/>
			<AssessmentItemList
				title="Upcoming Labs"
				subtitle="View or Submit"
				items={labs.filter(lab => lab.course_id === course_id)}
			/>
		</div>
	);
	renderFeedback = (course_id, submissions) => (
		<div className="feedback-wrapper">
			<AssessmentFeedbackList
				title="Assessment Feedback"
				subtitle=""
				submissions={submissions}
			/>
		</div>
	);
	render() {
		let { course_id } = this.props.match.params;
		let { courses, assignments, labs, submissions } = this.props;
		console.log(JSON.stringify({courses, assignments, labs, submissions}));
		let current_course = "";
		courses.forEach(course => {
			if (course.id === course_id) {
				current_course = course.name;	
			}	
		});
		return (
			<div className="course-wrapper">
				<Header
					currentCourse={current_course}
					courses={courses.filter(course => course.id !== course_id)}
				/>
				{this.renderAssessment(course_id, assignments, labs)}
				{this.renderFeedback(course_id, submissions)}
			</div>
		);
	}
}

Course.propTypes = {
	courses:     PropTypes.array.isRequired,
	assignments: PropTypes.array.isRequired,
	labs:        PropTypes.array.isRequired,
	submissions: PropTypes.array.isRequired,
};

const mapStateToProps = state => ({
	courses:     state.data.courses,
	assignments: state.data.assessment.assignments,
	labs:        state.data.assessment.labs,
	submissions: state.data.assessment.submissions
});

const CourseScreen = connect(
	mapStateToProps,
	null
)(Course);

export default CourseScreen;
