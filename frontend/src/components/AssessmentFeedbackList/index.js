import React from "react";
import { Link } from "react-router-dom";
import { AssessmentFeedbackItem, Loader } from "components";
import PropTypes from "prop-types";
import "./style.css";

const NotFound = () => (
	<span className="not-found">There is nothing here.</span>
);

const List = ({ submissions }) => (
	<div>
		{submissions.map((submission, i) => (
			<Link
				key={i}
				to={`/course/${submission.id}`}
			>
				<AssessmentFeedbackItem
					title={submission.title}
					due_date={new Date()}
					feedback={submission.feedback}
				/>
			</Link>
		))}
	</div>
);

const AssessmentFeedbackList = ({ title, subtitle, submissions }) => {
	let has_submissions = submissions.length > 0;
	return (
		<div className="assessment-feedback-list-wrapper">
			<div className="assessment-feedback-list-header">
				<span className="assessment-list-title">{title}</span>
				<span className="assessment-list-subtitle">{subtitle}</span>
			</div>
			{has_submissions ? <List submissions={submissions}/> : <NotFound/>}
		</div>
	);
}

AssessmentFeedbackList.propTypes = {
	title: PropTypes.string.isRequired,
	subtitle: PropTypes.string.isRequired,
	submissions: PropTypes.array.isRequired,
};

export default AssessmentFeedbackList;
