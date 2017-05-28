import React, { Component } from 'react';

class App extends Component {

  render() {
    return <div className="full-height">
      {this.props.children}
    </div>;
  }
}

export default App;
