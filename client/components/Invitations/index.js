import React, { Component } from 'react';
import { connect } from 'react-redux';
import CopyToClipboard from 'react-copy-to-clipboard';
import { Row,Col,Table,Button,Modal } from 'react-bootstrap';

import {
  toggleInvitationFormVisibility,
  toggleInvitationDeleteConfirmation,
  submitInvitationDelete
} from '../../actions/invitation';

import {
  flashOperationResult
} from '../../actions/general';

import {
  INVITATION_FORM_NEW_MODE,
  INVITATION_FORM_EDIT_MODE,
  INVITATION_STATUS_NOT_SENT,
  INVITATION_STATUS_SENT,
  INVITATION_STATUS_REPLIED_ATTENDING,
  INVITATION_STATUS_REPLIED_NOT_ATTENDING
} from '../../constants';

import {
  translateStatusCode,
  formatDateForDisplay
} from '../../helpers';

import InvitationForm from '../InvitationForm';

class Invitations extends Component {
    categoryNameByID(categories, categoryID) {
      for (let i = 0; i < categories.length; i++) {
        if (categories[i].id === categoryID) {
          return categories[i].tag;
        }
      }

      return "Unknown";
    }

    constructPrivateLink(privateID) {
      return `${location.origin}/rsvp/${privateID}`;
    }

