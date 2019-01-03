import React, { Component } from 'react';
import logo from './logo.svg';
import search from "./search.png";

class App extends Component {
  render() {
    return (
      <div className="app">
		  <header>
			  {/* Header Left */}
			  <div>
				  <h1>Concurrent & Distributed Systems</h1>
			  </div>
			  {/* Header Right */}
			  <div>
				  <h4>LOGOUT</h4>
			  </div>
		  </header>
		  <div className="tab-bar">
			  {/* Tab Left */}
			  <div>
				  <ul className="tabs">
					  <li><a href="">Overview</a></li>
					  <li><a href="">Participation</a></li>
					  <li><a href="">Feedback</a></li>
				  </ul>
			  </div>
			  {/* Tab Right */}
			  <div>
				  <div className="searchbar">
					  <input type="text" className="searchbar-input"/>
					  <img src={search} className="search-icon"/>
				  </div>
			  </div>
		  </div>
      </div>
    );
  }
}

export default App;
