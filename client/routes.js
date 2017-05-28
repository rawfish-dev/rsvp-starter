import React from 'react'
import ReactDOM from 'react-dom'
import thunkMiddleware from 'redux-thunk'
import { createStore,combineReducers,applyMiddleware } from 'redux'
import { Provider } from 'react-redux'
import { Router,IndexRedirect,Route,browserHistory } from 'react-router'
import { syncHistoryWithStore,routerReducer } from 'react-router-redux'

import gatheredReducers from './reducers';

import App from '#components/App';
import Login from '#components/Login';
import Details from '#components/Details';
import ControlPanel from '#components/ControlPanel';

import { requireAuthentication } from '#components/Auth';

// Cause css to get loaded
require('./css');

// Add the reducer to your store on the `routing` key
const store = createStore(
  combineReducers(
    Object.assign({}, gatheredReducers, { 
      routing: routerReducer
    })
  ),
  window.devToolsExtension && window.devToolsExtension(),
  applyMiddleware(
    thunkMiddleware, // lets us dispatch() functions
  )
)

// Create an enhanced history that syncs navigation events with the store
const history = syncHistoryWithStore(browserHistory, store)

ReactDOM.render(
  <Provider store={store}>
    { /* Tell the Router to use our enhanced history */ }
    <Router history={history}>
      <Route path="/" component={App}>
        <IndexRedirect to="/rsvp" />
        <Route path="/login" component={Login} />
        <Route path="/rsvp(/:id)" component={Details} />
        <Route path="/control_panel" component={requireAuthentication(ControlPanel)} />
      </Route>
    </Router>
  </Provider>,
  document.getElementById('app')
)
