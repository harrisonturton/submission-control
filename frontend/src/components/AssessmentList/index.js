import React from "react";
import { Link } from "react-router-dom";
import { AssessmentItem, Loader } from "components";
import PropTypes from "prop-types";
import "./style.css";

const NotFound = () => (
	<span className="not-found">There is nothing here.</span>
);

const List = ({ items }) => (
	<div>
		{items.map((item, i) => (
			<Link 
				key={i}
				to={`/course/${item.course_id}/${item.id}`}
			>
				<AssessmentItem
					name={item.name}
					due_date={new Date()}
					comments={item.comments}
				/>
			</Link>
		))}
	</div>
);

const AssessmentList = ({ title, subtitle, items }) => {
	let has_items = items.length > 0;
	return (
		<div className="assessment-list-wrapper">
			<div className="assessment-list-header">
				<span className="assessment-list-title">{title}</span>
				<span className="assessment-list-subtitle">{subtitle}</span>
			</div>
			{has_items ? <List items={items}/> : <NotFound/>}
		</div>
	);
};

AssessmentList.propTypes = {
	title: PropTypes.string.isRequired,
	subtitle: PropTypes.string.isRequired,
	items: PropTypes.array.isRequired
};

export default AssessmentList;
