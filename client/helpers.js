var moment = require('moment');

import {
    INVITATION_STATUS_LOOKUP
} from './constants';

export function translateStatusCode(status) {
    if (status in INVITATION_STATUS_LOOKUP) {
    return INVITATION_STATUS_LOOKUP[status];
    }

    return "Unknown";
}

export function formatDateForDisplay(rfc3339Timestamp) {
    // 8:23pm Fri, 20th Sept
    return moment(rfc3339Timestamp, 'YYYY-MM-DDTHH:mm:ssZ').format('HH:mm a ddd, Do MMM');
}
