import React from "react";
import { Link } from "react-router-dom";
import { AssessmentFeedbackItem } from "components";
import PropTypes from "prop-types";
import "./style.css";

const AssessmentFeedbackList = ({ title, subtitle, items }) => (
	<div className="assessment-feedback-list-wrapper">
		<div className="assessment-feedback-list-header">
			<span className="assessment-list-title">{title}</span>
			<span className="assessment-list-subtitle">{subtitle}</span>
		</div>
		{items.map((item, i) => (
			<Link to={`/course/${item.course_code}`}>
				<AssessmentFeedbackItem
					key={i}
					title={item.title}
					due_date={item.due_date}
					feedback={item.feedback}
				/>
			</Link>
		))}
	</div>
);

AssessmentFeedbackList.propTypes = {
	title: PropTypes.string.isRequired,
	subtitle: PropTypes.string.isRequired,
	items: PropTypes.array.isRequired,
};

export default AssessmentFeedbackList;
