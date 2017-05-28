import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Row,Col,Button,Alert } from 'react-bootstrap';

import {
  logoutUser
} from '../../actions/logout';

import {
  fetchRSVPs
} from '../../actions/rsvp';

import {
  fetchCategories
} from '../../actions/category';

import {
  fetchInvitations
} from '../../actions/invitation';

import RSVPs from '../RSVPs';
import Invitations from '../Invitations';
import Categories from '../Categories';

class ControlPanel extends Component {
  componentDidMount() {
    this.props.onFetchRSVPs()
    this.props.onFetchCategories()
    this.props.onFetchInvitations()
  }

  render() {
    return <div>
      <Row className="padding-top-lg">
        <Col lg={10} lgOffset={1}>
          <Row>
            <Col lg={2}><h5>Confirmed Guests: <span className="label label-primary">170</span></h5>
            </Col>

            <Col lg={2}><h5>Invitations Not Sent: <span className="label label-warning">12</span></h5>
            </Col>

            <Col lg={2}><h5>Special Diet: <span className="label label-success">4</span></h5>
            </Col>

            <Col lg={6} className="text-right">
              <span className="margin-right-md">Hello {this.props.username}!</span>
              <Button bsStyle="default" bsSize="small" onClick={this.props.onLogoutClick}>Logout</Button>
            </Col>
          </Row>

          <div className="margin-top-md margin-bottom-md">
            {this.props.operation && <Row>
              <Col lg={6} lgOffset={3}>
                <Alert bsStyle={this.props.operation.success ? 'success' : 'danger'} className="text-center">
                  <p>{this.props.operation.message}</p>
                </Alert>
              </Col>
            </Row>}
          </div>

          <Row>
            <Col lg={12}>
              <div className="tabs-container">
                <ul className="nav nav-tabs">
                  <li className="active">
                    <a data-toggle="tab" href="#tab-1">
                      <img src="/static/icons/Inbox.png" height="32" width="32" /> RSVPs
                    </a>
                  </li>

                  <li>
                    <a data-toggle="tab" href="#tab-2">
                      <img src="/static/icons/Star.png" height="32" width="32" /> Invitations
                    </a>
                  </li>

                  <li>
                    <a data-toggle="tab" href="#tab-3">
                      <img src="/static/icons/Tag.png" height="32" width="32" /> Categories
                    </a>
                  </li>
                </ul>

                <div className="tab-content">
                  <div id="tab-1" className="tab-pane active">
                    <RSVPs rsvps={this.props.rsvps} categories={this.props.categories} />
                  </div>

                  <div id="tab-2" className="tab-pane">
                    <Invitations invitations={this.props.invitations} categories={this.props.categories} />
                  </div>

                  <div id="tab-3" className="tab-pane margin-bottom-md">
                    <Categories categories={this.props.categories} />
                  </div>
                </div>
              </div>
            </Col>
          </Row>
        </Col>
      </Row>
    </div>;
  }
}

const mapStateToProps = (state) => {
  return {
    operation: state.operation,
    username: localStorage.getItem('username'),
    rsvps: state.rsvps,
    categories: state.categories,
    invitations: state.invitations
  };
};

const mapDispatchToProps = (dispatch) => {
  return {
    onFetchRSVPs: () => {
      dispatch(fetchRSVPs())
    },
    onFetchCategories: () => {
      dispatch(fetchCategories())
    },
    onFetchInvitations: () => {
      dispatch(fetchInvitations())
    },
    onLogoutClick: () => {
      dispatch(logoutUser())
    }
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(ControlPanel);
