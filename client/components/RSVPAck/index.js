import React, { Component } from 'react';
import { Row,Col,Alert } from 'react-bootstrap';

class RSVPAck extends Component {
  render() {
    return <div className="panel padding-bottom-md">
      <div className="panel-body">
        <p className="landing-sub-title margin-top-lg margin-left-md">Dear {this.props.rsvp.fullName} -</p>

        <Row>
          <Col xs={10} xsOffset={1}>
            <Alert bsStyle="success" className="text-center margin-top-lg margin-bottom-md">
              <p>Thank you for RSVP-ing! If there are any changes to be made, please contact the Bride or Groom <i className="fa fa-smile-o fa-fw fa-lg"></i></p>
            </Alert>
          </Col>
        </Row>

        <hr className="small" />

        <Row className="margin-top-md">
          <Col xs={2} xsOffset={3}>
            Attending:
          </Col>

          <Col xs={4}>
            <strong>{this.props.rsvp.attending ? 'Yes' : 'No'}</strong>
          </Col>
        </Row>

        {this.props.rsvp.attending && (() => {
          return <div>
            <Row className="margin-top-md">
              <Col xs={2} xsOffset={3}>
                Guest Count:
              </Col>

              <Col xs={4}>
                <strong>{this.props.rsvp.guestCount}</strong>
              </Col>
            </Row>

            <Row className="margin-top-md">
              <Col xs={2} xsOffset={3}>
                Special Diet:
              </Col>

              <Col xs={4}>
                <strong>{this.props.rsvp.specialDiet ? 'Yes' : 'No'}</strong>
              </Col>
            </Row>

            <Row className="margin-top-md">
              <Col xs={2} xsOffset={3}>
                Remarks:
              </Col>

              <Col xs={4}>
                <strong>{this.props.rsvp.remarks}</strong>
              </Col>
            </Row>

            <Row className="margin-top-md">
              <Col xs={2} xsOffset={3}>
                Mobile Phone Number:
              </Col>

              <Col xs={4}>
                <strong>{this.props.rsvp.mobilePhoneNumber}</strong>
              </Col>
            </Row>
          </div>;
        })()}
      </div>
    </div>;
  }
}

export default RSVPAck;