    render() {
      return <div className="panel-body">
        <Row className="margin-bottom-sm">
          <Col lg={8}>
            <h4>Invitations Not Replied</h4>
          </Col>

          <Col lg={4}>
            <div className="text-right">
              <Button bsStyle="primary" bsSize="small" className="margin-top-sm" onClick={this.props.onAddInvitationClick}><i className="fa fa-plus fa-fw fa-lg"></i> Add Invitation</Button>
            </div>
          </Col>
        </Row>

        {this.props.deleteInvitationConfirmation.visible && <div className="static-modal">
          <Modal.Dialog>
            <Modal.Header>
              <Modal.Title>Delete Invitation {this.props.deleteInvitationConfirmation.invitationID}</Modal.Title>
            </Modal.Header>

            <Modal.Body>
              Are you sure you want to delete this invitation?
            </Modal.Body>

            <Modal.Footer>
              <Button bsSize="sm" onClick={() => {this.props.onToggleDeleteInvitation(0)}}>Close</Button>
              <Button bsSize="sm" bsStyle="danger" onClick={() => {this.props.onDeleteInvitation(this.props.deleteInvitationConfirmation.invitationID)}}>Delete</Button>
            </Modal.Footer>
          </Modal.Dialog>
        </div>}

        {(() => {
          if (this.props.invitationForm && this.props.invitationForm.initialValues) {
            /* TODO:: Find some way to provide initial category if present */
            let currentCategoryID = this.props.invitationForm.initialValues.categoryID;
            if (this.props.categories && this.props.categories.length > 0 && (!currentCategoryID || currentCategoryID === 0)) {
              currentCategoryID = this.props.categories[0].id; 
            }

            return <InvitationForm 
              mode={this.props.invitationForm.mode}
              categories={this.props.categories}
              invitationID={this.props.invitationForm.initialValues.id}
              initialValues={Object.assign({}, this.props.invitationForm.initialValues, {
                categoryID: currentCategoryID,
              })} />;
          }
        })()}

        <Table striped={true} responsive>
          <thead>
            <tr>
              <th>Category</th>
              <th>Greeting</th>
              <th>Max guests no.</th>
              <th>Notes</th>
              <th>Mobile Number</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {this.props.invitations && this.props.invitations.map((invitation) => {
              if (invitation.status == INVITATION_STATUS_NOT_SENT || invitation.status == INVITATION_STATUS_SENT) {
                return <tr key={invitation.id}>
                  <td>{this.categoryNameByID(this.props.categories, invitation.categoryID)}</td>
                  <td>{invitation.greeting}</td>
                  <td>{invitation.maximumGuestCount}</td>
                  <td>{invitation.notes || '-'}</td>
                  <td>{invitation.mobilePhoneNumber}</td>
                  <td>{translateStatusCode(invitation.status)}<p><small>{formatDateForDisplay(invitation.updatedAt)}</small></p></td>
                  <td>
                    <Button bsStyle="danger" bsSize="xs" className="margin-right-sm" onClick={() => {this.props.onToggleDeleteInvitation(invitation.id)}}>Delete</Button>
                    <Button bsStyle="success" bsSize="xs" className="margin-right-sm">Send RSVP</Button>
                    <CopyToClipboard text={this.constructPrivateLink(invitation.privateID)} onCopy={() => this.props.onCopySuccess(invitation.greeting, invitation.privateID)}>
                      <Button bsStyle="primary" bsSize="xs" className="margin-right-sm">Copy Link</Button>
                    </CopyToClipboard>
                    <Button bsStyle="default" bsSize="xs" onClick={() => {this.props.onEditInvitationClick(invitation)}}>Edit</Button>
                  </td>
                </tr>;
              }
            })}
          </tbody>
        </Table>

        <Row className="margin-bottom-sm">
          <Col lg={12}>
            <h4>Invitations Replied</h4>
          </Col>
        </Row>

        <Table striped={true} responsive>
          <thead>
            <tr>
              <th>Category</th>
              <th>Greeting</th>
              <th>Max guests no.</th>
              <th>Notes</th>
              <th>Mobile Number</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {this.props.invitations && this.props.invitations.map((invitation) => {
              if (invitation.status == INVITATION_STATUS_REPLIED_ATTENDING || invitation.status == INVITATION_STATUS_REPLIED_NOT_ATTENDING) {
                return <tr key={invitation.id}>
                  <td>{this.categoryNameByID(this.props.categories, invitation.categoryID)}</td>
                  <td>{invitation.greeting}</td>
                  <td>{invitation.maximumGuestCount}</td>
                  <td>{invitation.notes || '-'}</td>
                  <td>{invitation.mobilePhoneNumber}</td>
                  <td>{translateStatusCode(invitation.status)}<p><small>{formatDateForDisplay(invitation.updatedAt)}</small></p></td>
                  <td>
                    <Button bsStyle="danger" bsSize="xs" className="margin-right-sm" onClick={() => {this.props.onToggleDeleteInvitation(invitation.id)}}>Delete</Button>
                    <Button bsStyle="success" bsSize="xs" className="margin-right-sm">Send RSVP</Button>
                    <CopyToClipboard text={this.constructPrivateLink(invitation.privateID)} onCopy={() => this.props.onCopySuccess(invitation.greeting, invitation.privateID)}>
                      <Button bsStyle="primary" bsSize="xs" className="margin-right-sm">Copy Link</Button>
                    </CopyToClipboard>
                    <Button bsStyle="default" bsSize="xs" onClick={() => {this.props.onEditInvitationClick(invitation)}}>Edit</Button>
                  </td>
                </tr>;
              }
            })}
          </tbody>
        </Table>
      </div>;
    }
}

const mapStateToProps = (state) => {
  return {
    invitationForm: state.invitationForm,
    deleteInvitationConfirmation: state.deleteInvitationConfirmation
  };
};

const mapDispatchToProps = (dispatch) => {
  return {
    onAddInvitationClick: () => {
      dispatch(toggleInvitationFormVisibility(INVITATION_FORM_NEW_MODE, {}))
    },
    onEditInvitationClick: (invitation) => {
      dispatch(toggleInvitationFormVisibility(INVITATION_FORM_EDIT_MODE, invitation))
    },
    onToggleDeleteInvitation: (invitationID) => {
      dispatch(toggleInvitationDeleteConfirmation(invitationID))
    },
    onDeleteInvitation: (invitationID) => {
      dispatch(submitInvitationDelete(invitationID))
    },
    onCopySuccess: (greeting, privateID) => {
      let copySuccessMessage = 'Dear ' + greeting + ', we are tying the knot and would love to have you to attend our wedding lunch on Saturday, 18th Feb 2017!\n'
        + 'Please find the RSVP form and further details at https://jennykevinweddingbells.com/rsvp/' + privateID + '\n'
        + '- Jenny & Kevin'
      dispatch(flashOperationResult(copySuccessMessage, true))
    }
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(Invitations);
