import React from "react";
import { Link } from "react-router-dom";
import { AssessmentItem } from "components";
import PropTypes from "prop-types";
import "./style.css";

const AssessmentItemList = ({ title, subtitle, items }) => (
	<div className="assessment-list-wrapper">
		<div className="assessment-list-header">
			<span className="assessment-list-title">{title}</span>
			<span className="assessment-list-subtitle">{subtitle}</span>
		</div>
		{items.map((item, i) => (
			<Link to={`/course/${item.course_code}`}>
				<AssessmentItem
					key={i}
					name={item.name}
					due_date={new Date()}
					comments={item.comments}
				/>
			</Link>
		))}
	</div>
);

AssessmentItemList.propTypes = {
	title: PropTypes.string.isRequired,
	subtitle: PropTypes.string.isRequired,
	items: PropTypes.array.isRequired,
	// Shape: { title, course_code, due_date  }
};

export default AssessmentItemList;
