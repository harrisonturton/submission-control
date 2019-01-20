import React from "react";
import PropTypes from "prop-types";
import { Link } from "react-router-dom";
import { SubmissionItem } from "components";

const NotFound = () => (
	<span className="not-found">There is nothing here.</span>
);

const List = ({ courseID, assessmentID, submissions }) => (
	<div>
		{submissions.map((sub, i) => (
			<Link
				key={i}
				to={`/course/${courseID}/${assessmentID}/${sub.id}`}
			>
				<SubmissionItem
					title={sub.title}
					dueDate={new Date()}
					description={sub.description}
					hasFeedback={sub.feedback !== ""}
					feedback={sub.feedback}
				/>
			</Link>
		))}
	</div>
);

const SubmissionList = ({ title, subtitle, submissions, courseID, assessmentID }) => {
	let has_loaded = submissions.length > 0;
	return (
		<div className="assessment-list-wrapper">
			<div className="assessment-list-header">
				<span className="assessment-list-title">{title}</span>
				<span className="assessment-list-subtitle">{subtitle}</span>
			</div>
			{has_loaded
					? <List
						submissions={submissions}
						courseID={courseID}
						assessmentID={assessmentID}
					  />
					: <NotFound/>}
		</div>
	);
};

SubmissionList.propTypes = {
	title:        PropTypes.string.isRequired,
	subtitle:     PropTypes.string.isRequired,
	courseID:     PropTypes.string.isRequired,
	assessmentID: PropTypes.string.isRequired,
	submissions:  PropTypes.array.isRequired
};

export default SubmissionList;
