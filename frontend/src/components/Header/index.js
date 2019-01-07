import React, { Component } from "react";
import chevronDown from "assets/chevron-down.png";
import "./style.css";

export default class Header extends Component {
	// Header expects an array `courses` of course names
	// the student is enrolled in, and a single prop
	// `currentCourse` which is the title of the current page.
	constructor(props) {
		super(props);
		this.state = { isExpanded: false };
		this.onClick = this.onClick.bind(this);
	}
	onClick() {
		this.setState(prev => ({
			isExpanded: !prev.isExpanded,
		}));
	}
	render() {
		let { isExpanded } = this.state;
		let { courses, currentCourse } = this.props;
		return (
			<header
				className={isExpanded ? "expanded" : ""}
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
						<li><a href="">{course}</a></li>	
					))}
				</ul>
			</header>
		);
	}
}
