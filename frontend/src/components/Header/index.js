import React, { Component } from "react";
import chevronDown from "assets/chevron-down.png";
import "./style.css";

export default class Header extends Component {
	// props:
	//   courses: array of human-readable course names
	//   currentCourse: the course to display in the header
	constructor(props) {
		super(props);
		this.state = { isExpanded: false };
		this.onClick = this.onClick.bind(this);
		this.onMouseLeave = this.onMouseLeave.bind(this);
	}
	onClick() {
		this.setState(prev => ({
			isExpanded: !prev.isExpanded,
		}));
	}
	onMouseLeave() {
		this.setState(prev => ({
			isExpanded: false
		}));
	}
	render() {
		let { isExpanded } = this.state;
		let { courses, currentCourse } = this.props;
		return (
			<header
				className={isExpanded ? "expanded" : ""}
				onMouseLeave={this.onMouseLeave}
				style={{
					height: isExpanded ? courses.length * 40 + 10 : 15
				}}
			>
				<span class="current" onClick={this.onClick}>
					{currentCourse}
					<img class="chevron" src={chevronDown}/>
				</span>
				<ul>
					{courses.map((course, i) => (
						<li className="course-name"><a href="">{course}</a></li>	
					))}
				</ul>
			</header>
		);
	}
}
