import React from "react";
import { Link } from "react-router-dom";
import { AssessmentFeedbackItem } from "components";
import PropTypes from "prop-types";
import "./style.css";

const AssessmentFeedbackList = ({ title, subtitle, submissions }) => (
	<div className="assessment-feedback-list-wrapper">
		<div className="assessment-feedback-list-header">
			<span className="assessment-list-title">{title}</span>
			<span className="assessment-list-subtitle">{subtitle}</span>
		</div>
		{submissions.map((submission, i) => (
			<Link to={`/course/${submission.id}`}>
				<AssessmentFeedbackItem
					key={i}
					title={submission.title}
					due_date={new Date()}
					feedback={submission.feedback}
				/>
			</Link>
		))}
	</div>
);

AssessmentFeedbackList.propTypes = {
	title: PropTypes.string.isRequired,
	subtitle: PropTypes.string.isRequired,
	submissions: PropTypes.array.isRequired,
};

export default AssessmentFeedbackList;
