import React, { Component } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { WithHeader } from "containers";
import { AssessmentList } from "components";
import "./style.css";

const AdminButton = ({ title, subtitle }) => (
	<a href="/" className="admin">
		<span className="title">{title}</span>
		<span className="subtitle">{subtitle}</span>
	</a>
);

class _ConvenorCourseScreen extends Component {
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
	renderAdmin() {
		return (
			<div className="column-right admin-wrapper">
				<div className="admin-header">
					<span className="assessment-list-title">Course Administation</span>
					<span className="assessment-list-subtitle">View or Submit</span>
				</div>
				<div className="admin-buttons">
					<AdminButton
						title="Manage Tutorial Groups"
						subtitle="Assignment times, students & tutors"
					/>
					<AdminButton
						title="Manage Students"
						subtitle="Enrol or remove students"
					/>
					<AdminButton
						title="Manage Tutors"
						subtitle="Add or remove tutors"
					/>
					<AdminButton
						title="Manage Test Suite"
						subtitle="Create & debug test configurations"
					/>
					<AdminButton
						title="Manage Student Feedback"
						subtitle="Create & schedule feedback items"
					/>
				</div>
			</div>
		);
	}
	filterData(courseID, assignments, labs) {
		const withCourseID = courseID => item => item.course_id == courseID;
		return {
			assignments: assignments.filter(withCourseID(courseID)),
			labs: labs.filter(withCourseID(courseID))
		};
	}
	render() {
		var { courseID, assignments, labs } = this.props;
		var { assignments, labs } = this.filterData(courseID, assignments, labs);
		return (
			<WithHeader className="column-parent" currentCourseID={courseID}>
				{this.renderAssessment(assignments, labs)}
				{this.renderAdmin()}
			</WithHeader>
		);
	}
}

_ConvenorCourseScreen.propTypes = {
	assignments: PropTypes.array.isRequired,
	labs:        PropTypes.array.isRequired,
};

const mapStateToProps = state => ({
	assignments: state.data.assessment.assignments,
	labs:        state.data.assessment.labs,
});

const ConvenorCourseScreen = connect(
	mapStateToProps,
	null
)(_ConvenorCourseScreen);

export default ConvenorCourseScreen;
