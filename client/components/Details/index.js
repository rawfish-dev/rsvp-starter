import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Row,Col } from 'react-bootstrap';
import Scroll, { Link,Element } from 'react-scroll';

import RSVPForm from '../RSVPForm';
import RSVPAck from '../RSVPAck';

import {
  fetchRSVP
} from '../../actions/guest';

class Details extends Component {

  constructor(props) {
    super(props);
    this.props.loadApplicationState(this.props.params.id);
  }

  render() {
    return <div className="full-height">
      <header id="top" className="header">
        <div className="text-vertical-center fade-in-5">
          <span className="landing-title">## Names ##</span>
          <p className="landing-sub-title margin-top-lg">## Event Title ##</p>

          <Link to="details" spy={true} smooth={true} offset={50} duration={500} className="btn btn-light-rounded btn-lg margin-right-sm">Details</Link>
          {this.props.guestRSVP && <Link to="form" spy={true} smooth={true} offset={70} duration={2200} className="btn btn-dark-rounded btn-lg">RSVP</Link>}
        </div>
      </header>

      <Element name="details">
        <section className="about">
          <div className="container">
            <Row>
              <Col lg={12} className="text-center">
                <h2>## Day, Date ## @ ## Event Location ##</h2>
              </Col>
            </Row>
          </div>
        </section>
      </Element>

      <aside className="callout">
        <div className="text-vertical-center">
          <h1>## Event Location Details ##</h1>
        </div>
      </aside>

      <section id="services" className="services details-bg">
        <div className="container">
          <Row className="text-center">
            <Col lg={10} lgOffset={1}>
              <Row>
                <Col md={4} sm={12}>
                  <div className="service-item">
                    <span className="fa-stack fa-4x">
                      <i className="fa fa-circle fa-stack-2x"></i>
                      <i className="fa fa-glass fa-stack-1x text-dark"></i>
                    </span>
                    <h4>
                      <strong>## Item 1 ##</strong>
                    </h4>
                    <p>## Item 1 Details ##</p>
                  </div>
                </Col>

                <Col md={4} sm={12}>
                  <div className="service-item">
                    <span className="fa-stack fa-4x">
                      <i className="fa fa-circle fa-stack-2x"></i>
                      <i className="fa fa-ship fa-stack-1x text-dark"></i>
                    </span>
                    <h4>
                      <strong>## Item 2 ##</strong>
                    </h4>
                    <p>## Item 2 Details ##</p>
                  </div>
                </Col>

                <Col md={4} sm={12} className="col-md-4 col-sm-12">
                  <div className="service-item">
                    <span className="fa-stack fa-4x">
                      <i className="fa fa-circle fa-stack-2x"></i>
                      <i className="fa fa-rocket fa-stack-1x text-dark"></i>
                    </span>
                    <h4>
                      <strong>## Item 3 ##</strong>
                    </h4>
                    <p>## Item 3 Details ##</p>
                  </div>
                </Col>
              </Row>
            </Col>
          </Row>
        </div>
      </section>

      {/* Google Maps Section
      <section id="contact" className="map">
        <iframe width="100%" height="100%" frameBorder="0" scrolling="no" marginHeight="0" marginWidth="0" src="## Google Map Location URL ##"></iframe>
        <br />
        <small>
          <a href="## Google Map Location URL ##"></a>
        </small>
      </section>*/}

      <section className="margin-top-lg margin-bottom-lg padding-top-md padding-bottom-md">
        <Row>
          <Col sm={12} className="text-center">
            <h3>
              <i className="fa fa-building-o fa-fw fa-lg"></i> <strong>Event Address</strong>
            </h3>

            <p>## Event Address ## 
              <br />## Event Postal Code ##
            </p>
            <br />

            <h4><i className="fa fa-child fa-fw fa-lg"></i> ## Contact 1 Title ##</h4>
            <ul className="list-unstyled">
              <li>## Contact 1 Phone Number ##</li>
            </ul>
            <br />

            <h4><i className="fa fa-diamond fa-fw fa-lg"></i> ## Contact 2 Title ##</h4>
            <ul className="list-unstyled">
              <li>## Contact 2 Phone Number ##</li>
            </ul>
          </Col>
        </Row>
      </section>


      {this.props.guestRSVP && <Element name="form">
        <section>
          <div className="rsvp-form">
            <Row className="padding-top-lg">
              <Col lg={6} lgOffset={3}>
                <div className="padding-top-lg">
                  {(() => {
                    if (this.props.guestRSVP && this.props.guestRSVP.completed) {
                      return <RSVPAck rsvp={this.props.guestRSVP} />
                    }

                    return <RSVPForm initialValues={this.props.guestRSVP} />
                  })()}  
                </div>
              </Col>
            </Row>

            <div className="text-center margin-top-md padding-bottom-md">
              <p className="text-muted"><i className="fa fa-github fa-fw fa-lg"></i> <a href="https://github.com/rawfish-dev/rsvp-starter">Source Code</a></p>
            </div>
          </div>
        </section>
      </Element>}
    </div>;
  }
}

const mapStateToProps = (state) => {
  return {
    guestRSVP: state.guestRSVP
  };
};

const mapDispatchToProps = (dispatch) => {
  return {
    loadApplicationState: (id) => {
      if (id) {
        dispatch(fetchRSVP(id)); 
      }
    }
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(Details);
