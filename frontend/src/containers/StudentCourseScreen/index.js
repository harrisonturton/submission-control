import React, { Component } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { WithHeader } from "containers";
import { AssessmentList, FeedbackList } from "components";
import "./style.css";

class _StudentCourseScreen extends Component {
	renderAssessment(assignments, labs) {
		return (
			<div className="column-left">
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
			<div className="column-right">
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
		var { courseID, assignments, labs, submissions } = this.props;
		var { assignments, labs, submissions } = this.filterData(courseID, assignments, labs, submissions);
		return (
			<WithHeader className="column-parent" currentCourseID={courseID}>
				{this.renderAssessment(assignments, labs)}
				{this.renderFeedback(submissions)}
			</WithHeader>
		);
	}
}

_StudentCourseScreen.propTypes = {
	courseID:    PropTypes.number.isRequired,
	assignments: PropTypes.array.isRequired,
	labs:        PropTypes.array.isRequired,
	submissions: PropTypes.array.isRequired,
};

const mapStateToProps = state => ({
	assignments: state.data.assessment.assignments,
	labs:        state.data.assessment.labs,
	submissions: state.data.submissions
});

const StudentCourseScreen = connect(
	mapStateToProps,
	null
)(_StudentCourseScreen);

export default StudentCourseScreen;
