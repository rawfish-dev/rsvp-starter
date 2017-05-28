import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Row,Col,Table,Button,Modal } from 'react-bootstrap';

import {
  toggleRSVPFormVisibility,
  toggleRSVPDeleteConfirmation,
  submitRSVPDelete
} from '../../actions/rsvp';

import {
  RSVP_FORM_NEW_MODE,
  RSVP_FORM_EDIT_MODE
} from '../../constants';

import RSVPAdminForm from '../RSVPAdminForm';

class RSVPs extends Component {
    render() {
        return <div className="panel-body">
          <Row className="margin-bottom-sm">
            <Col lg={8}>
              <h4>RSVPs</h4>
            </Col>

            <Col lg={4}>
              <div className="text-right">
                <Button bsStyle="primary" bsSize="small" className="margin-top-sm" onClick={this.props.onAddRSVPClick}><i className="fa fa-plus fa-fw fa-lg"></i> Add RSVP</Button>
              </div>
            </Col>
          </Row>

          {this.props.deleteRSVPConfirmation.visible && <div className="static-modal">
            <Modal.Dialog>
              <Modal.Header>
                <Modal.Title>Delete RSVP {this.props.deleteRSVPConfirmation.rsvpID}</Modal.Title>
              </Modal.Header>

              <Modal.Body>
                Are you sure you want to delete this RSVP?
              </Modal.Body>

              <Modal.Footer>
                <Button bsSize="sm" onClick={() => {this.props.onToggleDeleteRSVP(0)}}>Close</Button>
                <Button bsSize="sm" bsStyle="danger" onClick={() => {this.props.onDeleteRSVP(this.props.deleteRSVPConfirmation.rsvpID)}}>Delete</Button>
              </Modal.Footer>
            </Modal.Dialog>
          </div>}

          {(() => {
            if (this.props.rsvpForm !== null && this.props.rsvpForm.initialValues !== null) {
              /* Seems hackish? */
              let currentCategoryID = this.props.rsvpForm.initialValues.categoryID;
              if (!currentCategoryID || currentCategoryID === 0) {
                currentCategoryID = this.props.categories[0].id; 
              }

              return <RSVPAdminForm 
                mode={this.props.rsvpForm.mode}
                categories={this.props.categories}
                rsvpID={this.props.rsvpForm.initialValues.id}
                initialValues={Object.assign({}, this.props.rsvpForm.initialValues, {
                  categoryID: currentCategoryID,
                })} />;
            }
          })()}

          <Table striped={true} responsive>
            <thead>
              <tr>
                <th>Category</th>
                <th>Greeting / Full Name</th>
                <th>Attending</th>
                <th>Guest Count</th>
                <th>Special Dietiary</th>
                <th>Remarks</th>
                <th>Mobile Number</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {this.props.rsvps && this.props.rsvps.map((rsvp) => {
                return <tr>
                  <td>Bride & Groom</td>
                  <td>{rsvp.fullName}</td>
                  <td className="text-center">
                    {(() => {
                      if (rsvp.attending) {
                        return <span className="label label-success">Yes</span>;
                      }

                      return <span className="label label-danger">No</span>;
                    })()}
                  </td>
                  <td className="text-center">{rsvp.guestCount}</td>
                  <td  className="text-center">{rsvp.specialDiet ? 'Yes' : 'No'}</td>
                  <td>{rsvp.remarks}</td>
                  <td>{rsvp.mobilePhoneNumber}</td>
                  <td>
                    <Button bsStyle="danger" bsSize="xs" className="margin-right-sm" onClick={() => {this.props.onToggleDeleteRSVP(rsvp.id)}}>Delete</Button>
                    <Button bsStyle="default" bsSize="xs" onClick={() => {this.props.onEditRSVPClick(rsvp)}}>Edit</Button>
                  </td>
                </tr>;
              })}
            </tbody>
          </Table>
        </div>;
    }
}

const mapStateToProps = (state) => {
  return {
    rsvpForm: state.rsvpForm,
    deleteRSVPConfirmation: state.deleteRSVPConfirmation
  };
};

const mapDispatchToProps = (dispatch) => {
  return {
    onAddRSVPClick: () => {
      dispatch(toggleRSVPFormVisibility(RSVP_FORM_NEW_MODE, {}))
    },
    onEditRSVPClick: (rsvp) => {
      dispatch(toggleRSVPFormVisibility(RSVP_FORM_EDIT_MODE, rsvp))
    },
    onToggleDeleteRSVP: (rsvpID) => {
      dispatch(toggleRSVPDeleteConfirmation(rsvpID))
    },
    onDeleteRSVP: (rsvpID) => {
      dispatch(submitRSVPDelete(rsvpID))
    },
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(RSVPs);
