import React, { Component } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import {
	Header,
	AssessmentItemList,
	AssessmentFeedbackList
} from "components";
import "./style.css";

class Course extends Component {
	renderAssessment = (assignments, labs) => (
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
	renderFeedback = submissions => (
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
		return (
			<div className="course-wrapper">
				<Header
					currentCourse={courses.find(course => course.id == course_id).name}
					courses={courses.filter(course => course.id !== course_id)}
				/>
				{this.renderAssessment(assignments, labs)}
				{this.renderFeedback(submissions)}
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

const mapStateToProps = state => {
	console.log("current state: ", JSON.stringify(state));
	return {
		courses:     state.data.courses,
		assignments: state.data.assessment.assignments,
		labs:        state.data.assessment.labs,
		submissions: state.data.submissions
	};
};

const CourseScreen = connect(mapStateToProps, null)(Course);

export default CourseScreen;
