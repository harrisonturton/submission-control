import React, { Component } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { WithHeader } from "containers";
import { SubmissionList } from "components";
import "./style.css";

class _AssessmentScreen extends Component {
	renderLeftColumn(assessmentName) {
		return (
			<div className="column-left">
				<h1>{assessmentName}</h1>
				<p>You can make multiple submissions to a single assignment or lab. Each submission will be tested against the test suite, and you’ll see the result. Your tutors and convenors can see & give feedback on any submission you make.</p>
				<form className="assessment-form">
					<label>Title</label>
					<input
						name="title"
						type="text"
						placeholder="Title of your submission"
					/>
					<label>Comments</label>
					<textarea
						name="comments"
						type="text"
						placeholder="Some more info about your submission"
					/>
					<label>Zipped submission folder</label>
					<input type="file"/>
					<input type="submit" value="Submit"/>
				</form>
			</div>
		);
	}
	renderRightColumn(courseID, assessmentID, submissions) {
		return (
			<div className="column-right">
				<SubmissionList
					title="Your Submissions"
					subtitle=""
					courseID={courseID}
					assessmentID={assessmentID}
					submissions={submissions}
				/>
			</div>
		);
	}
	render() {
		var { match, assessment, submissions } = this.props;
		let { course_id, assessment_id } = match.params;
		var submissions = submissions.filter(sub => sub.assessment_id == assessment_id);
		let assessmentName = assessment.find(ass => ass.id == assessment_id).name;
		return (
			<WithHeader className="column-parent" currentCourseID={course_id}>
				{this.renderLeftColumn(assessmentName)}
				{this.renderRightColumn(course_id, assessment_id, submissions)}
			</WithHeader>
		);
	}
}

_AssessmentScreen.propTypes = {
	assessment:       PropTypes.array.isRequired,
	submissions:      PropTypes.array.isRequired,
};

const mapStateToProps = state => ({
	assessment: [ ...state.data.assessment.assignments, ...state.data.assessment.labs ],
	submissions: state.data.submissions,
});

const AssessmentScreen = connect(
	mapStateToProps,
	null
)(_AssessmentScreen);

export default AssessmentScreen;
