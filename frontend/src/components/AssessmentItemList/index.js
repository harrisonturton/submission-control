import React from "react";
import { Link } from "react-router-dom";
import { AssessmentItem, Loader } from "components";
import PropTypes from "prop-types";
import "./style.css";

const LoadingList = () => (
	<Loader/>
);

const LoadedList = ({ items }) => (
	<div>
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

const AssessmentItemList = ({ title, subtitle, items }) => (
	<div className="assessment-list-wrapper">
		<div className="assessment-list-header">
			<span className="assessment-list-title">{title}</span>
			<span className="assessment-list-subtitle">{subtitle}</span>
		</div>
		{items === undefined ? <LoadingList/> : <LoadedList items={items}/>}
	</div>
);

AssessmentItemList.propTypes = {
	title: PropTypes.string.isRequired,
	subtitle: PropTypes.string.isRequired,
	items: PropTypes.array.isRequired
};

export default AssessmentItemList;
