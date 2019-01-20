import React, { Component } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { WithHeader } from "containers";
import { AssessmentList, FeedbackList } from "components";
import "./style.css";

class _CourseScreen extends Component {
	renderAssessment(assignments, labs) {
		return (
			<div className="column">
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
	}
	renderFeedback(submissions) {
		return (
			<div className="column">
				<FeedbackList
					title="Assessment Feedback"
					subtitle=""
					submissions={submissions}
				/>
			</div>
		);
	}
	filterData(courseID, assignments, labs, submissions) {
		const withCourseID = courseID => item => item.course_id == courseID;
		return {
			assignments: assignments.filter(withCourseID(courseID)),
			labs:        labs.filter(withCourseID(courseID)),
			submissions: submissions.filter(withCourseID(courseID))
		};
	}
	render() {
		let { course_id } = this.props.match.params;
		var { assignments, labs, submissions } = this.props;
		var { assignments, labs, submissions } = this.filterData(course_id, assignments, labs, submissions);
		return (
			<WithHeader className="column-parent" currentCourseID={course_id}>
				{this.renderAssessment(assignments, labs)}
				{this.renderFeedback(submissions)}
			</WithHeader>
		);
	}
}

_CourseScreen.propTypes = {
	assignments:      PropTypes.array.isRequired,
	labs:             PropTypes.array.isRequired,
	submissions:      PropTypes.array.isRequired,
};

const mapStateToProps = state => ({
	assignments: state.data.assessment.assignments,
	labs:        state.data.assessment.labs,
	submissions: state.data.submissions
});

const CourseScreen = connect(
	mapStateToProps,
	null
)(_CourseScreen);

export default CourseScreen;
